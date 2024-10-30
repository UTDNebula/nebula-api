package responses

import "github.com/UTDNebula/nebula-api/api/schema"

type AutocompleteResponse struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    []schema.Autocomplete `json:"data"`
}
