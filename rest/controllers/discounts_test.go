package controllers

import (
	"errors"
	"reflect"
	"testing"
)

type stubDistinctResult struct {
	err          error
	decodeErr    error
	results      []any
	decodeCalled bool
}

func (r *stubDistinctResult) Err() error {
	return r.err
}

func (r *stubDistinctResult) Decode(value any) error {
	r.decodeCalled = true
	if r.decodeErr != nil {
		return r.decodeErr
	}
	results := value.(*[]any)
	*results = r.results
	return nil
}

func TestDecodeDiscountCategories(t *testing.T) {
	t.Run("SkipsNonStringValues", func(t *testing.T) {
		result := &stubDistinctResult{results: []any{"food", int32(12), nil, "retail"}}
		want := []string{"food", "retail"}

		got, err := decodeDiscountCategories(result)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Expected %v, got %v", want, got)
		}
	})

	t.Run("ReturnsQueryErrorWithoutDecoding", func(t *testing.T) {
		queryErr := errors.New("distinct query failed")
		result := &stubDistinctResult{err: queryErr}

		_, err := decodeDiscountCategories(result)
		if !errors.Is(err, queryErr) {
			t.Fatalf("Expected query error %v, got %v", queryErr, err)
		}
		if result.decodeCalled {
			t.Fatal("Expected Decode not to be called after a query error")
		}
	})

	t.Run("ReturnsDecodeError", func(t *testing.T) {
		decodeErr := errors.New("distinct decode failed")
		result := &stubDistinctResult{decodeErr: decodeErr}

		_, err := decodeDiscountCategories(result)
		if !errors.Is(err, decodeErr) {
			t.Fatalf("Expected decode error %v, got %v", decodeErr, err)
		}
		if !result.decodeCalled {
			t.Fatal("Expected Decode to be called")
		}
	})
}
