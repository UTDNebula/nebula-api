package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"encoding/json"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/schema"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var clubsDatabase *sql.DB = configs.ConnectClubsDB()
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

	// Since filtering by ID, return the single club
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
      ) ORDER BY paradedb.score(id) DESC) as club
    FROM club
    JOIN LATERAL (
      SELECT jsonb_agg(jsonb_build_object(
        'platform',contacts.platform,
        'url', contacts.url
      ) ORDER BY contacts.platform) as contacts from contacts where contacts.club_id = club.id
    ) as contacts on TRUE
    JOIN LATERAL (
      SELECT jsonb_agg(jsonb_build_object('name',officers.name, 'position',officers.position)) as officers FROM officers where officers.club_id = club.id
    ) as officers on TRUE where id @@@
      paradedb.boolean(
        should => ARRAY[
          paradedb.boost(20,paradedb.match('alias',$1,distance=>2)),
          paradedb.boost(10,paradedb.match('name',$1,distance=>2)),
          paradedb.boost(1,paradedb.match('description',$1,distance=>1)),
          paradedb.boost(5,paradedb.match('tags',$1,distance=>1))
        ]) and id @@@ 
      paradedb.const_score(0.0, paradedb.term('approved','approved'::approved_enum));
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
