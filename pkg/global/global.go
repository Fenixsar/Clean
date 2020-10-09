package global

import (
	"sync"

	"gitlab.q123123.net/ligmar/console-back/models"

	"gitlab.q123123.net/ligmar/boot"
)

const (
	SocketReceiver   = "socketReceiver"
	SocketDisconnect = "socketDisconnect"
	SocketNewConnect = "socketNewConnect"

	BroadcastMessage = "broadcastMessage"
)

type userList struct {
	sync.RWMutex
	List map[int64]boot.Channel
}

var users = userList{
	List: make(map[int64]boot.Channel, 0),
}

func AddToUsers(id int64, channel boot.Channel) {
	users.Lock()
	if c, ok := users.List[id]; ok {
		go boot.ChannelEmitter(c, boot.ChannelForm{
			Name: SocketReceiver,
			Data: models.SocketReceiveChannel{
				Event: SocketNewConnect,
			},
		})
	}
	users.List[id] = channel
	users.Unlock()
}

func RemoveFromUsers(id int64, channel boot.Channel) {
	users.Lock()
	if c, ok := users.List[id]; ok && c == channel {
		delete(users.List, id)
	}

	users.Unlock()
}

func BroadcastToUsers(form boot.ChannelForm) {
	users.RLock()
	for _, channel := range users.List {
		go boot.ChannelEmitter(channel, form)
	}
	users.RUnlock()
}
