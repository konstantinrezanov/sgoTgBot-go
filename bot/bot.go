package bot

import (
	"fmt"
	"log"
	"sgoTgBot-go/conversion"
	"strings"
	"sync"
	"io/ioutil"
	"os/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StartBot(schedule conversion.Schedule, m *sync.Mutex) {
	bot, err := tgbotapi.NewBotAPI(getToken())
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message.IsCommand() { // ignore any non-command Messages
			switch update.Message.Command() {
			case "start":
				msg.Text = "Запрос в форме: класс(ex.10Б) дата(ex.15.02.2021)\nКомманда /get показывает доступные дни"
			case "help":
				msg.Text = "Запрос в форме: класс(ex.10Б) дата(ex.15.02.2021)\nКомманда /get показывает доступные дни"
			case "get":
				m.Lock()
				dates := schedule.AvailableDates()
				m.Unlock()
				ret := "Доступные даты:\n"
				for _, day := range dates {
					ret += fmt.Sprintf("%v\n", day)
				}
				msg.Text = ret
			default:
				msg.Text = "I don't know that command"
			}
		} else {
			text := strings.Split(update.Message.Text, " ")
			if len(text) < 2 {
				msg.Text = "Not correct"
			} else {
				m.Lock()
				msg.Text = schedule.PrettyPrint(strings.ToUpper(text[0]), text[1])
				m.Unlock()
			}
		}

		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func getToken() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := usr.HomeDir + "/.config/sgobot/cred.config"
	content,err:=ioutil.ReadFile(path)
	if err!=nil {
		log.Fatal(err)
	}

	return strings.Split(string(content[:len(content)-1]), " ")[2]
}