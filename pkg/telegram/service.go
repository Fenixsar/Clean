package telegram

import (
	"errors"
	"fmt"

	"gitlab.q123123.net/ligmar/console-back/pkg/global"
	"gitlab.q123123.net/ligmar/console-back/pkg/repository"

	bootModels "gitlab.q123123.net/ligmar/boot/models"
	"gitlab.q123123.net/ligmar/console-back/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.q123123.net/ligmar/boot"
)

type Telegram struct {
	Bot        *tgbotapi.BotAPI
	Repository *repository.Repository
}

func NewTelegram(repos *repository.Repository) *Telegram {
	bot, err := tgbotapi.NewBotAPI(boot.TelegramSupportBotAPI)
	if err != nil {
		boot.Log.Error(err)
	}

	bot.Debug = false

	boot.Log.Printf("Authorized on account %s", bot.Self.UserName)
	msg := tgbotapi.NewMessage(boot.TelegramNotificationID, "<b>❗️ Рестарт Telegram сервера!!!</b>")
	msg.ParseMode = "HTML"
	bot.Send(msg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	t := &Telegram{
		Bot:        bot,
		Repository: repos,
	}
	go t.checkUpdates(updates)

	return t
}

func (t Telegram) checkUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "/start":
			text := fmt.Sprintf("<b>Приветствуем!✋</b>\nЭто служба поддержки MMORPG игры LIGMAR @LigmarBot\n\n<b>Файловые ограничения:</b>\nРазмер: <i>до 20мб</i>\nИзображения: <i>jpeg,png,svg,webp,...</i>\nВидео: <i>mp4,mpeg,quicktime,avi,flv,wmv,...</i>\n\nОставьте сообщение и мы ответим в ближайшее время ⬇️")
			_, _ = t.SendMessage(update.Message.Chat.ID, text)
		default:
			user, err := t.checkUser(update.Message)
			if err != nil {
				text := fmt.Sprintf("К сожалению, вы не можете отправлять сообщения, пока не активировали бота @Ligmar")
				_, _ = t.SendMessage(update.Message.Chat.ID, text)
				continue
			}

			err = t.Repository.Telegram.Store(user)
			if err != nil {
				boot.Log.Warnf("checkUpdates() - not save user: %s", err)
			}

			message, err := t.processingMessage(update.Message)
			if err != nil {
				boot.Log.Infof("checkUpdates() - not processing message: %s", err)
				continue
			}

			if message.File != nil {

			}

			err = t.Repository.Support.StoreMessage(message)
			if err != nil {
				boot.Log.Warnf("checkUpdates() - not save message: %s", err)
			}

			go global.BroadcastToUsers(boot.ChannelForm{
				Name: global.BroadcastMessage,
				Data: models.Chat{
					Key:       user.Key,
					ChatID:    user.TelegramID,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Message:   message,
				},
			})

		}
	}
}
func (t Telegram) checkUser(message *tgbotapi.Message) (user bootModels.UserTelegram, err error) {
	user, err = t.Repository.Telegram.GetByTelegramID(message.Chat.ID)
	if err != nil {
		return
	}

	user.FirstName = message.Chat.FirstName
	user.LastName = message.Chat.LastName
	user.Username = message.Chat.UserName
	user.Language = message.From.LanguageCode

	return
}
func (t Telegram) processingMessage(message *tgbotapi.Message) (newMessage models.Message, err error) {

	var fileID, mimeType, emoji string
	var fileSize int

	if message.Photo != nil { // фото
		for key, value := range *message.Photo {
			// сохранение большого разрешения
			if key == len(*message.Photo)-1 {
				fileID = value.FileID
				fileSize = value.FileSize
				mimeType = "image/jpeg"
			}
		}
	} else if message.Document != nil {
		fileID = message.Document.FileID
		mimeType = message.Document.MimeType
		fileSize = message.Document.FileSize

	} else if message.Video != nil {
		fileID = message.Video.FileID
		mimeType = message.Video.MimeType
		fileSize = message.Video.FileSize

	} else if message.Sticker != nil {
		emoji = message.Sticker.Emoji
	}

	newMessage.ChatID = message.Chat.ID
	newMessage.MessageID = message.MessageID
	newMessage.Date = boot.GetTimeStamp()
	newMessage.Unread = true

	if fileID != "" {
		if mimeTypeCheck, ok := MimeTypeCheck[mimeType]; !ok || mimeTypeCheck == "" {
			err = errors.New("mimeType not found")
			return
		}

		// файл должен быть меньше 20 мб
		if fileSize > 1024*1024*20 {
			err = errors.New("file size so big")
			return
		}
		newMessage.File = &models.File{
			ID:       fileID,
			Caption:  message.Caption,
			MimeType: mimeType,
		}
	} else {
		if emoji != "" {
			newMessage.Text = emoji
		} else {
			newMessage.Text = message.Text
		}
	}

	return
}

func (t Telegram) SendMessage(chatID int64, message string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true // отключить показ веб ссылок

	answer, err := t.Bot.Send(msg)
	return answer, err
}

// MimeTypeCheck - допустимые форматы, а так же идентификация
var MimeTypeCheck = map[string]string{
	"text/html":                     ".html",
	"text/css":                      ".css",
	"text/xml":                      ".xml",
	"image/gif":                     ".gif",
	"image/jpeg":                    ".jpg",
	"application/x-javascript":      ".js",
	"application/atom+xml":          ".atom",
	"application/rss+xml":           ".rss",
	"text/plain":                    ".txt",
	"image/png":                     ".png",
	"image/tiff":                    ".tiff",
	"image/vnd.wap.wbmp":            ".wbmp",
	"image/x-icon":                  ".ico",
	"image/x-ms-bmp":                ".bmp",
	"image/svg+xml":                 ".svg",
	"image/webp":                    ".webp",
	"application/msword":            ".doc",
	"application/pdf":               ".pdf",
	"application/rtf":               ".rtf",
	"application/vnd.ms-excel":      ".xls",
	"application/vnd.ms-powerpoint": ".ppt",
	"application/x-7z-compressed":   ".7z",
	"application/x-rar-compressed":  ".rar",
	"application/zip":               ".zip",
	"video/mp4":                     ".mp4",
	"video/mpeg":                    ".mpeg",
	"video/quicktime":               ".mov",
	"video/webm":                    ".webm",
	"video/x-flv":                   ".flv",
	"video/x-m4v":                   ".m4v",
	"video/x-mng":                   ".mng",
	"video/x-ms-asf":                ".asf",
	"video/x-ms-wmv":                ".wmv",
	"video/x-msvideo":               ".avi",
}
