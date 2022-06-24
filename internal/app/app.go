package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"math/rand"
	"realty_bot/internal/models"
	"time"
)

type UserServiceInterface interface {
	List() []models.User
	Subscribed(int64) bool
	Subscribe(models.User) error
	UnSubscribe(int64) error
}

type AdServiceInterface interface {
	NewAds() ([]models.Ad, error)
}

type App struct {
	config      *Config
	bot         *tgbotapi.BotAPI
	logger      *zerolog.Logger
	userService UserServiceInterface
	adService   AdServiceInterface
}

// NewApp returns app instance
func NewApp(logger *zerolog.Logger, config *Config,
	userService UserServiceInterface, adService AdServiceInterface) *App {
	// bot init
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		logger.Panic().Err(err)
	}
	app := &App{
		config:      config,
		bot:         bot,
		userService: userService,
		adService:   adService,
		logger:      logger,
	}
	return app
}

func (a *App) CommandHandler() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := a.bot.GetUpdatesChan(u)

	for update := range updates {
		// ignore any non-Message updates and non-command Messages
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "start":
			a.commandHelp(update.Message.Chat.ID)
		case "help":
			a.commandHelp(update.Message.Chat.ID)
		case "query":
			a.commandQuery(update.Message.Chat.ID)
		case "subscribers":
			a.commandSubscribers(update.Message.Chat.ID)
		case "subscribe":
			a.commandSubscribe(update.Message.Chat.ID, update.Message.Chat.UserName, update.Message.Chat.FirstName, update.Message.Chat.LastName)
		case "unsubscribe":
			a.commandUnsubscribe(update.Message.Chat.ID)
		}
	}
}

func (a *App) commandHelp(chatID int64) {
	message := "<strong>Available commands:</strong>\n" +
		"/query - current query to parse\n" +
		"/subscribe - subscribe\n" +
		"/unsubscribe - unsubscribe\n" +
		"/subscribers - list of subscribers\n" +
		"\n/help - this help information\n"
	if err := a.sendMessageToChat(chatID, message, false); err != nil {
		a.logger.Error().Err(err)
	}
}

func (a *App) commandQuery(chatID int64) {
	message := fmt.Sprintf("<strong>Current search query is:</strong>\n%s", a.config.Query)
	if err := a.sendMessageToChat(chatID, message, false); err != nil {
		a.logger.Error().Err(err)
	}
}

func (a *App) commandSubscribers(chatID int64) {
	list := a.userService.List()
	message := "<strong>No subscribers found</strong>\nUse /subscribe command to add yourself =)"
	if len(list) > 0 {
		message = "<strong>Current subscribers:</strong>\n"
		for _, user := range list {
			line := ""
			if user.UserName != "" {
				line += fmt.Sprintf("@%s ", user.UserName)
			}
			if user.Name != "" {
				line += user.Name
			}
			message += line + "\n"
		}
	}
	err := a.sendMessageToChat(chatID, message, false)
	if err != nil {
		a.logger.Error().Err(err)
	}
}

func (a *App) commandSubscribe(chatID int64, username string, firstName string, lastName string) {
	if a.userService.Subscribed(chatID) {
		err := a.sendMessageToChat(chatID,
			"<strong>You already subscribed! =)</strong>\nWait for notifications",
			false)
		if err != nil {
			a.logger.Error().Msg(err.Error())
		}
		return
	}
	user := models.User{
		ChatID:   chatID,
		UserName: username,
	}

	user.Name = fmt.Sprintf("%s %s",
		lastName, firstName)

	if err := a.userService.Subscribe(user); err != nil {
		err = a.sendMessageToChat(chatID,
			"<strong>Error! =)</strong>\nSomething went wrong: "+err.Error(),
			false)
		if err != nil {
			a.logger.Error().Err(err)
		}
		a.logger.Error().Err(err)
		return
	}
	err := a.sendMessageToChat(chatID,
		"<strong>You successfully subscribed! =)</strong>\nWait for new notifications",
		false)
	if err != nil {
		a.logger.Error().Err(err)
	}
}

func (a *App) commandUnsubscribe(chatID int64) {
	if !a.userService.Subscribed(chatID) {
		err := a.sendMessageToChat(chatID,
			"<strong>You are not subscribed! =)</strong>\nUse /subscribe command to subscribe notifications",
			false)
		if err != nil {
			a.logger.Error().Err(err)
		}
		return
	}
	if err := a.userService.UnSubscribe(chatID); err != nil {
		err = a.sendMessageToChat(chatID,
			"<strong>Error! =)</strong>\nSomething went wrong: "+err.Error(),
			false)
		if err != nil {
			a.logger.Error().Msg(err.Error())
		}
		a.logger.Error().Msg(err.Error())
		return
	}
	err := a.sendMessageToChat(chatID,
		"<strong>You successfully unsubscribed! =(</strong>",
		false)
	if err != nil {
		a.logger.Error().Msg(err.Error())
	}
}

func (a *App) sendMessageToChat(chatID int64, message string, showPreview bool) error {
	msgConf := tgbotapi.NewMessage(chatID, message)
	msgConf.ParseMode = "HTML"
	if !showPreview {
		msgConf.DisableWebPagePreview = true
	}
	_, err := a.bot.Send(msgConf)
	return err
}

func (a *App) Run() {
	a.logger.Info().Msgf("Bot started...")
	a.logger.Info().Msgf("Config is %#v", a.config)
	for {
		if len(a.userService.List()) == 0 {
			a.logger.Warn().Msgf("No subscribers found=(")
			time.Sleep(time.Second * 5)
			continue
		}
		time.Sleep(time.Second*10 + time.Duration(rand.Intn(1000))*time.Millisecond)
		a.logger.Info().Msgf("Parsing new ads...")
		newAds, err := a.adService.NewAds()
		if err != nil {
			a.logger.Error().Msg(err.Error())
			continue
		}
		if len(newAds) == 0 {
			a.logger.Info().Msgf("No new ads found =(")
			continue
		}
		a.logger.Info().Msgf("New ads=%d found!", len(newAds))

		for _, user := range a.userService.List() {
			for _, ad := range newAds {
				err = a.sendMessageToChat(user.ChatID,
					fmt.Sprintf("<strong>%s</strong>\n<code>€%d</code>\n<i>%s</i>\n\n%s", ad.Title, ad.Price, ad.Location, ad.Link),
					true)
				if err != nil {
					a.logger.Err(err)
					continue
				}
				time.Sleep(time.Second)
			}
		}
	}
}
