package core

import (
	"bytes"
	"errors"
	"html/template"
)

var (
	ErrTemplateLoadFail = errors.New("Could not load template file")
	ErrRecipeLoadFail   = errors.New("Could not load recipe")
)

type recipeService struct {
	recipeRepo RecipeRepository
}

// NewRecipeService creates new service for core domain
func NewRecipeService(recipeRepo RecipeRepository) RecipeService {
	return &recipeService{
		recipeRepo,
	}
}

// SingleRecipePage builds html page for a recipe given its id
func (r *recipeService) SingleRecipePage(id string) (*bytes.Buffer, error) {
	recipe, err := r.recipeRepo.Find(id)
	if err != nil {
		return nil, ErrRecipeLoadFail
	}

	tmpl, err := template.ParseFiles("./templates/recipe.html")
	if err != nil {
		return nil, ErrTemplateLoadFail
	}

	page := new(bytes.Buffer)

	tmpl.Execute(page, recipe)

	return page, nil
}

func (r *recipeService) IndexPage() (*bytes.Buffer, error) {
	// index highlights 12 most recent recipes
	recipes, err := r.recipeRepo.MostRecent(9)
	if err != nil {
		return nil, ErrRecipeLoadFail
	}

	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		return nil, ErrTemplateLoadFail
	}

	page := new(bytes.Buffer)

	tmpl.Execute(page, recipes)

	return page, nil
}
