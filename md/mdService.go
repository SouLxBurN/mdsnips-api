package md

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	database   = "soulxsnips"
	collection = "markdown"
)

type MDService struct {
	client *mongo.Client
}

// InitMDService Creates an instance of a MDService
// Requires a reference to a mongo.Client instance
func InitMDService(mClient *mongo.Client) *MDService {
	return &MDService{client: mClient}
}

// CreateMarkdownSnippet
// Errors are returned to the caller
func (m *MDService) CreateMarkdownSnippet(mdSnip *CreateMDReq) (*MarkdownSnippet, error) {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	newSnip := &MarkdownSnippet{
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
func (m *MDService) GetMarkdownSnippet(uuid string) (*MarkdownSnippet, error) {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	snippet := new(MarkdownSnippet)
	filter := bson.D{{Key: "id", Value: uuid}}
	if err := mdCollection.FindOne(ctx, filter).Decode(snippet); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return snippet, nil
}

// GetAllMarkdownSnippets
// Gets all Markdown Snippets without body
// Errors are returned to the caller
func (m *MDService) GetAllMarkdownSnippets() (*[]MDListItem, error) {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	snippets := make([]MDListItem, 0)
	filter := bson.D{}
	opts := options.Find().SetProjection(bson.M{"id": 1, "createDate": 1})
	cursor, err := mdCollection.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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
func (m *MDService) UpdateMarkdownSnippet(patch *UpdateMDReq) (*MarkdownSnippet, error) {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Get Existing Snippet
	snippet, err := m.GetMarkdownSnippet(patch.ID)
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
func (m *MDService) DeleteMarkdownSnippet(uuid string) {

}

// getMarkdownCollection
// Accepts a mongo.Client reference
// Returns a reference to the `soulxsnips.markdown` collection.
func getMarkdownCollection(mClient *mongo.Client) *mongo.Collection {
	return mClient.Database("soulxsnips").Collection("markdown")
}
