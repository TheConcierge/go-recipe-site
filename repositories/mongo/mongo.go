package mongo

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func getUniqueID(name string, timestamp int64) string {
	nameLower := strings.ToLower(name)
	nameDash := strings.Replace(nameLower, " ", "-", -1)

	dateDash := time.Unix(timestamp, 0).Format("2006-01-02")

	return nameDash + "-" + dateDash
}

func (r *mongoRepository) Store(recipe *core.Recipe) error {
	recipe.UniqueID = getUniqueID(recipe.Name, recipe.CreatedAt)

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
	fo.SetSort(bson.M{"created_at": -1})

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
