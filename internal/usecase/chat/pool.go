package chat

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akshay0074700747/project-and-company-management-chat-service/entities"
	"github.com/akshay0074700747/project-and-company-management-chat-service/helpers"
)

type Pool struct {
	ID           string
	RegisterChan chan *Client
	UnRegister   chan *Client
	Broadcast    chan entities.Message
	Clients      map[string]*Client
}

func NewPool(id string) *Pool {
	return &Pool{
		ID:           id,
		RegisterChan: make(chan *Client),
		UnRegister:   make(chan *Client),
		Broadcast:    make(chan entities.Message),
		Clients:      make(map[string]*Client),
	}
}

type RegisterNUnregister struct {
	Message string    `json:"Message"`
	Time    time.Time `json:"Time"`
}

func (pool *Pool) Serve(insertChan chan<- entities.InsertIntoRoomMessage) {
	defer func() {
		close(pool.RegisterChan)
		close(pool.UnRegister)
		close(pool.Broadcast)
	}()
	for {
		select {
		case client := <-pool.RegisterChan:
			for _, v := range pool.Clients {
				reg := RegisterNUnregister{
					Time:    time.Now(),
					Message: fmt.Sprintf("%s is Online", client.Name),
				}
				if err := v.Conn.WriteJSON(reg); err != nil {
					helpers.PrintErr(err, "error happened at sending")
					continue
				}
			}
			pool.Clients[client.ClientID] = client
		case client := <-pool.UnRegister:
			for _, v := range pool.Clients {
				unReg := RegisterNUnregister{
					Time:    time.Now(),
					Message: fmt.Sprintf("%s is Offline", client.Name),
				}
				if err := v.Conn.WriteJSON(unReg); err != nil {
					helpers.PrintErr(err, "error happened at sending")
					continue
				}
			}
			delete(pool.Clients, client.ClientID)
		case message := <-pool.Broadcast:
			for _, v := range pool.Clients {
				jsonDta, err := json.Marshal(message)
				if err != nil {
					helpers.PrintErr(err, "error happened at sending")
					continue
				}
				if err := v.Conn.WriteMessage(message.Type, jsonDta); err != nil {
					helpers.PrintErr(err, "error happened at sending")
					continue
				}
			}
			masg := entities.InsertIntoRoomMessage{
				RoomID:   pool.ID,
				Messages: message,
			}
			insertChan <- masg
		}
	}
}
