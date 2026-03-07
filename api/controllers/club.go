package controllers

import (
	"context"
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
// @Description	"gets the directory info for given club."
// @Produce		json
// @Param			id	path		string							true	"ID of course to get grades for"
// @Success		200	{object}	schema.APIResponse[schema.Club]	"class object"
// @Failure		500	{object}	schema.APIResponse[string]		"A string describing the error"
// @Failure		400	{object}	schema.APIResponse[string]		"A string describing the error"
func ClubDirectoryInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := c.Param("id")
	db := configs.ConnectClubsDB()

	var raw []byte
	err := db.QueryRowContext(ctx, `
        SELECT
          jsonb_agg(
          jsonb_build_object(
          'slug',slug,
          'id', club.id,
          'name',club.name,
          'tags',tags,
          'profile_image', profile_image,
          'updated_at', (updated_at AT TIME ZONE 'UTC'),
          'officers', officers,
          'contacts', contacts
          ))
         FROM club
         JOIN LATERAL (
          SELECT jsonb_agg(
              jsonb_build_object(
                'platform',contacts.platform,
                'url', contacts.url
              ) ORDER BY contacts.platform) as contacts from contacts where contacts.club_id = club.id
          )as contacts on TRUE
         JOIN LATERAL (
            SELECT jsonb_agg(jsonb_build_object('name', officers.name, 'position', officers.position)) as officers FROM officers where officers.club_id = club.id
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
