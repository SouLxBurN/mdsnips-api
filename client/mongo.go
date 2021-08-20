package client

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoClient Returns a reference to a mongo.Client
// host - required
// port - required
// user - optional pass is required, if provided
// pass - optional user is required, if provided
func InitMongoClient(host string, port string, user string, pass string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := ""
	if user != "" && pass != "" {
		auth = fmt.Sprintf("%s:%s@", user, pass)
	}
	connString := fmt.Sprintf("mongodb://%s%s:%s", auth, host, port)
	clientOptions := options.Client().ApplyURI(connString)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
