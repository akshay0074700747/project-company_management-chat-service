package usecase

import (
	"github.com/akshay0074700747/project-and-company-management-chat-service/internal/repository"
	"github.com/akshay0074700747/project-and-company-management-chat-service/internal/usecase/chat"
)

type ChatUsecase struct {
	Repo repository.ChatRepoInterface
	Chat *chat.ChatPool
}

func NewChatUsecase(repo repository.ChatRepoInterface) *ChatUsecase {
	return &ChatUsecase{
		Repo: repo,
		Chat: chat.NewChatPool(),
	}
}

func (chatt *ChatUsecase) SpinupPoolifnotalreadyExists(poolid string) *chat.Pool {

	if chatt.Chat.Pool[poolid] == nil {
		pool := chat.NewPool(poolid)
		go pool.Serve()
		chatt.Chat.Pool[poolid] = pool
		return pool
	}

	return chatt.Chat.Pool[poolid]
}