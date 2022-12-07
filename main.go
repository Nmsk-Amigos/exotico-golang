package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func TelegramBot() {

	bot, err := tgbotapi.NewBotAPI(goDotEnvVariable("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			DiscordSend(update.Message.From.UserName, update.Message.Text)
			bot.Send(msg)
		}
	}
}

func main() {
	// Echo instance
	e := echo.New()
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Exotico")
	})
	// Start server
	go TelegramBot()
	log.Println("http://localhost:8080")

	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func DiscordSend(username, text string) {
	var jsonData = []byte(`{
		"username": "` + username + `",
		"avatar_url": "https://media.discordapp.net/attachments/1049891409878069248/1049914341589270538/2048px-Telegram_logo.png",
		"content": "` + text + `"

	}`)
	fmt.Println(bytes.NewBuffer(jsonData))
	req, err := http.NewRequest("POST", goDotEnvVariable("DISCORD_WEBHOOK"), bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	response, error := client.Do(req)
	if error != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}
