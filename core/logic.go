package core

import (
	"bytes"
	"errors"
    "time"
    "strings"

	"html/template"
)

var (
	ErrTemplateLoadFail = errors.New("Could not load template file")
	ErrRecipeLoadFail   = errors.New("Could not load recipe")
    ErrRecipeSearchFail = errors.New("Could not search for recipe")
    ErrRecipeStoreFail = errors.New("Could not store new recipe")
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
	recipes, err := r.recipeRepo.MostRecent(12)
	if err != nil {
        ep := buildErrorPage(ErrRecipeLoadFail)
		return ep, ErrRecipeLoadFail
	}
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
        ep := buildErrorPage(ErrTemplateLoadFail)
		return ep, ErrTemplateLoadFail
	}

	page := new(bytes.Buffer)

	tmpl.Execute(page, recipes)

	return page, nil
}

func (r *recipeService) DefaultRecipeSearch() (*bytes.Buffer, error) {
    tmpl, err := template.ParseFiles("/.templates/search.html")
    if err != nil {
        ep := buildErrorPage(ErrRecipeSearchFail)
        return ep, nil
    }

    page := new(bytes.Buffer)

    tmpl.Execute(page, nil)

    return page, nil
}

func (r *recipeService) SearchRecipes(name string) (*bytes.Buffer, error) {
    recipes, err := r.recipeRepo.Search(name)
    if err != nil {
        ep := buildErrorPage(ErrRecipeSearchFail)
        return ep, nil
    }

    tmpl, err := template.ParseFiles("./templates/search.html")
    if err != nil {
        ep := buildErrorPage(ErrTemplateLoadFail)
        return ep, nil
    }

    page := new(bytes.Buffer)

    tmpl.Execute(page, recipes)

    return page, nil
}

func (r *recipeService) Inject() (*bytes.Buffer, error) {
    ut := time.Now().UTC().Unix()
    uid := getUniqueID("Chicken Parm", ut)
    recipe := &Recipe{
		Name:        "Chicken Parm",
        UniqueID:    uid,
		CreatedAt:   ut,
		Picture:     "https://imagesvc.meredithcorp.io/v3/mm/image?url=https%3A%2F%2Fimages.media-allrecipes.com%2Fuserphotos%2F3371942.jpg&w=958&h=958&c=sc&poi=face&q=85",
		Description: "Chicken parm is quite possibly the best meal ever invented.",
		Blurb:       "Prepare for chicken goodness",
		Parts:       []RecipePart{},
    }
	recipe.Parts = append(recipe.Parts, RecipePart{
		PartName: "Chicken",
		Instructions: []string{
			"cook chicken",
			"eat chicken",
		},
		Ingredients: []string{
			"chicken",
			"pepper",
		},
	})

	recipe.Parts = append(recipe.Parts, RecipePart{
		PartName: "Tomato Sauce",
		Instructions: []string{
			"cook sauce",
			"put sauce on chicken",
		},
		Ingredients: []string{
			"can of tomatos",
			"basil",
		},
	})

    err := r.recipeRepo.Store(recipe)
    if err != nil {
        ep := buildErrorPage(ErrRecipeStoreFail)
        return ep, nil
    }

    tmpl, err := template.ParseFiles("./templates/newrecipe.html")
    if err != nil {
        ep := buildErrorPage(ErrTemplateLoadFail)
        return ep, nil
    }

    page := new(bytes.Buffer)

    type ts struct {
        Msg string
    }
    tmpl.Execute(page, &ts{Msg: "injected stock recipe"})

    return page, nil
}
// helper functions for exposed service functions

func buildErrorPage(e error) *bytes.Buffer {
    tmpl, err := template.ParseFiles("./templates/error.html")
    if err != nil {
        superError := "Error page generation broke, literally everything is broken"
        ep := new(bytes.Buffer)
        ep.WriteString(superError)
        return ep
    }

    page := new(bytes.Buffer)
    type es struct {
        Error string
    }
    tmpl.Execute(page, &es{Error: e.Error()})

    return page
}

func getUniqueID(name string, timestamp int64) string {
	nameLower := strings.ToLower(name)
	nameDash := strings.Replace(nameLower, " ", "-", -1)

	dateDash := time.Unix(timestamp, 0).Format("2006-01-02")

	return nameDash + "-" + dateDash
}
