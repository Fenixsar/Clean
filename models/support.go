package models

type Chat struct {
	Key       string  `json:"key" bson:"key"`
	ChatID    int64   `json:"chatID" bson:"chatID"`
	FirstName string  `json:"firstName" bson:"firstName"`
	LastName  string  `json:"lastName" bson:"lastName"`
	Message   Message `json:"message" bson:"message"`
}

type ChatList []Chat
type MessageList []Message

type Message struct {
	MessageID int    `json:"messageID" bson:"messageID"`
	Text      string `json:"text" bson:"text"`
	File      *File  `json:"file,omitempty" bson:"file,omitempty"`
	Date      int    `json:"date" bson:"date"`
	ChatID    int64  `json:"-" bson:"chatID"`
	Sender    int64  `json:"sender" bson:"sender"`
	Unread    bool   `json:"unread" bson:"unread"`
}

type File struct {
	ID       string `json:"ID" bson:"ID"`
	Name     string `json:"name" bson:"name"`
	Caption  string `json:"caption" bson:"caption"`
	MimeType string `json:"mimeType" bson:"mimeType"`
}
