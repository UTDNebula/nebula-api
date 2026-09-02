package controllers

func stringCategories(results []any) []string {
	categories := make([]string, 0, len(results))
	for _, result := range results {
		category, ok := result.(string)
		if !ok {
			continue // Skip invalid category
		}
		categories = append(categories, category)
	}
	return categories
}
