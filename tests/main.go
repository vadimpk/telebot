package main

import (
	"gopkg.in/telebot.v3"
)

func main() {
	webhook := &telebot.Webhook{
		Listen: ":8080",
		AllowedUpdates: []string{
			"message",
			"callback_query",
			"chat_member",
			"photo",
			"video",
			"animation",
		},
		IP:                "",
		DropUpdates:       false,
		SecretToken:       "",
		HasCustomCert:     false,
		PendingUpdates:    0,
		ErrorUnixtime:     0,
		ErrorMessage:      "",
		SyncErrorUnixtime: 0,
		TLS:               nil,
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: "https://9499-176-36-201-123.ngrok-free.app",
		},
	}

	handler := telebot.NewHandler(telebot.HandlerSettings{})

	handler.Handle("/start", func(c telebot.Context) error {
		_, err := c.Bot().Send(c.Sender(), "Hello, "+c.Sender().FirstName+"!")
		return err
	})

	updates := make(chan telebot.Update, 100)
	go webhook.Start(updates)

	settings := telebot.Settings{
		Token:   "7013258382:AAEjpmkP2YpGEoeI-EfjebTBl5xkzAVLVPg",
		Handler: handler,
	}

	b, err := telebot.NewBot(settings)
	if err != nil {
		panic(err)
	}

	err = b.SetWebhook(webhook, map[string]string{
		"shopID": "1",
	})
	if err != nil {
		panic(err)
	}

	b.Updates = updates

	settings.Token = "7050792454:AAGcxoWUQkRphw9SQ-ixgPG-onsLGiVzgn4"
	b2, err := telebot.NewBot(settings)
	if err != nil {
		panic(err)
	}

	err = b2.SetWebhook(webhook, map[string]string{
		"shopID": "2",
	})
	if err != nil {
		panic(err)
	}

	b2.Updates = updates

	for update := range b.Updates {
		if update.Args["shopID"] == "1" {
			b.ProcessUpdate(update)
		} else {
			b2.ProcessUpdate(update)
		}
	}
}
