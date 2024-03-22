package initializer

import (
	"github.com/akshay0074700747/project-and-company-management-chat-service/config"
	"github.com/akshay0074700747/project-and-company-management-chat-service/db"
	"github.com/akshay0074700747/project-and-company-management-chat-service/internal/handlers"
	"github.com/akshay0074700747/project-and-company-management-chat-service/internal/repository"
	"github.com/akshay0074700747/project-and-company-management-chat-service/internal/usecase"
)

func Inject(cfg *config.Config) *handlers.ChatHandlers {
	db := db.ConnectMongo(cfg)
	repo := repository.NewChatRepo(db)
	usecase := usecase.NewChatUsecase(repo)
	insertRoom := usecase.InsertIntoDB()
	return handlers.NewChatHandlers(usecase, "company-service:50003", "project-service:50002","user-service:50001",insertRoom)
}
