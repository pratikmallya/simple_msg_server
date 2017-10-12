package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPostMsgRoute(t *testing.T) {
	testCases := []struct {
		postMessages         []Message
		expectedResponse     []ReturnMessage
		expectedResponseCode int
		waitDuration         int
	}{
		{postMessages: []Message{},
			expectedResponse:     nil,
			expectedResponseCode: 200,
		},
		{
			postMessages: []Message{
				Message{
					Username: "Spongebob",
					Text:     "I'm Ready!",
				},
			},
			waitDuration: 0,
			expectedResponse: []ReturnMessage{
				ReturnMessage{
					Text: "I'm Ready!",
				},
			},
		},
	}

	for _, testCase := range testCases {
		for _, msg := range testCase.postMessages {
			b, err := json.Marshal(msg)
			assert.NoErr(t, err)
			req := httptest.NewRequest("POST", "/chat", b)
		}
	}
}
