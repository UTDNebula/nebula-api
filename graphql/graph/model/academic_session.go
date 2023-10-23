package model

type AcademicSession struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date" bson:"start_date"`
	EndDate   string `json:"end_date" bson:"end_date"`
}
