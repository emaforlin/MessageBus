package core

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type MessageHandler func(msg string) error

type Topic struct {
	subscriptions []Subscription
}

type InMemoryBus struct {
	mu     *sync.RWMutex
	topics map[string]*Topic
}

type Subscription struct {
	id      string
	handler MessageHandler
}

func (b *InMemoryBus) Subscribe(topicName string, handler MessageHandler) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	id := uuid.NewString()

	sub := Subscription{
		id:      id,
		handler: handler,
	}

	if topic, exists := b.topics[topicName]; exists {
		topic.subscriptions = append(topic.subscriptions, sub)
	} else {
		b.topics[topicName] = &Topic{
			subscriptions: []Subscription{sub},
		}
	}

	return id, nil
}

func (b *InMemoryBus) Publish(ctx context.Context, topicName string, message string) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topic, exists := b.topics[topicName]
	if !exists {
		return fmt.Errorf("topic %q does not exists", topicName)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(topic.subscriptions))

	for _, sub := range topic.subscriptions {
		wg.Add(1)
		go func(handler MessageHandler) {
			defer wg.Done()
			err := handler(message)
			if err != nil {
				fmt.Printf("handling error: %v\n", err)
				errCh <- err
			}
		}(sub.handler)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var handlingErrs []error
	for err := range errCh {
		handlingErrs = append(handlingErrs, err)
	}

	return errors.Join(handlingErrs...)
}

func (b *InMemoryBus) Unsubscribe(topicName, subscriptionID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	topic, exists := b.topics[topicName]
	if !exists {
		return fmt.Errorf("topic %q does not exist", topicName)
	}

	for i, sub := range topic.subscriptions {
		if sub.id == subscriptionID {
			topic.subscriptions = append(topic.subscriptions[:i], topic.subscriptions[i+1:]...)

			if len(topic.subscriptions) == 0 {
				delete(b.topics, topicName)
			}

			return nil
		}
	}

	return fmt.Errorf("subscription %q not found in topic %q", subscriptionID, topicName)
}

func newInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		mu:     &sync.RWMutex{},
		topics: make(map[string]*Topic, 8),
	}
}
