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

// @Id				clubGet
// @Router			/club/{id} [get]
// @Tags			Clubs
// @Description	"Returns the directory info for given club."
// @Produce		json
// @Param			id	path		string							true	"ID of the club to get"
// @Success		200	{object}	schema.APIResponse[schema.Club]	"A club"
// @Failure		500	{object}	schema.APIResponse[string]		"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]		"A string describing the error"
func ClubDirectoryInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var clubsDatabase *sql.DB = configs.GetClubsDB()
	id := c.Param("id")

	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, `
    SELECT
        jsonb_agg(jsonb_build_object(
            'slug', slug,
            'id', club.id,
            'name', club.name,
            'description', club.description,
            'tags', tags,
            'profile_image', profile_image,
            'updated_at', (updated_at AT TIME ZONE 'UTC'),
            'officers', officers,
            'contacts', contacts
        ))
    FROM club
    JOIN LATERAL (
        SELECT jsonb_agg(jsonb_build_object(
            'platform', contacts.platform,
            'url', contacts.url
        ) ORDER BY contacts.platform) as contacts from contacts where contacts.club_id = club.id
    ) as contacts on TRUE
    JOIN LATERAL (
        SELECT jsonb_agg(jsonb_build_object(
            'name', officers.name, 
            'position', officers.position
        )) as officers FROM officers where officers.club_id = club.id
    ) as officers on TRUE
    WHERE club.id = $1;
  `, id).Scan(&raw)

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

	respond(c, http.StatusOK, "success", clubs[0])
}

// @Id				clubSearch
// @Router			/club/search [get]
// @Tags			Clubs
// @Description	"Returns list of clubs matching the search string"
// @Produce		json
// @Param			q	query		string								true	"Search string"
// @Success		200	{object}	schema.APIResponse[[]schema.Club]	"List of matching clubs"
// @Failure		500	{object}	schema.APIResponse[string]			"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]			"A string describing the error"
func ClubSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var clubsDatabase *sql.DB = configs.GetClubsDB()
	search := c.Query("q")

	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, `
	SELECT
        jsonb_agg(jsonb_build_object(
            'slug',slug,
            'id', club.id,
            'name',club.name,
            'description', club.description,
            'tags',tags,
            'profile_image', profile_image,
            'updated_at', (updated_at AT TIME ZONE 'UTC'),
            'officers', officers,
            'contacts', contacts
        ) ORDER BY (
            (20.0 * word_similarity($1, coalesce(club.alias, '')))
            + (10.0 * word_similarity($1, club.name))
            + (5.0 * word_similarity($1, coalesce(array_to_string(club.tags, ' '), '')))
            + (coalesce(-1.0 * (club.search_tsv <@> to_bm25query(to_tsvector('english', $1), 'club_search_idx')), 0.0))
        ) DESC) as club
    FROM club
    JOIN LATERAL (
        SELECT jsonb_agg(jsonb_build_object(
            'platform',contacts.platform,
            'url', contacts.url
        ) ORDER BY contacts.platform) as contacts from contacts where contacts.club_id = club.id
    ) as contacts on TRUE
    JOIN LATERAL (
        SELECT jsonb_agg(jsonb_build_object('name',officers.name, 'position',officers.position)) as officers FROM officers where officers.club_id = club.id
    ) as officers on TRUE
    WHERE (
        club.search_tsv @@ websearch_to_tsquery('english', $1)
        OR word_similarity($1, club.name) >= 0.2
        OR word_similarity($1, COALESCE(club.alias, '')) >= 0.2
        OR club.name ILIKE '%' || $1 || '%'
        OR club.alias ILIKE '%' || $1 || '%'
    ) AND (
        'approved' = 'approved'::approved_enum
    )
  `, search).Scan(&raw)

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

	// Since filtering by ID, return the single club
	respond(c, http.StatusOK, "success", clubs)
}
