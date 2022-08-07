package main

import (
	"github.com/bopoh24/realty-bot/internal/app"
	"github.com/bopoh24/realty-bot/internal/service"
	"github.com/bopoh24/realty-bot/internal/store/filestore"
	"github.com/bopoh24/realty-bot/pkg/log"
)

func main() {

	logger := log.NewLogger()

	// user service init
	config, err := app.NewConfig()
	if err != nil {
		logger.Fatal().Msgf("Unable to load config: %s", err)
	}

	// user service init
	userService, err := service.NewChatService(filestore.NewChatStore(config.FileUsers))
	if err != nil {
		logger.Fatal().Msgf("Unable to init user service:", err.Error())
	}

	// ad service init
	adParseService, err := service.NewAdParseService(config.Query, logger, filestore.NewAdStore(config.FileAds))
	if err != nil {
		logger.Fatal().Msgf("Unable to init ad service:", err.Error())
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
