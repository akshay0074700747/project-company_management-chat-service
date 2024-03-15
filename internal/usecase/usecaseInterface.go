package usecase

import (
	"github.com/akshay0074700747/project-and-company-management-chat-service/entities"
	"github.com/akshay0074700747/project-and-company-management-chat-service/internal/usecase/chat"
)

type ChatUseaseInterface interface {
	SpinupPoolifnotalreadyExists(string, chan<- entities.InsertIntoRoomMessage,bool) (*chat.Pool, []entities.Message)
}
