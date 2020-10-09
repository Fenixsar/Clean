package service

import (
	"encoding/json"
	"errors"

	"gitlab.q123123.net/ligmar/console-back/pkg/global"

	"gitlab.q123123.net/ligmar/boot"
	"gitlab.q123123.net/ligmar/console-back/models"

	"github.com/gorilla/websocket"
)

const (
	// User
	getAccount = "getAccount"

	// Support
	getChatList   = "getChatList"
	getMessages   = "getMessages"
	sendMessage   = "sendMessage"
	deleteMessage = "deleteMessage"
	markAsRead    = "markAsRead"

	newMessage = "newMessage"
)

func (s *Service) SocketHandler(data models.SocketReceiveChannel) {
	//User
	ok := s.UserRouter(data)
	if !ok {
		return
	}

	// Support
	ok = s.SupportRouter(data)
	if !ok {
		return
	}

	boot.Log.Errorf("SocketHandler() - %s", errors.New("socket handler not found"))
}

func (s *Service) UserRouter(data models.SocketReceiveChannel) (ok bool) {
	switch data.Event {
	case getAccount:
		s.Emit(data.Cb, s.User)
		return
	}

	return true
}
func (s *Service) SupportRouter(data models.SocketReceiveChannel) (ok bool) {
	switch data.Event {
	case getChatList:
		answer, err := s.SupportService.GetChatList(0)
		if err != nil {
			boot.Log.Errorf("SupportRouter() - %s - %s", getChatList, err)
			s.Emit(data.Cb, false)
			return
		}

		s.Emit(data.Cb, answer)
		return
	case getMessages:
		var param struct {
			ChatID int `json:"chatID"`
			Skip   int `json:"skip"`
		}

		if err := json.Unmarshal([]byte(data.Data), &param); err != nil {
			boot.Log.Errorf("SupportRouter() - %s - %s", getMessages, err)
			s.Emit(data.Cb, false)
			return
		}

		answer, err := s.SupportService.GetMessages(param.ChatID, param.Skip)
		if err != nil {
			boot.Log.Errorf("SupportRouter() - %s - %s", getMessages, err)
			s.Emit(data.Cb, false)
			return
		}
		s.Emit(data.Cb, answer)
		return
	case sendMessage:
		var param struct {
			ChatID int64  `json:"chatID"`
			Text   string `json:"text"`
		}

		if err := json.Unmarshal([]byte(data.Data), &param); err != nil {
			boot.Log.Errorf("SupportRouter() - %s - %s", sendMessage, err)
			s.Emit(data.Cb, false)
			return
		}

		answer, err := s.telegram.SendMessage(param.ChatID, param.Text)
		if err != nil {
			boot.Log.Errorf("SupportRouter() - %s - %s", sendMessage, err)
			s.Emit(data.Cb, false)
			return
		}

		message := models.Message{
			MessageID: answer.MessageID,
			Text:      param.Text,
			Date:      boot.GetTimeStamp(),
			ChatID:    param.ChatID,
			Sender:    s.TelegramID,
			Unread:    true,
		}

		err = s.SupportService.WriteMessage(message)
		if err != nil {
			boot.Log.Errorf("SupportRouter() - %s - %s", sendMessage, err)
			s.Emit(data.Cb, false)
			return
		}

		s.Emit(data.Cb, true)
		global.BroadcastToUsers(boot.ChannelForm{
			Name: global.BroadcastMessage,
			Data: models.Chat{
				ChatID:  param.ChatID,
				Message: message,
			},
		})
		return

	}

	return true
}

func (s *Service) Emit(event string, data interface{}) {
	var temp string
	msg, err := json.Marshal(data)
	if err != nil {
		boot.Log.Warnf("SocketEmit: can't marshal:  Data. \n%v", data)
		return
	}
	temp += "[\"" + event + "\"," + string(msg) + "]"
	w, err := s.ws.NextWriter(websocket.TextMessage)
	if err == nil {
		_, _ = w.Write([]byte(temp))
		_ = w.Close()
		return
	}

	boot.Log.Warn(err)
}
