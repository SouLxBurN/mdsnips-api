package md

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/soulxburn/soulxsnips/client"
	"github.com/soulxburn/soulxsnips/testutils"

	"github.com/stretchr/testify/assert"
)

var mdService *MDService

func TestMain(m *testing.M) {
	// Initialize
	mCont, err := testutils.SetupMongoTestContainer()
	if err != nil {
		log.Fatal("Failed to initialize mongo container")
	}

	mClient, err := client.InitMongoClient(mCont.Host, mCont.Port, "", "")
	if err != nil {
		log.Fatal("Failed to connection to mongo container")
	}
	mdService = InitMDService(mClient)

	// Execute all tests in this file.
	retCode := m.Run()

	// Cleanup
	mCont.Container.Terminate(context.Background())
	os.Exit(retCode)
}

func TestCreateMarkdownSnippet(t *testing.T) {
	req := &CreateMDReq{Body: "# Title\n##Subhead\nDetails...Details...Details..."}
	snippet, err := mdService.CreateMarkdownSnippet(req)
	assert.Nil(t, err)
	assert.NotNil(t, snippet)
}
