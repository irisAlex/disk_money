package api

import (
	"fmt"
	"money/pkg/prese"
	"time"

	"github.com/gin-gonic/gin"
)

type ResChatData struct {
	MessageList *ChatMessage `json:"message_list"`
}

type ChatMessage struct {
	IsMe    bool   `json:"isMe"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

type RequestData struct {
	NewMessage string `json:"newMessage"`
}

func Chat(c *gin.Context) {

	//var list []*ChatMessage
	var req RequestData
	prese.ParseJSON(c, &req)

	fmt.Println(req.NewMessage)

	currentTime := time.Now().Format("15:04")

	prese.ResJSON(c, 200, &ResChatData{
		MessageList: &ChatMessage{
			IsMe:    false,
			Content: "我是一个人工智能",
			Time:    currentTime,
		},
	})
}
