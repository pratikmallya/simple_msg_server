package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMsgRoute(t *testing.T) {
	testCases := []struct {
		postMessages         []Message
		username             string
		expectedResponse     string
		expectedResponseCode int
		waitDuration         time.Duration
	}{
		{postMessages: []Message{},
			expectedResponse:     "",
			expectedResponseCode: 404,
		},
		{
			postMessages: []Message{
				Message{
					Username: "Spongebob",
					Text:     "I'm Ready!",
				},
			},
			username:             "Spongebob",
			waitDuration:         0 * time.Second,
			expectedResponseCode: 200,
			expectedResponse:     "I'm Ready!",
		},
	}

	for _, testCase := range testCases {
		engine := InitializeEngine()
		for _, msg := range testCase.postMessages {
			b, err := json.Marshal(msg)
			assert.NoError(t, err)
			req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(b))
			rr := httptest.NewRecorder()
			engine.ServeHTTP(rr, req)
		}

		time.Sleep(testCase.waitDuration)

		// Get response
		req := httptest.NewRequest("GET", fmt.Sprintf("/chat/%s", testCase.username), nil)
		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, testCase.expectedResponseCode)
		assert.Contains(t, rr.Body.String(), testCase.expectedResponse)
	}
}

func TestPostMsgRoute(t *testing.T) {
	testCases := []struct {
		description          string
		postMessages         []Message
		expectedResponseCode []int
	}{
		{
			description:          "empty post",
			postMessages:         []Message{},
			expectedResponseCode: []int{400},
		},
		{
			description: "multiple messages",
			postMessages: []Message{
				Message{
					Username: "Spongebob",
					Text:     "I'm Ready!",
				},
				Message{},
				Message{
					Username: "Spongebob",
					Text:     "I'm Ready!",
					Timeout:  600,
				},
			},
			expectedResponseCode: []int{201, 400, 201},
		},
		{
			description: "bad timeout",
			postMessages: []Message{
				Message{
					Username: "Spongebob",
					Text:     "I'm Ready!",
					Timeout:  -1,
				},
			},
			expectedResponseCode: []int{400},
		},
	}

	for _, testCase := range testCases {
		engine := InitializeEngine()
		for i, msg := range testCase.postMessages {
			b, err := json.Marshal(msg)
			assert.NoError(t, err)
			req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(b))
			rr := httptest.NewRecorder()
			engine.ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, testCase.expectedResponseCode[i])
		}
	}
}

func TestPerfWrite(t *testing.T) {
	numreqs := 100
	var reqs []*http.Request
	for i := 0; i < numreqs; i++ {
		msg := Message{
			Username: "Spongebob",
			Text:     "I'm Ready!",
			Timeout:  60,
		}
		b, err := json.Marshal(msg)
		assert.NoError(t, err)
		req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(b))
		reqs = append(reqs, req)
	}

	engine := InitializeEngine()

	for i := 0; i < numreqs; i++ {
		rr := httptest.NewRecorder()
		engine.ServeHTTP(rr, reqs[i])
		assert.Equal(t, 201, rr.Code)
	}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chat/Spongebob", nil)
	engine.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
}
