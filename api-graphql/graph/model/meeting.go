package model

import (
	"fmt"
	"io"
	"strconv"
)

type Meeting struct {
	StartDate   string       `json:"start_date"`
	EndDate     string       `json:"end_date"`
	MeetingDays []string     `json:"meeting_days"`
	StartTime   string       `json:"start_time"`
	EndTime     string       `json:"end_time"`
	Modality    ModalityType `json:"modality"`
	Location    *Location    `json:"location"`
}

type ModalityType string

const (
	ModalityTypePending     ModalityType = "PENDING"
	ModalityTypeTraditional ModalityType = "TRADITIONAL"
	ModalityTypeHybrid      ModalityType = "HYBRID"
	ModalityTypeFlexible    ModalityType = "FLEXIBLE"
	ModalityTypeRemote      ModalityType = "REMOTE"
	ModalityTypeOnline      ModalityType = "ONLINE"
)

var AllModalityType = []ModalityType{
	ModalityTypePending,
	ModalityTypeTraditional,
	ModalityTypeHybrid,
	ModalityTypeFlexible,
	ModalityTypeRemote,
	ModalityTypeOnline,
}

func (e ModalityType) IsValid() bool {
	switch e {
	case ModalityTypePending, ModalityTypeTraditional, ModalityTypeHybrid, ModalityTypeFlexible, ModalityTypeRemote, ModalityTypeOnline:
		return true
	}
	return false
}

func (e ModalityType) String() string {
	return string(e)
}

func (e *ModalityType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ModalityType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ModalityType", str)
	}
	return nil
}

func (e ModalityType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
