package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"encoding/json"

	"github.com/UTDNebula/nebula-api/rest/configs"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
)

const ClubsClubJSONB = `jsonb_build_object(
    'slug', slug,
    'id', club.id,
    'name', club.name,
    'description', club.description,
    'tags', tags,
    'profile_image', profile_image,
    'updated_at', (updated_at AT TIME ZONE 'UTC'),
    'officers', officers,
    'contacts', contacts
)`

const ClubsContactJSONB = `jsonb_build_object(
    'platform', contacts.platform,
    'url', contacts.url
)`

const ClubsOfficerJSONB = `jsonb_build_object(
    'name', officers.name,
    'position', officers.position
)`

const ClubsEventJSONB = `jsonb_build_object(
    'id', events.id,
    'club_id', events.club_id,
    'name', events.name,
    'description', events.description,
    'start_time', (start_time AT TIME ZONE 'UTC'),
    'end_time', (end_time AT TIME ZONE 'UTC'),
    'location', events.location,
    'image', events.image,
    'created_at', (created_at AT TIME ZONE 'UTC'),
    'updated_at', (updated_at AT TIME ZONE 'UTC')
)`

// @Id				clubById
// @Router			/clubs/{id} [get]
// @Tags			Clubs
// @Description	Returns the listing info for the club with given ID
// @Produce		json
// @Param			id	path		string							true	"ID of the club to get"
// @Success		200	{object}	schema.APIResponse[schema.Club]	"A club"
// @Failure		500	{object}	schema.APIResponse[string]		"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]		"A string describing the error"
func ClubById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var clubsDatabase *sql.DB = configs.ConnectClubsDB()
	id := c.Param("id")

	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, `
        SELECT jsonb_agg(`+ClubsClubJSONB+`)
        FROM club
            JOIN LATERAL (
                SELECT jsonb_agg(`+ClubsContactJSONB+`
                        ORDER BY contacts.platform
                    ) AS contacts
                FROM contacts
                WHERE contacts.club_id = club.id
            ) AS contacts ON TRUE
            JOIN LATERAL (
                SELECT jsonb_agg(`+ClubsOfficerJSONB+`) AS officers
                FROM officers
                WHERE officers.club_id = club.id
            ) AS officers ON TRUE
        WHERE club.id = $1;`,
		id).Scan(&raw)

	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if raw == nil {
		respond(c, http.StatusNotFound, "error", "Club not found")
		return
	}

	var clubs []schema.Club
	if err := json.Unmarshal(raw, &clubs); err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Return single club since filtering by unique ID
	respond(c, http.StatusOK, "success", clubs[0])
}

// @Id				clubEvents
// @Router			/clubs/{id}/events [get]
// @Tags			Clubs
// @Description	Returns the upcoming events for the club with given ID
// @Produce		json
// @Param			id	path		string									true	"ID of the club to get events for"
// @Success		200	{object}	schema.APIResponse[schema.ClubsEvent]	"An event"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]				"A string describing the error"
func ClubEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var clubsDatabase *sql.DB = configs.ConnectClubsDB()
	id := c.Param("id")

	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, `
        SELECT jsonb_agg(`+ClubsEventJSONB+`) AS events
        FROM events
        WHERE events.club_id = $1
            AND events.approved = 'approved'::status_enum;`,
		id).Scan(&raw)

	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if raw == nil {
		respond(c, http.StatusNotFound, "error", "Club not found")
		return
	}

	var clubsEvents []schema.ClubsEvent
	if err := json.Unmarshal(raw, &clubsEvents); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", clubsEvents)
}

