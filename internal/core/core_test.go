package core

import (
	"context"
	"fmt"
	"slices"
	"sync"
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

func TestPublishMultipleSubscribers(t *testing.T) {
	bus := newInMemoryBus()

	var mu sync.Mutex
	var received []string

	handler1 := func(msg string) error {
		mu.Lock()
		received = append(received, "handler1:"+msg)
		mu.Unlock()
		return nil
	}

	handler2 := func(msg string) error {
		mu.Lock()
		received = append(received, "handler2:"+msg)
		mu.Unlock()
		return nil
	}

	topic := "test-topic"
	bus.Subscribe(topic, handler1)
	bus.Subscribe(topic, handler2)

	msg := "hello world"
	err := bus.Publish(context.Background(), topic, msg)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// Wait for goroutines to finish
	// (goroutines are joined in Publish, so this is safe)
	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Errorf("Expected 2 messages received, got %d", len(received))
	}

	found1, found2 := false, false
	for _, r := range received {
		if r == "handler1:"+msg {
			found1 = true
		}
		if r == "handler2:"+msg {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Errorf("Both handlers should have received the message")
	}
}
