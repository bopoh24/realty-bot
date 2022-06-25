package main

import (
	"github.com/bopoh24/realty-bot/internal/app"
	"github.com/bopoh24/realty-bot/internal/service"
	"github.com/bopoh24/realty-bot/internal/store/filestore"
	"github.com/bopoh24/realty-bot/pkg/log"
)

// TODO: coverage with gomock

func main() {

	logger := log.NewLogger()
	config := app.MustConfig()

	// user service init
	userService, err := service.NewUserService(filestore.NewUserStore(config.FileUsers))
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	// ad service init
	adParseService, err := service.NewAdParseService(config.Query, logger, filestore.NewAdStore(config.FileAds))
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	a := app.NewApp(
		logger,
		config,
		userService,
		adParseService)

	// start command handler
	go a.CommandHandler()
	// run app
	a.Run()
}
