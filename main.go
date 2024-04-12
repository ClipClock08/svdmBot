package main

import (
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strconv"
)

func main() {

	token := os.Getenv("BOT_TOKEN")
	chatIdent := os.Getenv("CHAT_ID")
	forwardChat, _ := strconv.Atoi(chatIdent)

	botToken := flag.String("bot-token", token, "Telegram Bot API token")
	chatID := flag.Int64("chat-id", int64(forwardChat), "Telegram chat ID to forward messages to")
	flag.Parse()

	if *botToken == "" {
		log.Fatal("missing bot token argument")
	}

	if *chatID == 0 {
		log.Fatal("missing chat ID argument")
	}

	bot, err := tgbotapi.NewBotAPI(*botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		forwardMessage(bot, update.Message, *chatID)
	}
}

func forwardMessage(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	chatID int64,
) {
	msg := tgbotapi.NewForward(chatID, message.Chat.ID, message.MessageID)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error forwarding message: %v", err)
	}
}
