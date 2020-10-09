package service

import (
	"gitlab.q123123.net/ligmar/console-back/pkg/global"
	"gitlab.q123123.net/ligmar/console-back/pkg/telegram"

	"github.com/gorilla/websocket"
	"gitlab.q123123.net/ligmar/console-back/models"
	"gitlab.q123123.net/ligmar/console-back/pkg/repository"
)

type SupportService interface {
	GetChatList(skip int) (models.ChatList, error)
	GetMessages(chatID, skip int) (models.MessageList, error)
	WriteMessage(message models.Message) (err error)
	//DeleteMessage()
	//MarkAsReadChat(int) error

	//GetFile()
	//SendFile()
}
type UserService interface {
}

type Service struct {
	models.User
	ws       *websocket.Conn
	telegram telegram.Service
	UserService
	SupportService
}

func NewService(u models.User, ws *websocket.Conn, telegram telegram.Service, repos *repository.Repository) {
	s := Service{
		u,
		ws,
		telegram,
		NewUser(repos.User),
		NewSupport(repos),
	}

	s.App(u, ws)
}

func (s *Service) App(user models.User, ws *websocket.Conn) {
	s.User = user
	s.WS = ws

	go global.AddToUsers(user.TelegramID, user.Channel)

	for {
		data := <-s.Channel
		switch data.Name {
		case global.SocketReceiver:
			temp := data.Data.(models.SocketReceiveChannel)
			if temp.Event == global.SocketDisconnect {
				go global.RemoveFromUsers(user.TelegramID, s.Channel)
				return
			}

			if temp.Event == global.SocketNewConnect {
				s.WS.Close()
				break
			}

			s.SocketHandler(temp)
			break
		case global.BroadcastMessage:
			chat := data.Data.(models.Chat)
			s.Emit(newMessage, chat)
			break
		}
	}
}
