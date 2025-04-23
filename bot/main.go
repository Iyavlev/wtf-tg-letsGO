package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "os"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Message struct {
    Message string `json:"message"`
}

func main() {
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("BOT_TOKEN is not set")
    }

    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Авторизован: %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        text := update.Message.Text
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сообщение отправлено на сервер.")
        bot.Send(msg)

        go sendToServer(text)
    }
}

func sendToServer(message string) {
    payload := Message{Message: message}
    data, _ := json.Marshal(payload)

    // ВАЖНО: адрес сервиса в Kubernetes (через ClusterIP или сервис)
    resp, err := http.Post("http://wtf-api/posts", "application/json", bytes.NewBuffer(data))
    if err != nil {
        log.Println("Ошибка отправки:", err)
        return
    }
    defer resp.Body.Close()

    log.Println("Отправлено:", message)
}
