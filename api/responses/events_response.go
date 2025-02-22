package responses

import "github.com/UTDNebula/nebula-api/api/schema"

type MultiBuildingEventsResponse[T any] struct {
	Status  int                           `json:"status"`
	Message string                        `json:"message"`
	Data    schema.MultiBuildingEvents[T] `json:"data"`
}

type SingleBuildingEventsResponse[T any] struct {
	Status  int                            `json:"status"`
	Message string                         `json:"message"`
	Data    schema.SingleBuildingEvents[T] `json:"data"`
}
