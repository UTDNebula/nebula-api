package responses

import "github.com/UTDNebula/nebula-api/api/schema"

type MultiCourseResponse struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    []schema.Course `json:"data"`
}

type SingleCourseResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    schema.Course `json:"data"`
}
