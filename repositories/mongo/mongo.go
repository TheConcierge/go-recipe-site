package mongo

import (
	"context"
    "time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/theconcierge/recipes/core"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewMongoRespoitory initializes repository
func NewMongoRespoitory(mongoURL, mongoDB string, mongoTimeout int) (core.RecipeRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}

	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, err
	}
	repo.client = client
	return repo, nil
}


func (r *mongoRepository) Store(recipe *core.Recipe) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("recipes")
	_, err := collection.InsertOne(
		ctx,
		recipe,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *mongoRepository) Find(id string) (*core.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	recipe := &core.Recipe{}
	collection := r.client.Database(r.database).Collection("recipes")
	filter := bson.M{"unique_id": id}

	err := collection.FindOne(ctx, filter).Decode(&recipe)
	if err != nil {
		return nil, err
	}
	return recipe, err
}

func (r *mongoRepository) MostRecent(numResults int) ([]*core.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	recipes := []*core.Recipe{}

	fo := options.Find()
	// most recent <numResults> recipes
    fo.SetSort(bson.M{"created_at": -1})
    fo.SetLimit(int64(numResults))

	collection := r.client.Database(r.database).Collection("recipes")

	cur, err := collection.Find(ctx, bson.D{}, fo)
	if err != nil {
		return recipes, err
	}

	for cur.Next(ctx) {
		recipe := &core.Recipe{}

		err := cur.Decode(&recipe)
		if err != nil {
			return recipes, err
		}
		recipes = append(recipes, recipe)
	}
	cur.Close(ctx)

	if len(recipes) == 0 {
		return recipes, mongo.ErrNoDocuments
	}

	return recipes, nil
}


func (r *mongoRepository) Search(name string) ([]*core.Recipe, error) {
    ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
    defer cancel()

    recipes := []*core.Recipe{}

    fo := options.Find()
    fo.SetSort(bson.M{"created_at": -1})
    // setting some responsible limit
    // eventually I will paginate responses if I end up using this site enough
    fo.SetLimit(30)

    filter := bson.M{"recipe_name": primitive.Regex{Pattern: name, Options: "i"}}

	collection := r.client.Database(r.database).Collection("recipes")

    cur, err := collection.Find(ctx, filter, fo)
    if err != nil {
        return recipes, err
    }

    for cur.Next(ctx) {
        recipe := &core.Recipe{}

        err := cur.Decode(&recipe)
        if err != nil {
            // currently not sure what errors could appear here, so don't know
            // if proper behavior would be to continue or just return what we
            // have so far
            return recipes, err
        }
        recipes = append(recipes, recipe)
    }
    cur.Close(ctx)

    if len(recipes) == 0 {
        return recipes, mongo.ErrNoDocuments
    }

    return recipes, nil

}
