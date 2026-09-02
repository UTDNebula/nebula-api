package controllers

type distinctResult interface {
	Err() error
	Decode(any) error
}

func decodeDiscountCategories(result distinctResult) ([]string, error) {
	if err := result.Err(); err != nil {
		return nil, err
	}

	var results []any
	if err := result.Decode(&results); err != nil {
		return nil, err
	}

	categories := make([]string, 0, len(results))
	for _, result := range results {
		category, ok := result.(string)
		if !ok {
			continue // Skip invalid category
		}
		categories = append(categories, category)
	}
	return categories, nil
}
