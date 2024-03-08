package usecase

import "github.com/akshay0074700747/project-and-company-management-chat-service/internal/usecase/chat"

type ChatUseaseInterface interface {
	SpinupPoolifnotalreadyExists(string) *chat.Pool
}
