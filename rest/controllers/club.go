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
    'id', club.id,
    'slug', club.slug,
    'name', club.name,
    'alias', club.alias,
    'founding_date', (club.founding_date AT TIME ZONE 'UTC'),
    'updated_at', (club.updated_at AT TIME ZONE 'UTC'),
    'description', club.description,
    'tags', club.tags,
    'profile_image', club.profile_image,
    'banner_image', club.banner_image,
    'schools', club.schools,
    'officers', officers,
    'contacts', contacts,
    'membership_forms', membership_forms
)`

const ClubsOfficerJSONB = `jsonb_build_object(
    'name', officers.name,
    'position', officers.position
)`

const ClubsContactJSONB = `jsonb_build_object(
    'platform', contacts.platform,
    'url', contacts.url
)`

const ClubsMembershipFormJSONB = `jsonb_build_object(
    'name', membership_forms.name,
    'url', membership_forms.url
)`

const ClubsEventJSONB = `jsonb_build_object(
    'id', events.id,
    'club_id', events.club_id,
    'name', events.name,
    'description', events.description,
    'start_time', (events.start_time AT TIME ZONE 'UTC'),
    'end_time', (events.end_time AT TIME ZONE 'UTC'),
    'location', events.location,
    'image', events.image,
    'created_at', (events.created_at AT TIME ZONE 'UTC'),
    'updated_at', (events.updated_at AT TIME ZONE 'UTC')
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
                SELECT jsonb_agg(`+ClubsOfficerJSONB+`
                        ORDER BY officers.display_order
                    ) AS officers
                FROM officers
                WHERE officers.club_id = club.id
            ) AS officers ON TRUE
            JOIN LATERAL (
                SELECT jsonb_agg(`+ClubsContactJSONB+`
                        ORDER BY contacts.display_order
                    ) AS contacts
                FROM contacts
                WHERE contacts.club_id = club.id
            ) AS contacts ON TRUE
            JOIN LATERAL (
                SELECT jsonb_agg(`+ClubsMembershipFormJSONB+`
                        ORDER BY membership_forms.display_order
                    ) AS membership_forms
                FROM membership_forms
                WHERE membership_forms.club_id = club.id
            ) AS membership_forms ON TRUE
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

	now := time.Now()
	nowDateTime := now.Format(time.DateTime)
	inOneYear := time.Now().Add(365 * 24 * time.Hour)
	inOneYearDateTime := inOneYear.Format(time.DateTime)

	var clubsDatabase *sql.DB = configs.ConnectClubsDB()
	id := c.Param("id")

	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, `
        SELECT jsonb_agg(`+ClubsEventJSONB+`) AS events
        FROM events
        WHERE events.club_id = $1
            AND events.end_time >= '`+nowDateTime+`'
            AND events.start_time <= '`+inOneYearDateTime+`'
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
                SELECT jsonb_agg(`+ClubsOfficerJSONB+`
                        ORDER BY officers.display_order
                    ) AS officers
                FROM officers
                WHERE officers.club_id = club.id
            ) AS officers ON TRUE
            JOIN LATERAL (
                SELECT jsonb_agg(`+ClubsContactJSONB+`
                        ORDER BY contacts.display_order
                    ) AS contacts
                FROM contacts
                WHERE contacts.club_id = club.id
            ) AS contacts ON TRUE
            JOIN LATERAL (
                SELECT jsonb_agg(`+ClubsMembershipFormJSONB+`
                        ORDER BY membership_forms.display_order
                    ) AS membership_forms
                FROM membership_forms
                WHERE membership_forms.club_id = club.id
            ) AS membership_forms ON TRUE
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

	now := time.Now()
	nowDateTime := now.Format(time.DateTime)
	inOneYear := time.Now().Add(365 * 24 * time.Hour)
	inOneYearDateTime := inOneYear.Format(time.DateTime)

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
            AND events.end_time >= '`+nowDateTime+`'
            AND events.start_time <= '`+inOneYearDateTime+`'
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
