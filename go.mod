module gitlab.q123123.net/ligmar/console-back

go 1.15

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/gorilla/websocket v1.4.2
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	gitlab.q123123.net/ligmar/boot v1.0.23
	go.mongodb.org/mongo-driver v1.4.1
	golang.org/x/net v0.0.0-20200904194848-62affa334b73
	gopkg.in/telegram-bot-api.v4 v4.6.4
)

replace gitlab.q123123.net/ligmar/boot v1.0.23 => ../../boot
