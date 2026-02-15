package schema

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

var queryableCache sync.Map

// FilterQuery converts URL query parameters into a MongoDB BSON query filter.
// It validates that each query parameter corresponds to a field in type F that is
// marked as queryable. Additionally, parameter offest can be included, referring to
// mongo query pagination.
//
// Notes:
//   - Currently filter only supports exact matching and treats fields as strings
//   - More "meta" filters may need to be added in the future.
//
// Returns an error if:
//   - A query parameter key is not defined in the struct
//   - A field exists but is not marked as queryable
//   - Multiple values are provided for a single parameter
//   - The "offset" parameter cannot be parsed as an integer
func FilterQuery[F any](urlValues url.Values) (bson.M, error) {
	queryable, err := loadQueryable(reflect.TypeFor[F]())
	if err != nil {
		return nil, err
	}

	query := bson.M{}
	for key, values := range urlValues {
		if key == "offset" {
			if num, err := strconv.Atoi(values[0]); err == nil {
				query[key] = num
				continue
			} else {
				return nil, fmt.Errorf("offest must be an integer")
			}
		}

		if len(values) > 1 {
			return nil, fmt.Errorf("multi-value queries are not supported")
		}

		allowed, exists := queryable[key]
		if !exists {
			return nil, fmt.Errorf("unknown query parameter '%s'", key)
		}
		if !allowed {
			return nil, fmt.Errorf("field '%s' cannot be used for filtering", key)
		}
		query[key] = values[0]
	}

	return query, nil
}

// loadQueryable returns a map indicating which fields of the given type are queryable.
func loadQueryable(t reflect.Type) (map[string]bool, error) {
	if cached, ok := queryableCache.Load(t.String()); ok {
		//should literally never fail but its best practice to check casts
		if queryMap, ok := cached.(map[string]bool); ok {
			return queryMap, nil
		}
		queryableCache.Delete(t.String())
		return nil, fmt.Errorf("queryableCache was corrupted: %s was not of type map[string]bool", t.String())
	}

	queryable := make(map[string]bool)
	err := recBuild(t, "", queryable, make([]reflect.Type, 0))
	if err != nil {
		return nil, err

	}
	actual, _ := queryableCache.LoadOrStore(t.String(), queryable)
	if queryMap, ok := actual.(map[string]bool); ok {
		return queryMap, nil
	}
	queryableCache.Delete(t.String())
	return nil, fmt.Errorf("queryableCache was corrupted: %s was not of type map[string]bool", t.String())
}

// recBuild recursively traverses a struct type to build a map of queryable fields.
// It constructs dot-notation paths for nested fields and determines whether each field
// can be used for filtering based on the "queryable" tag.
func recBuild(t reflect.Type, prefix string, queryableMap map[string]bool, visited []reflect.Type) error {
	if willCreateLoop(visited, t) {
		return nil
	}

	newVisited := append(visited, t)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		json, hasJson := field.Tag.Lookup("json")
		if !hasJson {
			return fmt.Errorf("exported field '%s.%s' missing json tag (use json:\"-\" to exclude)", t.Name(), field.Name)
		} else if json == "-" {
			continue
		}

		// Determine the JSON path
		fullPath := strings.Split(json, ",")[0]
		if prefix != "" {
			fullPath = prefix + "." + fullPath
		}

		fieldType := field.Type
		for fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		_, queryable := field.Tag.Lookup("queryable")
		if fieldType.Kind() == reflect.Struct {
			if queryable {
				if err := recBuild(field.Type, fullPath, queryableMap, newVisited); err != nil {
					return err
				}
			} else {
				queryableMap[fullPath] = false
			}
		} else {
			queryableMap[fullPath] = queryable
		}
	}
	return nil
}

// willCreateLoop determines if adding `value` to the `visited` list would create
// a loop.
func willCreateLoop(visited []reflect.Type, value reflect.Type) bool {
	if value.Kind() != reflect.Ptr {
		return false
	}

	index := -1
	for i := len(visited) - 1; i >= 0; i-- {
		if visited[i] == value {
			index = i
			break
		}
	}

	return index > 0 && visited[index-1] == visited[len(visited)-1]
}
