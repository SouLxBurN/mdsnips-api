package md

import (
	"context"
	"log"
	"testing"

	"github.com/soulxburn/soulxsnips/client"
	"github.com/soulxburn/soulxsnips/testutils"

	"github.com/stretchr/testify/assert"
)

// SetupMDService
// Creates a Mongo Test Container, Initializes a mongo client
// returns a MDService connected to the Mongo Test Container,
// and a cleanup function for tearing down the container.
func SetupMDService(t *testing.T) (*MDService, func(t *testing.T)) {
	// Initialize
	mCont, err := testutils.SetupMongoTestContainer()
	if err != nil {
		log.Fatal("Failed to initialize mongo container")
	}

	mClient, err := client.InitMongoClient(mCont.ConnectionString)
	if err != nil {
		log.Fatal("Failed to connection to mongo container")
	}

	return InitMDService(mClient), func(t *testing.T) {
		mCont.Container.Terminate(context.Background())
	}
}

// Test_CreateMarkdownSnippet
// Happy path test for creating a MarkdownSnippet.
func Test_CreateMarkdownSnippet(t *testing.T) {
	mdService, cleanup := SetupMDService(t)
	defer cleanup(t)
	expectedBody := "# Title\n##Subhead\nDetails...Details...Details..."
	req := &CreateMDReq{Body: expectedBody}

	snippet, err := mdService.CreateMarkdownSnippet(req)
	assert.Nil(t, err)
	assert.NotNil(t, snippet)
	assert.NotEmpty(t, snippet.ID)
	assert.NotEmpty(t, snippet.CreateDate)
	assert.Equal(t, expectedBody, snippet.Body)

	persistedSnip, err := mdService.GetMarkdownSnippet(snippet.ID)
	assert.Nil(t, err)
	assert.NotNil(t, persistedSnip)
	assert.NotEmpty(t, persistedSnip.ID)
	assert.NotEmpty(t, persistedSnip.CreateDate)
	assert.Equal(t, expectedBody, snippet.Body)
}

// Test_UpdateMarkdownSnippet
//
func Test_UpdateMarkdownSnippet(t *testing.T) {
	mdService, cleanup := SetupMDService(t)
	defer cleanup(t)
	// Init
	initialBody := "# Title\n##Subhead\nDetails...Details...Details..."
	updateBody := "# Dead_again_kekw\n##Update\nTest Containers is really great."
	req := &CreateMDReq{Body: initialBody}

	// Execute
	initialSnip, err := mdService.CreateMarkdownSnippet(req)
	assert.Nil(t, err)
	assert.Equal(t, initialBody, initialSnip.Body)

	upReq := &UpdateMDReq{ID: initialSnip.ID, Body: updateBody}
	updatedSnip, err := mdService.UpdateMarkdownSnippet(upReq)
	assert.Nil(t, err)
	assert.Equal(t, initialSnip.ID, updatedSnip.ID)
	assert.Equal(t, initialSnip.CreateDate.Unix(), updatedSnip.CreateDate.Unix())
	assert.Equal(t, updateBody, updatedSnip.Body)
}

// Test_GetAllMarkdownSnippets
// Should return number of snippets created,
// after creating 5 snippets.
func Test_GetAllMarkdownSnippets(t *testing.T) {
	mdService, cleanup := SetupMDService(t)
	defer cleanup(t)
	expectedBody := "# Title\n##Subhead\nDetails...Details...Details..."
	req := &CreateMDReq{Body: expectedBody}

	for i := 0; i < 5; i++ {
		_, err := mdService.CreateMarkdownSnippet(req)
		assert.Nil(t, err)
	}

	persistedSnips, err := mdService.GetAllMarkdownSnippets()
	assert.Nil(t, err)
	assert.NotNil(t, persistedSnips)
	assert.NotEmpty(t, persistedSnips)
	assert.Len(t, persistedSnips, 5)
}
