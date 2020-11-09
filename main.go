package main

import (
	"log"
	"net/http"
	"os"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
	c "github.com/lithiferous/go-fundmgr/coms"
)

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm fundmgr bot!"))
}

func main() {
	bot, err := t.NewBotAPI(os.Getenv("TG"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	http.HandleFunc("/", MainHandler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	//init core
	const sep = " "
	const fp = "./data"
	l, s := c.InState(fp)

	ups := fetchUpdates(bot)

	for up := range ups {
		r := ""
		if up.Message == nil {
			continue
		}
		switch up.Message.Command() {
		case "sup":
			r = c.Status(*s)
		case "upd":
			e, delta := c.Delta(c.DropCmd(up.Message.Text, up.Message.Command()), sep, *l)
			switch e {
			case "":
				r = c.Eval(delta, l)
				c.OutState(fp, l)
			default:
				r = e
			}
		case "add":
			e, payer, person := c.Person(c.DropCmd(up.Message.Text, up.Message.Command()), sep, *l)
			switch e {
			case "":
				r = c.Payer(l, s, person, payer)
				c.OutState(fp, l)
			default:
				r = e
			}
		case "pay":
			r = c.Pay(c.DropCmd(up.Message.Text, up.Message.Command()), sep, &s)
			log.Printf("resulted %d\n", c.OutState(fp, l))
		case "help":
			r = "Список комманд:\n" +
				"/sup - покажет текущий галактический счёт магов\n" +
				"/pay `cумма` - вносит общую покупку для всех магов\n" +
				"/add `плательщик` `имя мага` - добавит плательщика\n" +
				"/upd `плательщик` `сумма` - частный перевод\n" +
				"/help - отобразить это сообщение\n" +
				"(p.s. писать без кавычек)\n"
		}
		if r != "" {
			msg := t.NewMessage(up.Message.Chat.ID, r)
			bot.Send(msg)
		}
	}
}

func fetchUpdates(bot *t.BotAPI) t.UpdatesChannel {
	_, err := bot.SetWebhook(t.NewWebhook("https://go-fundmgr.herokuapp.com/" + bot.Token))
	if err != nil {
		log.Printf("webhook err - %w ", err)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	return updates
}
