package repository

import (
	"context"
	"strings"

	"github.com/akshay0074700747/project-and-company-management-chat-service/entities"
	"github.com/akshay0074700747/project-and-company-management-chat-service/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepo struct {
	DB *mongo.Database
}

func NewChatRepo(db *mongo.Database) *ChatRepo {
	return &ChatRepo{
		DB: db,
	}
}

func (chat *ChatRepo) InsertMessage(msg entities.InsertIntoRoomMessage) error {

	coll := chat.DB.Collection("chatRoomCollection")
	filter := bson.D{{"room_id", msg.RoomID}}
	update := bson.D{{"$push", bson.D{{"messages", msg.Messages}}}}
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (chat *ChatRepo) LoadMesseges(roomID string) ([]entities.Message, error) {

	coll := chat.DB.Collection("chatRoomCollection")
	filter := bson.D{{"room_id", roomID}}
	sort := bson.D{{"messages.time", 1}}

	var result struct {
		Messages []entities.Message `bson:"messages"`
	}
	if err := coll.FindOne(context.TODO(), filter, options.FindOne().SetSort(sort)).Decode(&result); err != nil {
		if err != mongo.ErrNoDocuments {
			helpers.PrintErr(err, "error happened at adapter")
		}
		return []entities.Message{}, nil
	}

	return result.Messages, nil
}

func (chat *ChatRepo) LoadMessagesofPrivate(roomID string) ([]entities.Message, error) {
	coll := chat.DB.Collection("chatRoomCollection")
	filter := bson.D{{"room_id", roomID}}
	sort := bson.D{{"messages.time", 1}}

	var result struct {
		Messages []entities.Message `bson:"messages"`
	}
	if err := coll.FindOne(context.TODO(), filter, options.FindOne().SetSort(sort)).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			rooms := strings.Split(roomID, " ")
			if err = coll.FindOne(context.TODO(), bson.D{{"room_id", rooms[1] + " " + rooms[0]}}, options.FindOne().SetSort(sort)).Decode(&result); err != nil {
				return []entities.Message{}, nil
			}
			return result.Messages, nil
		}
		helpers.PrintErr(err, "error happened at adapter")
		return []entities.Message{}, nil
	}

	return result.Messages, nil
}
