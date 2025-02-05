package responses

import "github.com/UTDNebula/nebula-api/api/schema"

type MultiBuildingEventsResponse struct {
	Status  int                        `json:"status"`
	Message string                     `json:"message"`
	Data    schema.MultiBuildingEvents `json:"data"`
}

type SingleBuildingEventsResponse struct {
	Status  int                         `json:"status"`
	Message string                      `json:"message"`
	Data    schema.SingleBuildingEvents `json:"data"`
}
