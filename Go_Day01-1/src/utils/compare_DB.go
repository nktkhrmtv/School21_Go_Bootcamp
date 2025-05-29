package utils

import (
	"fmt"
)

func CompareDatabases(oldDB, newDB *Recipes) {
	oldCakes := make(map[string]Cake)
	for _, cake := range oldDB.Cakes {
		oldCakes[cake.Name] = cake
	}
	newCakes := make(map[string]Cake)
	for _, cake := range newDB.Cakes {
		newCakes[cake.Name] = cake
	}

	for name := range newCakes {
		oldCake, exists := oldCakes[name]; 
		if !exists {
			fmt.Printf("ADDED cake \"%s\"\n", name)
		}
		if exists{
			compareCakes(oldCake, newCakes[name], name)
		}
	}
	for name := range oldCakes {
		if _, exists := newCakes[name]; !exists {
			fmt.Printf("REMOVED cake \"%s\"\n", name)
		}
	}
}

func compareCakes(oldCake, newCake Cake, name string) {
	if newCake.StoveTime != oldCake.StoveTime {
		fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", name, newCake.StoveTime, oldCake.StoveTime)
	}

	oldIngredients := make(map[string]Ingredient)
	for _, ingredient := range oldCake.Ingredients {
		oldIngredients[ingredient.Name] = ingredient
	}
	newIngredients := make(map[string]Ingredient)
	for _, ingredient := range newCake.Ingredients {
		newIngredients[ingredient.Name] = ingredient
	}

	for ingredientName := range newIngredients {
		if _, exists := oldIngredients[ingredientName]; !exists {
			fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", ingredientName, name)
		}
	}
	for ingredientName := range oldIngredients {
		if _, exists := newIngredients[ingredientName]; !exists {
			fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", ingredientName, name)
		}
	}

	for ingredientName, newIngredient := range newIngredients {
		if oldIngredient, exists := oldIngredients[ingredientName]; exists {
			if newIngredient.Count != oldIngredient.Count && newIngredient.Unit == oldIngredient.Unit {
				fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", ingredientName, name, newIngredient.Count, oldIngredient.Count)
			}
			if newIngredient.Unit != oldIngredient.Unit {
				if oldIngredient.Unit == "" {
					fmt.Printf("ADDED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", newIngredient.Unit, ingredientName, name)
				} else if newIngredient.Unit == "" {
					fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", oldIngredient.Unit, ingredientName, name)
				} else {
					fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", ingredientName, name, newIngredient.Unit, oldIngredient.Unit)
				}
			}
		}
	}
}
