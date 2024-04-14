package main

import (
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	bot    *tgbotapi.BotAPI
	chatID *int64
)

func init() {
	var err error

	token := os.Getenv("BOT_TOKEN")
	chatIdent := os.Getenv("CHAT_ID")
	forwardChat, _ := strconv.Atoi(chatIdent)

	botToken := flag.String("bot-token", token, "Telegram Bot API token")
	chatID = flag.Int64("chat-id", int64(forwardChat), "Telegram chat ID to forward messages to")
	flag.Parse()

	if *botToken == "" {
		log.Printf("missing bot token argument")
	}

	if *chatID == 0 {
		log.Printf("missing chat ID argument")
	}

	bot, err = tgbotapi.NewBotAPI(*botToken)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("%s", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case "answer":
			values := strings.Split(update.Message.Text, ",")
			answerChatID := extractNumber(values[0])
			answerText := strings.TrimSpace(values[1])
			if err != nil {
				log.Printf("%s", err)
			}
			answerMessage(bot, fmt.Sprintf("Відповідь адміністратора: %s", answerText), int64(answerChatID))
		case "start":
			continue
		default:
			forwardMessage(bot, update.Message, *chatID)
			responseMessage(bot, update.Message)
		}
	}
}

func forwardMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, chatID int64) {
	msg := tgbotapi.NewForward(chatID, message.Chat.ID, message.MessageID)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error forwarding message: %v", err)
	}

	forResponseMessage := tgbotapi.NewMessage(chatID,
		fmt.Sprintf(
			`id відправника: [%d].
Імя користувача: <i>%s %s</i> 
Для того щоб відправити повідомлення введіть команду:
"/answer <b>chat_id</b>, ваша відповідь"
Наприклад для відповіді абоненту %s %s 
потрібно ввести:
<code>/answer %d, текст відповіді</code>
<b>Зверни увагу!</b> Знак <b>","</b> є обов'язковим після chat_id`,
			message.Chat.ID,
			message.Chat.FirstName,
			message.Chat.LastName,
			message.Chat.FirstName,
			message.Chat.LastName,
			message.Chat.ID,
		),
	)
	forResponseMessage.ParseMode = "HTML"
	_, err = bot.Send(forResponseMessage)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func responseMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	response := tgbotapi.NewMessage(message.Chat.ID, "Ваше повідомленяя вже обробляється нашим менеджером. Дякуємо за звернення❤")
	_, err := bot.Send(response)
	if err != nil {
		log.Printf("Error forwarding message: %v", err)
	}
}

func answerMessage(bot *tgbotapi.BotAPI, text string, answerChatID int64) {
	answer := tgbotapi.NewMessage(answerChatID, text)
	_, err := bot.Send(answer)
	if err != nil {
		log.Printf("Error forwarding message: %v", err)
	}
}

func extractNumber(text string) int {
	re := regexp.MustCompile(`\d+`)

	match := re.FindString(text)

	num, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}

	return num
}
