package testutils

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// MongoTestContainer
// Container - Reference to GenericContainer Object
// Host - IP of the running container
// Port - Port of the running container
// NOTE: Terminate the Container reference when done.
type MongoTestContainer struct {
	Container testcontainers.Container
	Host      string
	Port      string
}

// setupMongoTestContainer Creates a mongo test container
// for testing integrations with mongo.
func SetupMongoTestContainer() (*MongoTestContainer, error) {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println("filename: " + filename)
	// The ".." may change depending on you folder structure
	dir := path.Join(path.Dir(filename), "..")
	fmt.Println("directory: " + dir)
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017"},
		WaitingFor:   wait.ForListeningPort("27017"),
	}

	mCont, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := mCont.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := mCont.MappedPort(ctx, "27017")
	if err != nil {
		return nil, err
	}

	return &MongoTestContainer{Container: mCont, Host: ip, Port: port.Port()}, nil
}
