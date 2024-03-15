package usecase

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/akshay0074700747/project-and-company-management-chat-service/entities"
	"github.com/akshay0074700747/project-and-company-management-chat-service/helpers"
)

func (chatt *ChatUsecase) InsertIntoDB() chan<- entities.InsertIntoRoomMessage {

	insertChan := make(chan entities.InsertIntoRoomMessage, 100)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	var run = true

	go func() {
		defer func() {
			for v := range insertChan {
				if err := chatt.Repo.InsertMessage(v); err != nil {
					helpers.PrintErr(err,"error happened at InsertMessage adapter")
				}
			}
			close(insertChan)
			close(sigchan)
			helpers.PrintMsg("the insert into mongo routine is shutting down...")
		}()

		for run {
			select {
			case <-sigchan:
				run = false
			case msg := <-insertChan:
				if err := chatt.Repo.InsertMessage(msg); err != nil {
					helpers.PrintErr(err,"error happened at InsertMessage adapter")
				}
			}
		}
	}()

	return insertChan
}
