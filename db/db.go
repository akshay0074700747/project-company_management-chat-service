package db

import (
	"context"
	"log"

	"github.com/akshay0074700747/project-and-company-management-chat-service/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(cfg *config.Config) *mongo.Database {

	clientOptions := options.Client().ApplyURI(cfg.MongoUrl)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database("chatServiceDB")

}
