package responses

import "github.com/UTDNebula/nebula-api/api/schema"

type MultiProfessorResponse struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    []schema.Professor `json:"data"`
}

type SingleProfessorResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    schema.Professor `json:"data"`
}
