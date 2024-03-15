package chat

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akshay0074700747/project-and-company-management-chat-service/entities"
	"github.com/akshay0074700747/project-and-company-management-chat-service/helpers"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	ClientID string
	Name     string
	Pool     *Pool
}

func NewClient(conn *websocket.Conn, clientID, name string, pool *Pool) *Client {
	return  &Client{
		Conn:     conn,
		ClientID: clientID,
		Name:     name,
		Pool:     pool,
	}
}

func (client *Client) Serve(msgs []entities.Message) {

	client.Pool.RegisterChan <- client

	defer func() {
		client.Pool.UnRegister <- client
		client.Conn.Close()
	}()

	for _, v := range msgs {
		jsonDta, err := json.Marshal(v)
		if err != nil {
			helpers.PrintErr(err, "error happened at sending")
			continue
		}
		if err := client.Conn.WriteMessage(v.Type, jsonDta); err != nil {
			helpers.PrintErr(err, "error happened at sending")
			continue
		}
	}

	for {
		msgtype, p, err := client.Conn.ReadMessage()
		if err != nil {
			helpers.PrintErr(err, "error happened , closing connection")
			break
		}

		message := entities.Message{Type: msgtype,Message: string(p),Time: time.Now(),Name: client.Name}
		client.Pool.Broadcast <- message
		fmt.Printf("Message Received from %s", client.ClientID)
	}
}
