package core

import "bytes"

// RecipeService inteface for our core domain
type RecipeService interface {
	SingleRecipePage(id string) (*bytes.Buffer, error)
	IndexPage() (*bytes.Buffer, error)
}
