package md

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// ConfigureIndexes
// Creates/Updates the markdown collection indexes.
func ConfigureIndexes(mClient *mongo.Client) {
	collection := getMarkdownCollection(mClient)

	index := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "title", Value: bsonx.String("text")},
				{Key: "body", Value: bsonx.String("text")},
			},
		},
		{
			Keys: bsonx.Doc{{Key: "createDate", Value: bsonx.Int32(1)}},
		},
	}
	name, err := collection.Indexes().CreateMany(context.TODO(), index)
	if err != nil {
		fmt.Printf("Error Creating Text Index: %s", err)
		return
	}
	fmt.Printf("Index Created: %s\n", name)
}
