package handler

import (
	"net/http"
	"regexp"
	"strings"

	"gitlab.q123123.net/ligmar/console-back/pkg/global"

	"gitlab.q123123.net/ligmar/console-back/pkg/repository"

	"gitlab.q123123.net/ligmar/console-back/pkg/service"

	"gitlab.q123123.net/ligmar/console-back/models"

	"gitlab.q123123.net/ligmar/boot"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) ws(c *gin.Context) {
	var user models.User
	var err error

	if boot.Prod {
		user, err = h.userIdentity(c)
		if err != nil {
			boot.Log.Warn(err)
			return
		}
	} else {
		user, err = h.services.Authentication.Check(33835172)
		if err != nil {
			boot.Log.Warn(err)
			return
		}
		wsUpgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	user.Init()

	socket, err := wsHandler(c.Writer, c.Request, &user)
	if err != nil {
		return
	}

	service.NewService(user, socket, h.telegram, repository.NewRepository())
}

func wsHandler(w http.ResponseWriter, r *http.Request, user *models.User) (ws *websocket.Conn, err error) {

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		boot.Log.Warnf("Failed to set websocket upgrade: %+v", err)
		return nil, err
	}

	go func(с boot.Channel) {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}

			m := string(msg)
			if our, _ := regexp.MatchString(`^\d*\[\"([А-яёЁ\w]+)\"(\]|\,.+\])$`, m); !our {
				boot.Log.Error("Somebody try send wrong socket events!")
				_ = conn.Close()
				return
			}

			go boot.ChannelEmitter(с, parseReceiveMassage(m))
		}

		go boot.ChannelEmitter(с, boot.ChannelForm{
			Name: global.SocketReceiver,
			Data: models.SocketReceiveChannel{
				Event: global.SocketDisconnect,
			},
		})
	}(user.Channel)

	return conn, nil
}

func parseReceiveMassage(m string) (massage boot.ChannelForm) {
	var event, data, cb string
	if strings.Index(m, ",") >= 0 {
		event = m[strings.Index(m, "[")+2 : strings.Index(m, ",")-1]
		data = m[strings.Index(m, ",")+1 : len(m)-1]
	} else {
		event = m[strings.Index(m, "[")+2 : len(m)-2]
	}
	if strings.Index(m, "[") >= 0 {
		cb = m[0:strings.Index(m, "[")]
	}

	massage = boot.ChannelForm{
		Name: global.SocketReceiver,
		Data: models.SocketReceiveChannel{
			Data:  data,
			Event: event,
			Cb:    cb,
		},
	}

	return
}
