package model

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProfessorObjectIDReferences(t *testing.T) {
	professorID, sectionID := primitive.NewObjectID(), primitive.NewObjectID()
	data, err := bson.Marshal(bson.M{
		"_id": professorID, "first_name": "Ada", "sections": bson.A{sectionID},
	})
	if err != nil {
		t.Fatal(err)
	}
	var dbProfessor DBProfessor
	if err := bson.Unmarshal(data, &dbProfessor); err != nil {
		t.Fatal(err)
	}
	professor := TransformProfessor(&dbProfessor)
	if professor.ID != professorID.Hex() || professor.FirstName != "Ada" || len(professor.Sections) != 1 || professor.Sections[0].ID != sectionID.Hex() {
		t.Fatalf("unexpected transformed professor: %+v", professor)
	}

	data, err = bson.Marshal(bson.M{"_id": sectionID, "professors": bson.A{professorID}})
	if err != nil {
		t.Fatal(err)
	}
	var dbSection DBSection
	if err := bson.Unmarshal(data, &dbSection); err != nil {
		t.Fatal(err)
	}
	section := TransformSection(&dbSection)
	if section.ID != sectionID.Hex() || len(section.Professors) != 1 || section.Professors[0].ID != professor.ID {
		t.Fatalf("unexpected transformed section: %+v", section)
	}
}
