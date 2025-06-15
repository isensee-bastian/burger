package burger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRevenue(t *testing.T) {
	order := []IngredientType{IngBunBottom, IngTomatoes, IngPattyBeef, IngKetchup, IngCheese, IngSalad, IngBunTop}

	assert.Equal(t, 0, revenue(order, []IngredientType{}), "Empty burger should yield 0")
	assert.Equal(t, 0, revenue(order, []IngredientType{IngBunBottom}), "Just a single bun should yield 0")
	assert.Equal(t, 2, revenue(order, []IngredientType{IngBunBottom, IngBunTop}), "Just the outer buns should yield 2")
	assert.Equal(t, 2, revenue(order, []IngredientType{IngBunBottom, IngTomatoes, IngPattyBeef}), "Missing top bun should subtract 1")
	assert.Equal(t, 3, revenue(order, []IngredientType{IngTomatoes, IngPattyBeef, IngKetchup, IngBunTop}), "Missing bottom bun should subtract 1")
	assert.Equal(t, 3, revenue(order, []IngredientType{IngTomatoes, IngPattyBeef, IngKetchup, IngCheese, IngSalad}), "Missing outer buns should subtract 2")
	assert.Equal(t, 6, revenue(order, []IngredientType{IngBunBottom, IngTomatoes, IngPattyBeef, IngCheese, IngSalad, IngBunTop}), "Missing ingredient should subtract 1")

	assert.Equal(t, 7, revenue(order, []IngredientType{IngBunBottom, IngTomatoes, IngPattyBeef, IngKetchup, IngCheese, IngSalad, IngBunTop}), "All ingredients and outer buns present should yield full score")

	assert.Equal(t, 6, revenue(order, []IngredientType{IngBunBottom, IngTomatoes, IngPattyBeef, IngKetchup, IngMayo, IngCheese, IngSalad, IngBunTop}), "An unexpected ingredient shall subtract 1")
	assert.Equal(t, 5, revenue(order, []IngredientType{IngBunBottom, IngPattyBeef, IngKetchup, IngMayo, IngCheese, IngSalad, IngBunTop}), "A missing and an unexpected ingredient shall subtract 2")
	assert.Equal(t, 7, revenue(order, []IngredientType{IngBunBottom, IngPattyBeef, IngTomatoes, IngSalad, IngCheese, IngKetchup, IngBunTop}), "Ingredient order shall not matter for non-buns")
	assert.Equal(t, 5, revenue(order, []IngredientType{IngTomatoes, IngBunBottom, IngPattyBeef, IngKetchup, IngCheese, IngBunTop, IngSalad}), "Ingredient order shall matter for buns")
	assert.Equal(t, 7, revenue(order, []IngredientType{IngBunTop, IngPattyBeef, IngTomatoes, IngSalad, IngCheese, IngKetchup, IngBunBottom}), "Swapping outer bottom and top buns shall not matter")

	assert.Equal(t, 7, revenue(order, []IngredientType{IngBunBottom, IngTomatoes, IngPattyBeef, IngKetchup, IngCheese, IngSalad, IngCheese, IngBunTop}), "Duplicate expected ingredients shall not matter")
	assert.Equal(t, 6, revenue(order, []IngredientType{IngBunBottom, IngTomatoes, IngPattyBeef, IngMayo, IngKetchup, IngMayo, IngCheese, IngSalad, IngBunTop}), "Duplicate unexpected ingredients shall only subtract 1")
}
