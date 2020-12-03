package core

// RecipeRepository interface for recipe storage
type RecipeRepository interface {
	Find(id string) (*Recipe, error)
	Store(recipe *Recipe) error
	MostRecent(numResults int) ([]*Recipe, error)
}
