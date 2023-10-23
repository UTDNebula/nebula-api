package model

type Location struct {
	Building string `json:"building"`
	Room     string `json:"room"`
	MapURI   string `json:"map_uri" bson:"map_uri"`
}