// @Id				clubSearch
// @Router			/clubs/search [get]
// @Tags			Clubs
// @Description	Returns list of clubs matching the search string
// @Produce		json
// @Param			q	query		string								true	"Search string"
// @Success		200	{object}	schema.APIResponse[[]schema.Club]	"List of matching clubs"
// @Failure		500	{object}	schema.APIResponse[string]			"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]			"A string describing the error"
func ClubSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var clubsDatabase *sql.DB = configs.ConnectClubsDB()
	search := c.Query("q")

	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, `
        SELECT jsonb_agg(`+ClubsClubJSONB+`
                ORDER BY (
                    (20.0 * word_similarity($1, coalesce(club.alias, '')))
                    + (10.0 * word_similarity($1, club.name))
                    + (5.0 * word_similarity($1, coalesce(array_to_string(club.tags, ' '), '')))
                    + (coalesce(-1.0 * (club.search_tsv <@> to_bm25query(to_tsvector('english', $1), 'club_search_idx')), 0.0))
                ) DESC
            ) AS club
        FROM club
            JOIN LATERAL (
                SELECT jsonb_agg(`+ClubsContactJSONB+`
                        ORDER BY contacts.platform
                    ) AS contacts
                FROM contacts
                WHERE contacts.club_id = club.id
            ) AS contacts ON TRUE
            JOIN LATERAL (
                SELECT jsonb_agg(`+ClubsOfficerJSONB+`) AS officers
                FROM officers
                WHERE officers.club_id = club.id
            ) AS officers ON TRUE
        WHERE (
                club.search_tsv @@ websearch_to_tsquery('english', $1)
                OR word_similarity($1, club.name) >= 0.2
                OR word_similarity($1, COALESCE(club.alias, '')) >= 0.2
                OR club.name ILIKE '%' || $1 || '%'
                OR club.alias ILIKE '%' || $1 || '%'
            )
            AND 'approved' = 'approved'::approved_enum`,
		search).Scan(&raw)

	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if raw == nil {
		respond(c, http.StatusNotFound, "error", "Club not found")
		return
	}

	var clubs []schema.Club
	if err := json.Unmarshal(raw, &clubs); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", clubs)
}

// @Id				clubsEventById
// @Router			/clubs/events/{id} [get]
// @Tags			Clubs
// @Description	Returns the listing info for the event with given ID
// @Produce		json
// @Param			id	path		string									true	"ID of the event to get"
// @Success		200	{object}	schema.APIResponse[schema.ClubsEvent]	"An event"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]				"A string describing the error"
func ClubsEventById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var clubsDatabase *sql.DB = configs.ConnectClubsDB()
	id := c.Param("id")

	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, `
        SELECT jsonb_agg(`+ClubsEventJSONB+`) AS events
        FROM events
        WHERE events.id = $1`,
		id).Scan(&raw)

	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if raw == nil {
		respond(c, http.StatusNotFound, "error", "Event not found")
		return
	}

	var clubsEvents []schema.ClubsEvent
	if err := json.Unmarshal(raw, &clubsEvents); err != nil {
		respondWithInternalError(c, err)
		return
	}

	// Return single event since filtering by unique ID
	respond(c, http.StatusOK, "success", clubsEvents[0])
}

// @Id				clubsEventSearch
// @Router			/clubs/events/search [get]
// @Tags			Clubs
// @Description	Returns list of events matching the search string
// @Produce		json
// @Param			q	query		string									true	"Search string"
// @Success		200	{object}	schema.APIResponse[[]schema.ClubsEvent]	"List of matching events"
// @Failure		500	{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]				"A string describing the error"
func ClubsEventSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var clubsDatabase *sql.DB = configs.ConnectClubsDB()
	search := c.Query("q")

	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, `
		SELECT jsonb_agg(`+ClubsEventJSONB+`
                ORDER BY (
                    ts_rank(search_tsv, websearch_to_tsquery('english', $1))
                ) DESC
            ) AS events
        FROM events
        WHERE events.search_tsv @@ (websearch_to_tsquery('english', $1))
            AND events.approved = 'approved'::status_enum`,
		search).Scan(&raw)

	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if raw == nil {
		respond(c, http.StatusNotFound, "error", "Event not found")
		return
	}

	var clubsEvents []schema.ClubsEvent
	if err := json.Unmarshal(raw, &clubsEvents); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", clubsEvents)
}
