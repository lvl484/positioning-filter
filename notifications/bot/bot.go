package bot

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//const TgbotapiKey = "1192652390:AAEmF-btBEDG0uCkdOAEDwEXdTJM4mIqTWA"

//var chatID int64

type NotificationBot struct {
	bot *tgbotapi.BotAPI
}

func NewNotificationBot(config *Config) (*NotificationBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.TgbotapiKey)
	if err != nil {
		return nil, err
	}

	return &NotificationBot{
			bot: bot,
		},
		nil
}

// SendNotification sends chattable item to appropriate chat with information in message
func (n *NotificationBot) SendNotification(message string) error {
	_, err := n.bot.Send(tgbotapi.NewMessage(410287439, message))
	if err != nil {
		return err
	}
	return nil
}

func (n *NotificationBot) Bot() error {

	log.Printf("Authorized on account %s", n.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := n.bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}

	for update := range updates {
		if err != nil {
			log.Println(err)
		}
		reply := "I do not know how to answer"
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %v %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

		switch update.Message.Command() {
		case "start":
			reply = "Your chatID: " + strconv.FormatInt(update.Message.Chat.ID, 10) + ". Please point it in the app and I will send you notifications."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := n.bot.Send(msg)
		if err != nil {
			return err
		}

	}
	return nil
}

/*func main() {
	b, err := NewNotificationBot()
	if err != nil {
		log.Println(err)
		return
	}
	go b.Bot()
	for {
		b.SendNotification("yep")
		time.Sleep(time.Second)
	}
}
*/
