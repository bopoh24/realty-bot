package app

import (
	"fmt"
	"github.com/bopoh24/realty-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"math/rand"
	"time"
)

type UserServiceInterface interface {
	List() []models.User
	Exists(int64) bool
	Save(models.User) error
	Delete(int64) error
}

type AdServiceInterface interface {
	AdsToNotify() ([]models.Ad, error)
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
		logger.Fatal().Err(err)
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

func (a *App) Run() {
	a.logger.Info().Msgf("Bot started...")
	for {
		if len(a.userService.List()) == 0 {
			a.logger.Warn().Msgf("No subscribers found=(")
			time.Sleep(time.Second * 10)
			continue
		}
		a.logger.Info().Msgf("Parsing new ads...")
		newAds, err := a.adService.AdsToNotify()
		if err != nil {
			a.logger.Error().Msg(err.Error())
			a.parseTimeOut()
			continue
		}
		if len(newAds) == 0 {
			a.logger.Info().Msgf("No new ads found =(")
			a.parseTimeOut()
			continue
		}
		a.logger.Info().Msgf("New ads found: %d", len(newAds))
		for _, user := range a.userService.List() {
			for _, ad := range newAds {
				if err = a.sendAdNotification(user.ChatID, ad); err != nil {
					a.logger.Error().Msg(err.Error())
				}
			}
		}
		a.parseTimeOut()
	}
}

func (a App) parseTimeOut() {
	timeout := time.Second * 10
	if time.Now().UTC().Hour() >= 21 || time.Now().UTC().Hour() <= 4 {
		timeout = time.Minute * 30
	}
	time.Sleep(timeout + time.Duration(rand.Intn(3000))*time.Millisecond)
}

func (a App) sendAdNotification(chatID int64, ad models.Ad) error {
	return a.sendMessageToChat(chatID,
		fmt.Sprintf("<strong>%s</strong>\n"+
			"<code>€%d</code>\n"+
			"<i>%s</i>\n"+
			"%s"+
			"\n%s",
			ad.Title, ad.Price, ad.Location,
			ad.Datetime.Format("02.01.2006 15:04:05"), ad.Link), true)
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
		for i, user := range list {
			line := fmt.Sprintf("%d. ", i+1)
			if user.UserName != "" {
				line += fmt.Sprintf("@%s ", user.UserName)
			} else {
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
	if a.userService.Exists(chatID) {
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

	if err := a.userService.Save(user); err != nil {
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
	if !a.userService.Exists(chatID) {
		err := a.sendMessageToChat(chatID,
			"<strong>You are not subscribed! =)</strong>\nUse /subscribe command to subscribe notifications",
			false)
		if err != nil {
			a.logger.Error().Err(err)
		}
		return
	}
	if err := a.userService.Delete(chatID); err != nil {
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
