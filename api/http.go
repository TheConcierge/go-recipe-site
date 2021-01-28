package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/theconcierge/recipes/core"
)

type RecipeHandler interface {
	SingleRecipe(http.ResponseWriter, *http.Request)
	Index(http.ResponseWriter, *http.Request)
    Search(http.ResponseWriter, *http.Request)
    Inject(http.ResponseWriter, *http.Request)
}

type handler struct {
	recipeService core.RecipeService
}

func NewHandler(recipeService core.RecipeService) RecipeHandler {
	return &handler{recipeService: recipeService}
}

func (h *handler) SingleRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["recipe_id"]

	page, err := h.recipeService.SingleRecipePage(id)

	if err != nil {
		fmt.Fprintf(w, "Could not load page\n%s", err.Error())
	}
	//fmt.Println(page.Bytes())
	w.Write(page.Bytes())
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) {
	page, err := h.recipeService.IndexPage()

	if err != nil {
		fmt.Fprintf(w, "Could not load page\n%s", err.Error())
	}
	//fmt.Println(page.Bytes())
	w.Write(page.Bytes())
}

func (h *handler) Inject(w http.ResponseWriter, r *http.Request) {
    page, err := h.recipeService.Inject()

    if err != nil {
        fmt.Fprintf(w, "Could not load page\n$s", err.Error())
    }

    w.Write(page.Bytes())
}

func (h *handler) Search(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("search_name")
    fmt.Sprintf("form value %s", name)
    page, err := h.recipeService.SearchRecipes(name)

    if err != nil {
        fmt.Fprintf(w, "Could not load page\n%s", err.Error())
    }

    w.Write(page.Bytes())
}
