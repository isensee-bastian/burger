package burger

// revenue calculates the monetary value of a sold burger by comparing its ingredients with the customers expectations,
// i.e. the ordered burger. The closer the actual burger is to the customers expectations, the higher the revenue
// (score) will be.
func revenue(expectedIngredients, actualIngredients []IngredientType) int {
	expectedLookup := map[IngredientType]struct{}{}
	actualLookup := map[IngredientType]struct{}{}

	// Create map lookups for simpler and more efficient checking plus deduplication.
	for _, expected := range expectedIngredients {
		expectedLookup[expected] = struct{}{}
	}
	for _, actual := range actualIngredients {
		actualLookup[actual] = struct{}{}
	}

	revenue := 0

	// Add 1 for each ingredient that matches the customers expectation (order does not matter except for buns below).
	for expected := range expectedLookup {
		_, found := actualLookup[expected]

		if found {
			revenue += 1
		}
	}
	// Subtract 1 for each ingredient that the customer does not expect.
	for actual := range actualLookup {
		_, found := expectedLookup[actual]

		if !found {
			revenue -= 1
		}
	}

	// Check for outer buns. Subtract 1 for each missing outer bun. Specific bun type does not matter, though.
	if len(actualIngredients) <= 0 || actualIngredients[0] != IngBunBottom && actualIngredients[0] != IngBunTop {
		revenue -= 1
	}
	if len(actualIngredients) <= 1 || actualIngredients[len(actualIngredients)-1] != IngBunBottom && actualIngredients[len(actualIngredients)-1] != IngBunTop {
		revenue -= 1
	}

	// Ensure min revenue is capped at zero.
	return max(revenue, 0)
}
