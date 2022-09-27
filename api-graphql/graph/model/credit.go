package model

type Credit struct {
	Category    string `json:"category"`
	CreditHours int    `json:"credit_hours"`
}

func (Credit) IsOutcome() {}
