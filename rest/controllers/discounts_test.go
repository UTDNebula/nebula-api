package controllers

import (
	"reflect"
	"testing"
)

func TestStringCategoriesSkipsNonStringValues(t *testing.T) {
	results := []any{"food", int32(12), nil, "retail"}
	want := []string{"food", "retail"}

	if got := stringCategories(results); !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v, got %v", want, got)
	}
}
