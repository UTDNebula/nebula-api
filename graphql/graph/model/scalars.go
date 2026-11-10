package model

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// This file is for custom marshaller and unmarshaller for types not from gqlgen

// Marshaller for ObjectID scalar type
func MarshalObjectID(id bson.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%q", id.Hex()))
	})
}

// Unmarshaller for ObjectID scalar type
func UnmarshalObjectID(v any) (bson.ObjectID, error) {
	str, ok := v.(string)
	if !ok {
		return bson.NilObjectID, fmt.Errorf("ObjectID must be a string")
	}
	objectID, err := bson.ObjectIDFromHex(str)
	return objectID, err
}
