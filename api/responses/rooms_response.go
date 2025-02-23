package responses

import "github.com/UTDNebula/nebula-api/api/schema"

type RoomsResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    []schema.BuildingRooms `json:"data"`
}
