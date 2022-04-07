package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var authOpt = options.Credential{
	AuthMechanism: os.Getenv("AUTH"),
	AuthSource:    os.Getenv("DB_NAME"),
	Username:      os.Getenv("DB_USER"),
	Password:      os.Getenv("DB_PWD"),
}

func ConnectDB(host string, port string) (*mongo.Client, error) {

	uri := fmt.Sprintf("mongodb://%s:%s", host, port)
	clientopt := options.Client().SetAuth(authOpt).ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientopt)

	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
