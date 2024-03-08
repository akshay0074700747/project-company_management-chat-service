package repository

import "go.mongodb.org/mongo-driver/mongo"

type ChatRepo struct {
	DB *mongo.Database
}

func NewChatRepo(db *mongo.Database) *ChatRepo {
	return &ChatRepo{
		DB: db,
	}
}

