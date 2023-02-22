package model

type AcademicSession struct {
	Name      string `bson:"name"`
	StartDate string `bson:"start_date"`
	EndDate   string `bson:"end_date"`
}
