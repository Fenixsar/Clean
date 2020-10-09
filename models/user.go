package models

import (
	"github.com/gorilla/websocket"
	"gitlab.q123123.net/ligmar/boot"
)

type User struct {
	Channel    boot.Channel    `json:"-" bson:"-"`
	WS         *websocket.Conn `json:"-" bson:"-"`
	FirstName  string          `json:"firstName" bson:"firstName"`
	LastName   string          `json:"lastName" bson:"lastName"`
	Username   string          `json:"username" bson:"username"`
	TelegramID int64           `json:"-" bson:"key"`
	Role       string          `json:"role" bson:"role"`
}

type SocketReceiveChannel struct {
	Event string
	Data  string
	Cb    string
}

func (u *User) Init() {
	u.Channel = make(boot.Channel, 0)
}
