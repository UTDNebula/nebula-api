package controllers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"encoding/json"

	"github.com/UTDNebula/nebula-api/rest/configs"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

var clubTracer = otel.Tracer("club-controller")

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

	ctx, span := clubTracer.Start(ctx, "club.directory.info")
	defer span.End()

	var clubsDatabase *sql.DB = configs.ConnectClubsDB()
	id := c.Param("id")

	var raw []byte
	err := clubsDatabase.QueryRowContext(ctx, `
    SELECT
        jsonb_build_object(
            'slug', slug,
            'id', club.id,
            'name', club.name,
            'description', club.description,
            'tags', tags,
            'profile_image', profile_image,
            'updated_at', (updated_at AT TIME ZONE 'UTC'),
            'officers', officers,
            'contacts', contacts
        ) AS club
    FROM club
    JOIN LATERAL (
        SELECT jsonb_agg(jsonb_build_object(
            'platform', contacts.platform,
            'url', contacts.url
        ) ORDER BY contacts.platform) AS contacts FROM contacts WHERE contacts.club_id = club.id
    ) AS contacts ON TRUE
    JOIN LATERAL (
        SELECT jsonb_agg(jsonb_build_object(
            'name', officers.name,
            'position', officers.position
        )) AS officers FROM officers WHERE officers.club_id = club.id
    ) AS officers ON TRUE
    WHERE club.id = $1;
    `, id).Scan(&raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond(c, http.StatusNotFound, "error", "Club not found")
			return
		}
		respondWithInternalError(c, err)
		return
	}

	var club schema.Club
	if err := json.Unmarshal(raw, &club); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", club)
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

	ctx, span := clubTracer.Start(ctx, "club.search")
	defer span.End()

	var clubsDatabase *sql.DB = configs.ConnectClubsDB()
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
        ) ORDER BY paradedb.score(id) DESC) AS club
    FROM club
    JOIN LATERAL (
        SELECT jsonb_agg(jsonb_build_object(
            'platform',contacts.platform,
            'url', contacts.url
        ) ORDER BY contacts.platform) AS contacts FROM contacts WHERE contacts.club_id = club.id
    ) AS contacts ON TRUE
    JOIN LATERAL (
        SELECT jsonb_agg(jsonb_build_object(
            'name',officers.name, 
            'position',officers.position
        )) AS officers FROM officers WHERE officers.club_id = club.id
    ) AS officers ON TRUE WHERE id @@@
        paradedb.boolean(
            should => ARRAY[
            paradedb.boost(20,paradedb.match('alias',$1,distance=>2)),
            paradedb.boost(10,paradedb.match('name',$1,distance=>2)),
            paradedb.boost(1,paradedb.match('description',$1,distance=>1)),
            paradedb.boost(5,paradedb.match('tags',$1,distance=>1))
            ]) AND id @@@ 
        paradedb.const_score(0.0, paradedb.term('approved','approved'::approved_enum));
    `, search).Scan(&raw)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	if raw == nil {
		respond(c, http.StatusNotFound, "error", "Clubs not found")
		return
	}

	var clubs []schema.Club
	if err := json.Unmarshal(raw, &clubs); err != nil {
		respondWithInternalError(c, err)
		return
	}

	respond(c, http.StatusOK, "success", clubs)
}
