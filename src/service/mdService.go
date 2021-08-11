package service

import (
	"context"
	"errors"
	"log"
	"soulxsnips/src/client"
	"soulxsnips/src/model"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMarkdownSnippet
// Errors are returned to the caller
func CreateMarkdownSnippet(mdSnip *model.CreateMarkdownSnippet) (*model.MarkdownSnippet, error) {
	mdCollection := getMarkdownCollection()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	newSnip := &model.MarkdownSnippet{
		ID:         uuid.New().String(),
		Body:       mdSnip.Body,
		CreateDate: time.Now(),
	}

	_, err := mdCollection.InsertOne(ctx, newSnip)
	if err != nil {
		return nil, err
	}

	return newSnip, nil
}

// GetMarkdownSnippet
// Errors are returned to the caller
func GetMarkdownSnippet(uuid string) (*model.MarkdownSnippet, error) {
	mdCollection := getMarkdownCollection()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	snippet := new(model.MarkdownSnippet)
	filter := bson.D{{Key: "id", Value: uuid}}
	if err := mdCollection.FindOne(ctx, filter).Decode(snippet); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return snippet, nil
}

// GetAllMarkdownSnippets
// Gets all Markdown Snippets without body
// Errors are returned to the caller
func GetAllMarkdownSnippets() (*[]model.MarkdownSnippetListItem, error) {
	mdCollection := getMarkdownCollection()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	snippets := make([]model.MarkdownSnippetListItem, 0)
	filter := bson.D{}
	opts := options.Find().SetProjection(bson.M{"id": 1, "createDate": 1})
	cursor, err := mdCollection.Find(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	if err = cursor.All(ctx, &snippets); err != nil {
		return nil, err
	}

	return &snippets, nil
}

// UpdateMarkdownSnippet
// Errors are returned to the caller
func UpdateMarkdownSnippet(patch *model.UpdateMarkdownSnippet) (*model.MarkdownSnippet, error) {
	mdCollection := getMarkdownCollection()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Get Existing Snippet
	snippet, err := GetMarkdownSnippet(patch.ID)
	if snippet == nil && err == nil {
		return nil, errors.New("Markdown Snippet does not exist")
	}

	// Update Fields
	snippet.Body = patch.Body

	filter := bson.D{{Key: "id", Value: snippet.ID}}
	update := bson.M{"$set": snippet}
	_, err = mdCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return snippet, nil
}

// DeleteMarkdownSnippet
// Errors are returned to the caller
func DeleteMarkdownSnippet(uuid string) {

}

// getMarkdownCollection
// Returns a reference to the `soulxsnips.markdown` collection.
func getMarkdownCollection() *mongo.Collection {
	mongo, err := client.GetMongoClient()
	if err != nil {
		log.Fatal("Unable to Connect to MongoDB")
	}

	return mongo.Database("soulxsnips").Collection("markdown")
}
