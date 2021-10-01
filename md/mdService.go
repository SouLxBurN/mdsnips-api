package md

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	database   = "mdsnips"
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
		Title:      mdSnip.Title,
		UpdateKey:  createUpdateKey(mdSnip.Body),
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
	opts := options.FindOne().SetProjection(bson.M{"updateKey": 0})
	if err := mdCollection.FindOne(ctx, filter, opts).Decode(snippet); err != nil {
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
func (m *MDService) GetAllMarkdownSnippets() ([]MDListItem, error) {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	snippets := make([]MDListItem, 0)
	filter := bson.D{}
	opts := options.Find().SetProjection(bson.M{"id": 1, "title": 1, "createDate": 1})
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

	return snippets, nil
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
	snippet.Title = patch.Title
	snippet.Body = patch.Body

	filter := bson.D{{Key: "id", Value: snippet.ID}}
	update := bson.M{"$set": snippet}
	_, err = mdCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return snippet, nil
}

// ValidateIdAndKey
// Fetch by snippet by Id and validate against updateKey
func (m *MDService) ValidateIdAndKey(uuid string, updateKey string) bool {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	snippet := new(MarkdownSnippet)
	filter := bson.D{{Key: "id", Value: uuid}}
	opts := options.FindOne().SetProjection(bson.M{"updateKey": 1})
	if err := mdCollection.FindOne(ctx, filter, opts).Decode(snippet); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false
		}
		return false
	}
	return updateKey == snippet.UpdateKey
}

// DeleteMarkdownSnippet
// Errors are returned to the caller
func (m *MDService) DeleteMarkdownSnippet(uuid string) {

}

// getMarkdownCollection
// Accepts a mongo.Client reference
// Returns a reference to the `mdsnips.markdown` collection.
func getMarkdownCollection(mClient *mongo.Client) *mongo.Collection {
	return mClient.Database("mdsnips").Collection("markdown")
}

// createUpdateKey
// Generates an update key based on the markdown content
func createUpdateKey(content string) string {
	content += strconv.Itoa(time.Now().Nanosecond())
	algo := crc32.NewIEEE()
	algo.Write([]byte(content))
	return fmt.Sprintf("%x", algo.Sum32())
}
