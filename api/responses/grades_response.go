package responses

type GradeResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SectionGradeResponse struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	GradeData []GradeData `json:"grade_data"`
}

type GradeData struct {
	Id   string `bson:"_id" json:"_id"`
	Data []struct {
		Type              string      `bson:"type" json:"type"`
		GradeDistribution interface{} `bson:"grade_distribution" json:"grade_distribution"`
	} `bson:"data" json:"data"`
}
