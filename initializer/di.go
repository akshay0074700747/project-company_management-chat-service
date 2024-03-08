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
	return handlers.NewChatHandlers(usecase, ":50003", ":50002",":50001")
}
