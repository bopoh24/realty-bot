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

	config, err := app.NewConfig()
	if err != nil {
		logger.Panic().Msgf("Unable to load app config: %s", err)
	}

	// user service init
	userService, err := service.NewUserService(filestore.NewUserStore(config.FileUsers))
	if err != nil {
		logger.Panic().Msg(err.Error())
	}

	// ad service init
	adParseService := service.NewAdParseService(config.Query, logger, filestore.NewAdStore(config.FileAds))

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
