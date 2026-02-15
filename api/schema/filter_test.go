package schema

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
)

type _normal struct {
	Name   string `bson:"name" json:"name" queryable:""`
	Number int    `bson:"number" json:"number" queryable:""`
	Hidden bool   `bson:"hidden" json:"hidden"`
}

type _missingJson struct {
	Name   string
	Number int
	Hidden bool
}

type _missingQueryable struct {
	Name   string `bson:"name" json:"name"`
	Number int    `bson:"number" json:"number"`
	Hidden bool   `bson:"hidden" json:"hidden"`
}

type _nested struct {
	Name   string  `bson:"name" json:"name" queryable:""`
	Number int     `bson:"number" json:"number" queryable:""`
	Hidden bool    `bson:"hidden" json:"hidden"`
	Nested _normal `bson:"nested" json:"nested" queryable:""`
}

type _nestedAnonymous struct {
	Name   string `bson:"name" json:"name" queryable:""`
	Number int    `bson:"number" json:"number" queryable:""`
	Hidden bool   `bson:"hidden" json:"hidden"`
	Nested struct {
		Normal       _normal `bson:"normal" json:"normal" queryable:""`
		NormalHidden _normal `bson:"nested_hidden" json:"nested_hidden"`
	} `bson:"nested" json:"nested" queryable:""`
	NestedHidden struct {
	} `bson:"nested" json:"nested_hidden"`
}

type _nestedEmbedded struct {
	_normal `bson:",inline" json:",inline" queryable:""`
}

type _nestedPointer struct {
	Name   string   `bson:"name" json:"name" queryable:""`
	Number int      `bson:"number" json:"number" queryable:""`
	Hidden bool     `bson:"hidden" json:"hidden"`
	Nested *_normal `bson:"nested" json:"nested" queryable:""`
}

type _nestedDoublePointer struct {
	Name   string          `bson:"name" json:"name" queryable:""`
	Number int             `bson:"number" json:"number" queryable:""`
	Hidden bool            `bson:"hidden" json:"hidden"`
	Nested *_nestedPointer `bson:"nested" json:"nested" queryable:""`
}

type _nestedInfinite struct {
	Name   string           `bson:"name" json:"name" queryable:""`
	Number int              `bson:"number" json:"number" queryable:""`
	Hidden bool             `bson:"hidden" json:"hidden"`
	Nested *_nestedInfinite `bson:"nested" json:"nested" queryable:""`
}

type _slice struct {
	Names []string `bson:"names" json:"names" queryable:""`
}

type _array struct {
	Names [12]string `bson:"names" json:"names" queryable:""`
}

type _slicePointer struct {
	Names []*string `bson:"names" json:"names" queryable:""`
}

type _sliceDoublePointer struct {
	Names *[]*string `bson:"names" json:"names" queryable:""`
}

type _jsonExcluded struct {
	Name   string `bson:"name" json:"name" queryable:""`
	Number int    `bson:"-" json:"-"`
	Hidden bool   `bson:"hidden" json:"hidden"`
}

func TestFilterQuery(t *testing.T) {

	testCases := map[string]struct {
		Function func(values url.Values) (bson.M, error)
		UrlQuery map[string][]string
		Fail     bool
		Expected bson.M
	}{
		"Normal": {
			Function: FilterQuery[_normal],
			UrlQuery: map[string][]string{
				"name":   {"bob"},
				"number": {"0"},
			},
			Expected: bson.M{
				"name":   "bob",
				"number": "0",
			},
		},
		"Nested": {
			Function: FilterQuery[_nested],
			UrlQuery: map[string][]string{
				"name":          {"bob"},
				"number":        {"0"},
				"nested.name":   {"bob"},
				"nested.number": {"0"},
			},
			Expected: bson.M{
				"name":          "bob",
				"number":        "0",
				"nested.name":   "bob",
				"nested.number": "0",
			},
		},
		"Normal with Offest": {
			Function: FilterQuery[_normal],
			UrlQuery: map[string][]string{
				"name":   {"bob"},
				"offset": {"0"},
			},
			Expected: bson.M{
				"name":   "bob",
				"offset": 0,
			},
		},
		"Fail empty parameter": {
			Function: FilterQuery[_nested],
			UrlQuery: map[string][]string{
				"": {"false"},
			},
			Fail: true,
		},
		"Fail field cannot be queried": {
			Function: FilterQuery[_nested],
			UrlQuery: map[string][]string{
				"hidden": {"false"},
			},
			Fail: true,
		},
		"Fail multiple values": {
			Function: FilterQuery[_nested],
			UrlQuery: map[string][]string{
				"": {"false", "true"},
			},
			Fail: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := tc.Function(tc.UrlQuery)

			if tc.Fail {
				if err == nil {
					t.Fatal("expected error but got nil")
				}

			} else {
				if err != nil {
					t.Errorf("unexpected error %v ", err)
				}

				if diff := cmp.Diff(tc.Expected, result); diff != "" {
					t.Errorf("Failed (-expected +got)\n %s", diff)
				}
			}

			if tc.Fail && err == nil {
				t.Errorf("expected error, got nil")
			} else if !tc.Fail && err != nil {
				t.Errorf("unexpected error %v ", err)
			}
		})
	}

}

