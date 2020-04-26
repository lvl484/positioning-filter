// Package bot provides telegram bot and functionality to send notifications
package bot

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type NotificationBot struct {
	bot *tgbotapi.BotAPI
	log *logrus.Logger
}

func NewNotificationBot(config *Config, log *logrus.Logger) (*NotificationBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.TgbotapiKey)
	if err != nil {
		return nil, err
	}

	return &NotificationBot{
			bot: bot,
			log: log,
		},
		nil
}

// SendNotification sends chattable item to appropriate chat with filtered position in message
func (n *NotificationBot) SendNotification(chatID int64, message string) error {
	_, err := n.bot.Send(tgbotapi.NewMessage(chatID, message))
	if err != nil {
		return err
	}
	return nil
}

// Bot implements telegram bot configuration
func (n *NotificationBot) Bot() error {

	n.log.Infof("Authorized on account %s", n.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := n.bot.GetUpdatesChan(u)
	if err != nil {
		n.log.Errorf("Can't start channel for getting updates %v", err)
	}

	for update := range updates {
		if err != nil {
			n.log.Errorf("Can't handle updates %v", err)
		}
		reply := "I do not know how to answer"
		if update.Message == nil {
			continue
		}

		n.log.Infof("[%s] %v %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

		switch update.Message.Command() {
		case "start":
			reply = "Your chatID: " + strconv.FormatInt(update.Message.Chat.ID, 10) + ". Please point it in the app and I will send you notifications."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := n.bot.Send(msg)
		if err != nil {
			n.log.Errorf("Can't send message %v", err)
			return err
		}

	}
	return nil
}
