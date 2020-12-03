package mongo

import (
	"testing"
	"time"

	"github.com/theconcierge/recipes/core"
)

func TestStore(t *testing.T) {
	repo, err := NewMongoRespoitory("mongodb://localhost:27017", "recipes", 100)
	if err != nil {
		t.Errorf("got err %s", err.Error())
	}

	recipe := &core.Recipe{
		Name:        "Chicken Parmesan",
		CreatedAt:   time.Now().UTC().Unix(),
		Picture:     "https://imagesvc.meredithcorp.io/v3/mm/image?url=https%3A%2F%2Fimages.media-allrecipes.com%2Fuserphotos%2F3371942.jpg&w=958&h=958&c=sc&poi=face&q=85",
		Description: "Chicken parm is quite possibly the best meal ever invented.",
		Blurb:       "Prepare for chicken goodness",
		Parts:       []core.RecipePart{},
	}
	recipe.Parts = append(recipe.Parts, core.RecipePart{
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

	recipe.Parts = append(recipe.Parts, core.RecipePart{
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

	err = repo.Store(recipe)
	if err != nil {
		t.Errorf("got err %s", err.Error())
	}
}

func TestFind(t *testing.T) {
	repo, err := NewMongoRespoitory("mongodb://localhost:27017", "recipes", 100)
	if err != nil {
		t.Errorf("got err %s", err.Error())
	}

	recipe, err := repo.Find("chicken-parmesan-2020-11-29")
	if err != nil {
		t.Errorf("got err %s", err.Error())
		// if we get an error retrieving, no point testing anything else
		// until it is fixed
		return
	}

	expectedName := "Chicken Parmesan"
	actualName := recipe.Name
	if actualName != expectedName {
		t.Errorf("Expected %s and got %s", expectedName, actualName)
	}
}
