package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Message ...
type Message struct {
	Username       string `json:"username" binding:"required"`
	Text           string `json:"text" binding:"required"`
	Timeout        int    `json:"timeout"`
	ID             int
	ExpirationTime time.Time
}

// ReturnMessage ...
type ReturnMessage struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var m map[string][]Message

func main() {
	m = make(map[string][]Message)
	rand.Seed(42)
	r := gin.Default()

	r.GET("/chat/:username", func(c *gin.Context) {
		now := time.Now()
		username := c.Param("username")

		if _, ok := m[username]; !ok {
			c.String(http.StatusNotFound, "No unexpired messages found for %s", username)
			return
		}
		retmsgs := returnUnexpiredMessages(m[username], now)
		if retmsgs == nil {
			c.String(http.StatusNotFound, "No unexpired messages found for %s", username)
			return
		}
		delete(m, username)
		c.JSON(200, retmsgs)
	})

	r.POST("/chat", func(c *gin.Context) {
		var msg Message
		now := time.Now()
		err := c.BindJSON(&msg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if msg.Timeout < 0 {
			c.String(http.StatusBadRequest, "Invalid timeout: %d ", msg.Timeout)
			return
		}

		if msg.Timeout == 0 {
			msg.Timeout = 60
		}

		msg.ExpirationTime = now.Add(time.Duration(msg.Timeout) * time.Second)
		msg.ID = rand.Int()

		if msgs, ok := m[msg.Username]; ok {
			m[msg.Username] = append(msgs, msg)
			c.JSON(http.StatusCreated, gin.H{
				"id": msg.ID,
			})
			return
		}
		m[msg.Username] = []Message{msg}
		c.JSON(http.StatusCreated, gin.H{
			"id": msg.ID,
		})
		return
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

func returnUnexpiredMessages(msgs []Message, now time.Time) []ReturnMessage {
	var rmsgs []ReturnMessage
	for _, msg := range msgs {
		if now.After(msg.ExpirationTime) {
			continue
		}
		rmsgs = append(rmsgs, ReturnMessage{
			ID:   msg.ID,
			Text: msg.Text,
		})

	}
	return rmsgs
}
