package core

import (
	"context"
	"fmt"
	"slices"
	"testing"
)

func TestPublish(t *testing.T) {
	var topicName = "#test"
	bus := newInMemoryBus()

	spyHandlerFunc := MessageHandler(func(msg string) error {
		fmt.Printf("SPY HANDLER: %v", msg)
		return nil
	})

	bus.Subscribe(topicName, spyHandlerFunc)

	err := bus.Publish(context.Background(), topicName, "ping")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestSubscribe(t *testing.T) {
	var topicName = "#test"

	bus := newInMemoryBus()

	messagesToSend := []string{"tic", "tac", "toe"}
	sentMessages := []string{}

	spyHandlerFunc := MessageHandler(func(msg string) error {
		sentMessages = append(sentMessages, msg)
		return nil
	})

	bus.Subscribe(topicName, spyHandlerFunc)

	for _, msg := range messagesToSend {
		bus.Publish(context.Background(), topicName, msg)
	}

	if !slices.Equal(messagesToSend, sentMessages) {
		t.Errorf("got %v but expected %v", sentMessages, messagesToSend)
	}
}
