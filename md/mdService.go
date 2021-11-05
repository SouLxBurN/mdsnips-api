package md

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"strconv"
	"time"

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

type SortBy string

const (
	CreateDate_ASC  SortBy = "createDate_ASC"
	CreateDate_DESC        = "createDate_DESC"
)

func (s *SortBy) validate() error {
	switch *s {
	case CreateDate_ASC, CreateDate_DESC:
		return nil
	}
	return errors.New("SortBy value failed validation")
}

type MDSearchParams struct {
	Text   string
	Limit  int64
	Skip   int64
	SortBy SortBy
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
		ID:         createMDID(mdSnip.Title, mdSnip.Body),
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
func (m *MDService) GetMarkdownSnippet(mdID string) (*MarkdownSnippet, error) {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	snippet := new(MarkdownSnippet)
	filter := bson.D{{Key: "id", Value: mdID}}
	opts := options.FindOne().SetProjection(bson.M{"updateKey": 0})
	if err := mdCollection.FindOne(ctx, filter, opts).Decode(snippet); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return snippet, nil
}

// @Deprecated
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

// SearchMarkdownSnippets
// Searches through all Markdown Snippets and returns then without their body.
// Errors are returned to the caller
func (m *MDService) SearchMarkdownSnippets(searchParams MDSearchParams) ([]MDListItem, error) {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sortby := bson.D{}
	switch searchParams.SortBy {
	case CreateDate_ASC:
		sortby = bson.D{{Key: "createDate", Value: 1}}
	case CreateDate_DESC:
		sortby = bson.D{{Key: "createDate", Value: -1}}
	}

	snippets := make([]MDListItem, 0)
	filter := bson.D{}
	if searchParams.Text != "" {
		filter = bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: searchParams.Text}}}}
	}
	opts := options.Find()
	opts.SetProjection(bson.M{"id": 1, "title": 1, "createDate": 1})
	opts.SetSort(sortby)
	opts.SetSkip(searchParams.Skip)
	opts.SetLimit(searchParams.Limit)
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
func (m *MDService) UpdateMarkdownSnippet(patch *UpdateMDReq) error {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// This probably isn't the best way to do this.
	type updateSnippet struct {
		Title string `bson:"title"`
		Body  string `bson:"body"`
	}

	// Update Fields
	updates := updateSnippet{patch.Title, patch.Body}

	filter := bson.D{{Key: "id", Value: patch.ID}}
	update := bson.M{"$set": updates}
	if _, err := mdCollection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

// ValidateIdAndKey
// Fetch snippet by Id and validate against updateKey
func (m *MDService) ValidateIdAndKey(mdID string, updateKey string) bool {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	snippet := make(map[string]string)
	filter := bson.D{{Key: "id", Value: mdID}}
	opts := options.FindOne().SetProjection(bson.M{"updateKey": 1})
	if err := mdCollection.FindOne(ctx, filter, opts).Decode(snippet); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false
		}
		return false
	}
	return updateKey == snippet["updateKey"]
}

// DeleteMarkdownSnippet
// Errors are returned to the caller
func (m *MDService) DeleteMarkdownSnippet(mdID string, updateKey string) error {
	mdCollection := getMarkdownCollection(m.client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.D{{Key: "id", Value: mdID}}
	if _, err := mdCollection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

// getMarkdownCollection
// Accepts a mongo.Client reference
// Returns a reference to the `mdsnips.markdown` collection.
func getMarkdownCollection(mClient *mongo.Client) *mongo.Collection {
	return mClient.Database("mdsnips").Collection("markdown")
}

// createUpdateKey
// Generates an update key based on the markdown content.
func createUpdateKey(content string) string {
	content += strconv.Itoa(time.Now().Nanosecond())
	algo := crc32.NewIEEE()
	algo.Write([]byte(content))
	return fmt.Sprintf("%x", algo.Sum32())
}

// createMDID
// Generates a ID/URL hash for the markdown snippet.
func createMDID(title string, content string) string {
	seed := title + content + strconv.Itoa(time.Now().Nanosecond())
	algo := crc32.NewIEEE()
	algo.Write([]byte(seed))
	return fmt.Sprintf("%x", algo.Sum32())
}
