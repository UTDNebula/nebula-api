package responses

import "github.com/UTDNebula/nebula-api/api/schema"

type MultiSectionResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    []schema.Section `json:"data"`
}

type SingleSectionResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    schema.Section `json:"data"`
}
