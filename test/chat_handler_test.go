package test

import (
	"encoding/json"
	"justn0w-bot/config"
	"justn0w-bot/internal/request"
	"justn0w-bot/internal/router"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	config.Init()
}

func TestChatHandler(t *testing.T) {
	r := router.Router()
	w := httptest.NewRecorder()

	chatRequest := request.ChatRequest{
		Question: "第四条职业道德内容是什么呢",
	}
	requestStr, _ := json.Marshal(chatRequest)
	req, _ := http.NewRequest("POST", "/chat/generate/stream", strings.NewReader(string(requestStr)))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInVzZXJOYW1lIjoianVzdG5vdyIsImlzcyI6Imp1c3RuMHctYm90Iiwic3ViIjoiMSIsImV4cCI6MTc1NzY4NTk4NSwibmJmIjoxNzU3Njc4Nzg1LCJpYXQiOjE3NTc2Nzg3ODV9.QSbxAwKSjvIgxvIS9DARp6--Y7rWZqEQCvj1EbI5DEY")

	r.ServeHTTP(w, req)
}
