package repository

import (
	"context"

	"gitlab.q123123.net/ligmar/boot"
	"gitlab.q123123.net/ligmar/boot/models"
	bootModels "gitlab.q123123.net/ligmar/boot/models"
	"go.mongodb.org/mongo-driver/bson"
)

type Telegram struct {
}

func NewTelegramRepository() *Telegram {
	return &Telegram{}
}

func (t Telegram) GetByTelegramID(id int64) (userTelegram models.UserTelegram, err error) {
	err = boot.FindOneByMatch(boot.UserTelegramCollection, bson.M{"telegramID": id}, &userTelegram)
	return
}
func (t Telegram) Store(userTelegram models.UserTelegram) (err error) {
	err = boot.UpsertOne(boot.UserTelegramCollection, userTelegram.Key, userTelegram)
	return
}
func (t Telegram) GetByTelegramIDs(ids []int64) (userTelegrams []bootModels.UserTelegram, err error) {
	cursor, err := boot.UserTelegramCollection.Find(context.TODO(), bson.M{"telegramID": bson.M{"$in": ids}})
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var userTelegram models.UserTelegram
		err := cursor.Decode(&userTelegram)
		if err != nil {
			continue
		}

		userTelegrams = append(userTelegrams, userTelegram)
	}

	return
}
