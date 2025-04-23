package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "sync"
)

type Message struct {
    Message string `json:"message"`
}

var (
    messages []string
    mu       sync.Mutex
)

func postsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read body", http.StatusInternalServerError)
        return
    }

    var msg Message
    if err := json.Unmarshal(body, &msg); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    mu.Lock()
    messages = append(messages, msg.Message)
    mu.Unlock()

    log.Printf("Новое сообщение: %s", msg.Message)
    fmt.Fprintln(w, "OK")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintln(w, "<html><body><h1>Сообщения</h1><ul>")
    for _, msg := range messages {
        fmt.Fprintf(w, "<li>%s</li>", msg)
    }
    fmt.Fprintln(w, "</ul></body></html>")
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/posts", postsHandler)

    log.Println("Сервер запущен на порту 80")
    log.Fatal(http.ListenAndServe(":80", nil))
}
