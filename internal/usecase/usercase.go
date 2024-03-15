package usecase

import (
	"strings"

	"github.com/akshay0074700747/project-and-company-management-chat-service/entities"
	"github.com/akshay0074700747/project-and-company-management-chat-service/helpers"
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

func (chatt *ChatUsecase) SpinupPoolifnotalreadyExists(poolid string, insertChan chan<- entities.InsertIntoRoomMessage, isPrivate bool) (*chat.Pool, []entities.Message) {

	if isPrivate {

		res, err := chatt.Repo.LoadMessagesofPrivate(poolid)
		if err != nil {
			helpers.PrintErr(err, "error happened at LoadMesseges")
		}

		if chatt.Chat.Pool[poolid] == nil {
			ids := strings.Split(poolid, " ")
			if chatt.Chat.Pool[ids[1]+" "+ids[0]] == nil {
				pool := chat.NewPool(ids[1] + " " + ids[0])
				go pool.Serve(insertChan)
				chatt.Chat.Pool[ids[1]+" "+ids[0]] = pool
				return pool, res
			}
			return chatt.Chat.Pool[ids[1]+" "+ids[0]], res
		}
		return chatt.Chat.Pool[poolid], res
	}

	res, err := chatt.Repo.LoadMesseges(poolid)
	if err != nil {
		helpers.PrintErr(err, "error happened at LoadMesseges")
	}

	if chatt.Chat.Pool[poolid] == nil {
		pool := chat.NewPool(poolid)
		go pool.Serve(insertChan)
		chatt.Chat.Pool[poolid] = pool
		return pool, res
	}

	return chatt.Chat.Pool[poolid], res
}
