package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/UTDNebula/nebula-api/rest/configs"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func searchOffset(c *gin.Context) (int64, bool) {
	if raw := c.Query("offset"); raw != "" {
		offset, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respond(c, http.StatusBadRequest, "offset is not type integer", err.Error())
			return 0, false
		}
		return offset, true
	}
	return 0, true
}

func relationshipResults[Source, Result any](flag string, c *gin.Context, fetch func(context.Context, bson.M, map[string]bson.D) ([]Result, error), offsetMessage string) ([]Result, bool) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	var query bson.M
	var err error
	if flag == "Search" {
		params := c.Request.URL.Query()
		params.Del("former_offset")
		params.Del("latter_offset")
		query, err = schema.FilterQuery[Source](params)
		if err != nil {
			respond(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
			return nil, false
		}
	} else {
		query, err = getQuery[Source](flag, c)
		if err != nil {
			return nil, false
		}
	}
	pagination, err := configs.GetAggregateLimit(&query, c)
	if err != nil {
		respond(c, http.StatusBadRequest, offsetMessage, err.Error())
		return nil, false
	}
	records, err := fetch(ctx, query, pagination)
	if err != nil {
		respondWithInternalError(c, err)
		return nil, false
	}
	return records, true
}

func respondRelated[Source, Result any](flag string, c *gin.Context, fetch func(context.Context, bson.M, map[string]bson.D) ([]Result, error), offsetMessage string) {
	records, ok := relationshipResults[Source, Result](flag, c, fetch, offsetMessage)
	if ok {
		respond(c, http.StatusOK, "success", records)
	}
}

func gradesByID(c *gin.Context, fetch func(context.Context, bson.ObjectID) ([14]int, error)) {
	id, err := objectIDFromParam(c, "id")
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	grades, err := fetch(ctx, *id)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", grades)
}
