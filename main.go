package main

import (
	"gitlab.q123123.net/ligmar/boot"
	"gitlab.q123123.net/ligmar/console-back/models"
	"gitlab.q123123.net/ligmar/console-back/pkg/auth"
	"gitlab.q123123.net/ligmar/console-back/pkg/handler"
	"gitlab.q123123.net/ligmar/console-back/pkg/repository"
	"gitlab.q123123.net/ligmar/console-back/pkg/telegram"
)

func main() {
	boot.LoadBoot()
	boot.MongoDB()

	repos := repository.NewRepository()

	bot := telegram.NewTelegram(repos)

	services := auth.NewAuth(repos)
	handlers := handler.NewHandler(services, bot)

	srv := new(models.Server)

	if err := srv.Run(boot.ConsoleServerPort, handlers.InitRoutes()); err != nil {
		boot.Log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
