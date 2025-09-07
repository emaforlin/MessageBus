package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type MessageHandler func(msg string) error

type Topic struct {
	handlers []MessageHandler
}

type InMemoryBus struct {
	mu     *sync.RWMutex
	topics map[string]Topic
}

func (b *InMemoryBus) Subscribe(topicName string, handler MessageHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if topic, exists := b.topics[topicName]; exists {
		topic.handlers = append(topic.handlers, handler)
	} else {
		b.topics[topicName] = Topic{
			handlers: []MessageHandler{handler},
		}
	}

	return nil
}

func (b *InMemoryBus) Publish(ctx context.Context, topicName string, message string) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topic, exists := b.topics[topicName]
	if !exists {
		return fmt.Errorf("topic %q does not exists", topicName)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(topic.handlers))

	for _, handlerFn := range topic.handlers {
		wg.Add(1)
		go func(handler MessageHandler) {
			defer wg.Done()
			err := handler(message)
			if err != nil {
				fmt.Printf("handling error: %v\n", err)
				errCh <- err
			}
		}(handlerFn)
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

func newInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		mu:     &sync.RWMutex{},
		topics: make(map[string]Topic, 8),
	}
}
