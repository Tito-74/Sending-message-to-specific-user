package models

import (
	"time"

	"github.com/gofiber/websocket/v2"
)


type Client struct {
	Id   string
	Conn *websocket.Conn
}

type Message struct{
	ID string 
	CreatedAt time.Time
	Message []byte
	RoomId string
	To string
	From string
}