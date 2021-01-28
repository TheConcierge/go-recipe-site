package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/theconcierge/recipes/api"
	"github.com/theconcierge/recipes/core"
	"github.com/theconcierge/recipes/repositories/mongo"
)

var db core.RecipeRepository

// getDB connects to the mongo db
func getDB() core.RecipeRepository {
	db, _ = mongo.NewMongoRespoitory("mongodb://recipe-mongodb.rs.svc.cluster.local:27017", "recipes", 100)
	//db, _ = mongo.NewMongoRespoitory("mongodb://localhost:27017", "recipes", 100)
	return db
}

func main() {
	repo := getDB()
	service := core.NewRecipeService(repo)
	handler := api.NewHandler(service)

	r := mux.NewRouter()
	r.HandleFunc("/recipes/{recipe_id}", handler.SingleRecipe).Methods("GET")
	r.HandleFunc("/", handler.Index).Methods("GET")
	r.HandleFunc("/inject", handler.Inject).Methods("GET")
    r.HandleFunc("/search", handler.Search).Methods("GET")

	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))

	staticImageDirectory := http.Dir("./images/")
	staticImageHandler := http.StripPrefix("/images/", http.FileServer(staticImageDirectory))

	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	r.PathPrefix("/images/").Handler(staticImageHandler).Methods("GET")

	http.ListenAndServe(
		":3000",
		r,
	)
}
