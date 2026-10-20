package responses

import "github.com/UTDNebula/nebula-api/rest/schema"

type BucketResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    schema.BucketInfo `json:"data"`
}

type ObjectResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    schema.ObjectInfo `json:"data"`
}

type DeleteResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    int    `json:"data"`
}
