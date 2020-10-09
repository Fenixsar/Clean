package repository

import (
	"context"

	"gitlab.q123123.net/ligmar/console-back/models"

	"gitlab.q123123.net/ligmar/boot"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Support struct{}

func NewSupportRepository() *Support {
	return &Support{}
}

func (s Support) GetChatList(skip int) (list models.ChatList, err error) {

	//db.getCollection('support').aggregate([
	//{$sort: {"date":-1}},
	//{$group: {_id:"$chatID", "message":{$first:"$$ROOT"}}},
	//{$sort: {"message.date":-1}},
	//{$skip: 0},
	//{$limit: 30}
	//])

	pipeline := []bson.M{
		boot.Sort(bson.M{"date": -1}),
		boot.Group(bson.M{"_id": "$chatID", "message": bson.M{"$first": "$$ROOT"}}),
		boot.Sort(bson.M{"message.date": -1}),
		boot.Skip(int64(skip)),
		boot.Limit(100),
	}

	cursor, err := boot.SupportMessageConsoleCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var temp struct {
			User    models.User    `json:"user" bson:"user"`
			Message models.Message `json:"message" bson:"message"`
		}
		err := cursor.Decode(&temp)
		if err != nil {
			continue
		}
		list = append(list, models.Chat{
			ChatID:  temp.Message.ChatID,
			Message: temp.Message,
		})
	}

	return list, nil
}
func (s Support) GetMessages(chatID, skip int) (list models.MessageList, err error) {
	option := options.Find().SetSort(bson.D{{"date", -1}}).SetSkip(int64(skip)).SetLimit(30)
	cursor, err := boot.SupportMessageConsoleCollection.Find(context.TODO(), bson.M{"chatID": chatID}, option)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var message models.Message
		err := cursor.Decode(&message)
		if err != nil {
			continue
		}

		list = append(list, message)
	}

	return list, nil
}
func (s Support) DeleteMessage(messageID int) (err error) {
	err = boot.DeleteOneByMatch(boot.SupportMessageConsoleCollection, bson.M{"messageID": messageID})
	return
}
func (s Support) StoreMessage(message models.Message) (err error) {
	err = boot.InsertOne(boot.SupportMessageConsoleCollection, message)
	return
}
func (s Support) MarkAsRead(chatID int) (err error) {
	_, err = boot.UpdateFields(boot.SupportMessageConsoleCollection, bson.M{"chatID": chatID}, bson.M{"$set": bson.M{"unread": false}})
	return
}
