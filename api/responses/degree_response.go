package responses

import "github.com/UTDNebula/nebula-api/api/schema"

type MultiDegreeResponse struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    []schema.Degree `json:"data"`
}

type SingleDegreeResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    schema.Degree `json:"data"`
}
