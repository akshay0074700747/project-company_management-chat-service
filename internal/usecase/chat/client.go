package chat

import (
	"fmt"
	"time"

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

func (client *Client) Serve() {

	client.Pool.RegisterChan <- client

	defer func() {
		client.Pool.UnRegister <- client
		client.Conn.Close()
	}()
	for {
		msgtype, p, err := client.Conn.ReadMessage()
		if err != nil {
			helpers.PrintErr(err, "error happened , closing connection")
			break
		}

		message := Message{Type: msgtype,Message: string(p),Time: time.Now(),Name: client.Name}
		client.Pool.Broadcast <- message
		fmt.Printf("Message Received from %s", client.ClientID)
	}
}