func TestLoadQueryable(t *testing.T) {

	testcases := map[string]struct {
		Type     reflect.Type
		Expected map[string]bool
		Fail     bool
	}{
		"Normal": {
			Type: reflect.TypeFor[_normal](),
			Expected: map[string]bool{
				"name":   true,
				"number": true,
				"hidden": false,
			},
		},
		"Missing Json": {
			Type: reflect.TypeFor[_missingJson](),
			Fail: true,
		},
		"Missing Queryable": {
			Type: reflect.TypeFor[_missingQueryable](),
			Expected: map[string]bool{
				"name":   false,
				"number": false,
				"hidden": false,
			},
		},
		"Nested": {
			Type: reflect.TypeFor[_nested](),
			Expected: map[string]bool{
				"name":          true,
				"number":        true,
				"hidden":        false,
				"nested.name":   true,
				"nested.number": true,
				"nested.hidden": false,
			},
		},
		"Nested Embedded": {
			Type: reflect.TypeFor[_nestedEmbedded](),
			Expected: map[string]bool{
				"name":   true,
				"number": true,
				"hidden": false,
			},
		},
		"Nested Anonymous": {
			Type: reflect.TypeFor[_nestedAnonymous](),
			Expected: map[string]bool{
				"name":                 true,
				"number":               true,
				"hidden":               false,
				"nested_hidden":        false,
				"nested.normal.name":   true,
				"nested.normal.number": true,
				"nested.normal.hidden": false,
				"nested.nested_hidden": false,
			},
		},
		"Nested Pointer": {
			Type: reflect.TypeFor[_nestedPointer](),
			Expected: map[string]bool{
				"name":          true,
				"number":        true,
				"hidden":        false,
				"nested.name":   true,
				"nested.number": true,
				"nested.hidden": false,
			},
		},
		"Nested Double Pointer": {
			Type: reflect.TypeFor[_nestedDoublePointer](),
			Expected: map[string]bool{
				"name":                 true,
				"number":               true,
				"hidden":               false,
				"nested.name":          true,
				"nested.number":        true,
				"nested.hidden":        false,
				"nested.nested.name":   true,
				"nested.nested.number": true,
				"nested.nested.hidden": false,
			},
		},
		"Slice": {
			Type: reflect.TypeFor[_slice](),
			Expected: map[string]bool{
				"names": true,
			},
		},
		"Array": {
			Type: reflect.TypeFor[_array](),
			Expected: map[string]bool{
				"names": true,
			},
		},
		"Slice Pointer": {
			Type: reflect.TypeFor[_slicePointer](),
			Expected: map[string]bool{
				"names": true,
			},
		},
		"Slice Double Pointer": {
			Type: reflect.TypeFor[_sliceDoublePointer](),
			Expected: map[string]bool{
				"names": true,
			},
		},
		"Nested Infinite": {
			Type: reflect.TypeFor[_nestedInfinite](),
			Expected: map[string]bool{
				"name":                 true,
				"number":               true,
				"hidden":               false,
				"nested.name":          true,
				"nested.number":        true,
				"nested.hidden":        false,
				"nested.nested.name":   true,
				"nested.nested.number": true,
				"nested.nested.hidden": false,
			},
		},
		"Json Excluded": {
			Type: reflect.TypeFor[_jsonExcluded](),
			Expected: map[string]bool{
				"name":   true,
				"hidden": false,
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			result, err := loadQueryable(tc.Type)

			if (err != nil) != tc.Fail {
				t.Errorf("loadQueryable() error = %v, fail %v", err, tc.Fail)
				return
			}

			if diff := cmp.Diff(tc.Expected, result); diff != "" {
				t.Errorf("Failed (-expected +got)\n %s", diff)
			}
		})
	}

	t.Run("Cache Corruption", func(t *testing.T) {
		rType := reflect.TypeFor[_normal]()
		typeName := rType.String()

		t.Run("Corrupted on Load", func(t *testing.T) {
			queryableCache.Store(typeName, 14)

			_, err := loadQueryable(rType)
			if err == nil {
				t.Fatal("expected error when cache contains wrong type")
			}

			if _, exists := queryableCache.Load(typeName); exists {
				t.Error("corrupted cache entry should have been deleted")
			}

			if _, err = loadQueryable(rType); err != nil {
				t.Fatalf("unexpected failure after cache cleared: %v", err)
			}
		})

		t.Run("Recovery After Corruption", func(t *testing.T) {
			queryableCache.Store(typeName, "wrong type")

			if _, err := loadQueryable(reflect.TypeFor[_normal]()); err == nil {
				t.Fatal("expected corruption error")
			}

			// first will load, second will be cached
			for range 2 {
				if _, err := loadQueryable(reflect.TypeFor[_normal]()); err != nil {
					t.Fatalf("should recover after corruption: %v", err)
				}
			}

		})
	})

}
