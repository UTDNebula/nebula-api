package model

type Assistant struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
}
