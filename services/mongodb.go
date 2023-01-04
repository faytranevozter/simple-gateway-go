package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func ConnectMongoDB(timeout time.Duration, URI string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client()
	clientOptions.ApplyURI(URI)
	connection, _ := connstring.Parse(URI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return client.Database(connection.Database)
}
