package repository

import "github.com/akshay0074700747/project-and-company-management-chat-service/entities"

type ChatRepoInterface interface {
	InsertMessage(entities.InsertIntoRoomMessage) error
	LoadMesseges(string) ([]entities.Message, error)
	LoadMessagesofPrivate(string)([]entities.Message, error)
}
