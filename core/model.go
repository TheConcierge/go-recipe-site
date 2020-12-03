package core

// RecipePart data model for a subsection of the recipe
type RecipePart struct {
	PartName     string   `json:"part_name" bson:"part_name" msgpack:"part_name"`
	Instructions []string `json:"instructions" bson:"instructions" msgpack:"instructions"`
	Ingredients  []string `json:"ingredients" bson:"ingredients" msgpack:"ingredients"`
}

// Recipe data model for a single recipe
type Recipe struct {
	Name        string       `json:"recipe_name" bson:"recipe_name" msgpack:"recipe_name"`
	UniqueID    string       `json:"unique_id" bson:"unique_id" msgpack:"unique_id"`
	Picture     string       `json:"picture_url" bson:"picture_url" msgpack:"picture_url" validate:"format=url"`
	Description string       `json:"description" bson:"description" msgpack:"description"`
	Parts       []RecipePart `json:"parts" bson:"parts" msgpack:"parts"`
	CreatedAt   int64        `json:"created_at" bson:"created_at" msgpack:"created_at"`
	Blurb       string       `json:"blurb" bson:"blurb" msgpack:"blurb"`
}
