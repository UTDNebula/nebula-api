package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/UTDNebula/nebula-api/rest/configs"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
)

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
        SELECT jsonb_agg(`+schema.ClubsClubJSONB+`)
        FROM club
            JOIN LATERAL (
                SELECT jsonb_agg(`+schema.ClubsOfficerJSONB+`
                        ORDER BY officers.display_order
                    ) AS officers
                FROM officers
                WHERE officers.club_id = club.id
            ) AS officers ON TRUE
            JOIN LATERAL (
                SELECT jsonb_agg(`+schema.ClubsContactJSONB+`
                        ORDER BY contacts.display_order
                    ) AS contacts
                FROM contacts
                WHERE contacts.club_id = club.id
            ) AS contacts ON TRUE
            JOIN LATERAL (
                SELECT jsonb_agg(`+schema.ClubsMembershipFormJSONB+`
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

	var clubExists bool
	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, `
        SELECT
            EXISTS (SELECT 1 FROM club WHERE club.id = $1) AS club_exists,
            (
                SELECT jsonb_agg(`+schema.ClubsEventJSONB+`) AS events
                FROM events
                WHERE events.club_id = $1
                    AND events.end_time >= '`+nowDateTime+`'
                    AND events.start_time <= '`+inOneYearDateTime+`'
                    AND events.approved = 'approved'::status_enum
            ) AS events;`,
		id).Scan(&clubExists, &raw)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if !clubExists {
		respond(c, http.StatusNotFound, "error", "Club not found")
		return
	}
	if raw == nil {
		respond(c, http.StatusNotFound, "error", "No events found for club")
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
        SELECT jsonb_agg(`+schema.ClubsClubJSONB+`
                ORDER BY (
                    (20.0 * word_similarity($1, coalesce(club.alias, '')))
                    + (10.0 * word_similarity($1, club.name))
                    + (5.0 * word_similarity($1, coalesce(array_to_string(club.tags, ' '), '')))
                    + (coalesce(-1.0 * (club.search_tsv <@> to_bm25query(to_tsvector('english', $1), 'club_search_idx')), 0.0))
                ) DESC
            ) AS club
        FROM club
            JOIN LATERAL (
                SELECT jsonb_agg(`+schema.ClubsOfficerJSONB+`
                        ORDER BY officers.display_order
                    ) AS officers
                FROM officers
                WHERE officers.club_id = club.id
            ) AS officers ON TRUE
            JOIN LATERAL (
                SELECT jsonb_agg(`+schema.ClubsContactJSONB+`
                        ORDER BY contacts.display_order
                    ) AS contacts
                FROM contacts
                WHERE contacts.club_id = club.id
            ) AS contacts ON TRUE
            JOIN LATERAL (
                SELECT jsonb_agg(`+schema.ClubsMembershipFormJSONB+`
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
		respond(c, http.StatusNotFound, "error", "No clubs found")
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
        SELECT jsonb_agg(`+schema.ClubsEventJSONB+`) AS events
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
// @Description	Returns list of events matching the search string and filters
// @Produce		json
// @Param			q			query		string									false	"Search string to match event names and descriptions."
// @Param			sort		query		string									false	"Order of returned events. Default to upcoming. If `q` is provided, then defaults to relevant."	Enums(relevant, upcoming, past, updated)
// @Param			start_date	query		string									false	"Start date to begin searching for events. Defaults to current date. If `sort` is `past`, then defaults to no value.\nFormat: YYYY-MM-DD"
// @Param			end_date	query		string									false	"End date to search for events up to. Defaults to one year in the future. If `sort` is `past`, then defaults to current date.\nFormat: YYYY-MM-DD"
// @Param			limit		query		int										false	"Number of results to return."							minimum(1)	maximum(100)	default(20)
// @Param			offset		query		int										false	"The starting position of the current page of events"	minimum(0)	default(0)
// @Param			tag			query		[]string								false	"A club tag that should exist in the returned events' host clubs."
// @Success		200			{object}	schema.APIResponse[[]schema.ClubsEvent]	"List of matching events"
// @Failure		500			{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400			{object}	schema.APIResponse[string]				"A string describing the error"
func ClubsEventSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	now := time.Now()
	nowDateTime := now.Format(time.DateTime)
	inOneYear := time.Now().Add(365 * 24 * time.Hour)
	inOneYearDateTime := inOneYear.Format(time.DateTime)

	var params schema.ClubsEventQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		respond(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}
	params.TrimSpace()

	var clubsDatabase *sql.DB = configs.ConnectClubsDB()

	var args []any
	argc := 0
	var conditions []string
	join := ""
	orderBy := ""
	limit := "LIMIT 20"
	offset := ""

	// Handle Query
	if params.Query != "" {
		args = append(args, params.Query)
		argc++
		conditions = append(conditions, fmt.Sprint(`events.search_tsv @@ (websearch_to_tsquery('english', $`, argc, `))`))
		// By default, order by relevance if there's a search query
		orderBy = fmt.Sprint(`ORDER BY (ts_rank(events.search_tsv, websearch_to_tsquery('english', $`, argc, `))) DESC`)
	}

	// Handle Sort
	// By default, sort by upcoming. If Query is present, sort by relevant. If Sort is present, sort according to that
	switch params.Sort {
	case "relevant":
		// Use orderBy statement from Query handler
	case "past":
		orderBy = "ORDER BY events.start_time DESC"
	case "updated":
		orderBy = "ORDER BY events.updated_at DESC"
	default:
		if params.Query != "" {
			break
		}
		fallthrough
	case "upcoming":
		orderBy = "ORDER BY events.start_time ASC"
	}

	// Handle Tags
	if len(params.Tags) != 0 {
		args = append(args, params.Tags)
		argc++
		join = "JOIN club ON events.club_id = club.id"
		conditions = append(conditions, fmt.Sprint(`club.tags && $`, argc))
	}

	// Handle StartDate
	if params.StartDate.IsZero() {
		if params.Sort != "past" {
			conditions = append(conditions, `events.end_time >= '`+nowDateTime+`'`)
		}
	} else {
		args = append(args, params.StartDate)
		argc++
		conditions = append(conditions, fmt.Sprint(`events.end_time >= $`, argc, ``))
	}

	// Handle EndDate
	if params.EndDate.IsZero() {
		if params.Sort != "past" {
			conditions = append(conditions, `events.start_time <= '`+inOneYearDateTime+`'`)
		} else {
			conditions = append(conditions, `events.start_time <= '`+nowDateTime+`'`)
		}
	} else {
		args = append(args, params.EndDate)
		argc++
		conditions = append(conditions, fmt.Sprint(`events.start_time <= $`, argc, ``))
	}

	// Handle Limit
	if params.Limit > 0 {
		args = append(args, params.Limit)
		argc++
		limit = fmt.Sprint("LIMIT $", argc)
	}

	// Handle Offset
	if params.Offset > 0 {
		args = append(args, params.Offset)
		argc++
		offset = fmt.Sprint("OFFSET $", argc)
	}

	// Construct query
	var conditionsBuilder strings.Builder
	for _, condition := range conditions {
		conditionsBuilder.WriteString(" AND ")
		conditionsBuilder.WriteString(condition)
	}
	query := `
        SELECT jsonb_agg(subquery.event_data) AS events
        FROM (
            SELECT ` + schema.ClubsEventJSONB + ` AS event_data
            FROM events
            ` + join + `
            WHERE events.approved = 'approved'::status_enum` + conditionsBuilder.String() + `
            ` + orderBy + `
            ` + limit + `
            ` + offset + `
        ) AS subquery;`

	// // DEBUGGING
	// fmt.Println("QUERY:", query)
	// fmt.Println()
	// fmt.Println("  ", argc, "ARGS:", args)
	// fmt.Println(" ORDER BY:", orderBy)
	// fmt.Println("    LIMIT:", limit)
	// fmt.Println("   OFFSET:", offset)
	// fmt.Println("StartDate:", params.StartDate)
	// fmt.Println("  EndDate:", params.EndDate)

	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, query, args...).Scan(&raw)

	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if raw == nil {
		respond(c, http.StatusNotFound, "error", "No events found")
		return
	}

	var clubsEvents []schema.ClubsEvent
	if err := json.Unmarshal(raw, &clubsEvents); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", clubsEvents)
}
