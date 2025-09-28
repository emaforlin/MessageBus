package core

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type MessageHandler func(msg string) error

// MessageBus defines the core messaging interface
type MessageBus interface {
	Subscribe(topicName string, handler MessageHandler) (string, error)
	Unsubscribe(topicName, subscriptionID string) error
	Publish(ctx context.Context, topicName string, message string) error
	Close() error
}

// TopicManager handles topic lifecycle and subscription management
type TopicManager interface {
	AddSubscription(topicName string, sub Subscription) error
	RemoveSubscription(topicName, subscriptionID string) error
	GetSubscriptions(topicName string) ([]Subscription, error)
	ListTopics() []string
	DeleteTopic(topicName string) error
}

// SubscriptionStore manages subscription storage and retrieval
type SubscriptionStore interface {
	Store(topicName string, subscription Subscription) error
	Retrieve(topicName, subscriptionID string) (Subscription, error)
	Remove(topicName, subscriptionID string) error
	List(topicName string) ([]Subscription, error)
}

// MessagePublisher handles message publishing logic
type MessagePublisher interface {
	PublishToTopic(ctx context.Context, topic *Topic, message string) error
}

type Topic struct {
	subscriptions []Subscription
	mu            sync.RWMutex // Per-topic locking for better concurrency
}

type Subscription struct {
	id      string
	handler MessageHandler
}

// InMemoryBus implementation
type InMemoryBus struct {
	mu        *sync.RWMutex
	topics    map[string]*Topic
	publisher MessagePublisher
	closed    bool
}

// Ensure InMemoryBus implements all interfaces
var (
	_ MessageBus        = (*InMemoryBus)(nil)
	_ TopicManager      = (*InMemoryBus)(nil)
	_ SubscriptionStore = (*InMemoryBus)(nil)
	_ MessagePublisher  = (*InMemoryBus)(nil)
)

// Subscribe implements MessageBus interface
func (b *InMemoryBus) Subscribe(topicName string, handler MessageHandler) (string, error) {
	if b.closed {
		return "", fmt.Errorf("message bus is closed")
	}

	id := uuid.NewString()
	sub := Subscription{
		id:      id,
		handler: handler,
	}

	return id, b.AddSubscription(topicName, sub)
}

// AddSubscription implements TopicManager interface
func (b *InMemoryBus) AddSubscription(topicName string, sub Subscription) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return fmt.Errorf("message bus is closed")
	}

	topic, exists := b.topics[topicName]
	if !exists {
		topic = &Topic{subscriptions: make([]Subscription, 0)}
		b.topics[topicName] = topic
	}

	topic.mu.Lock()
	topic.subscriptions = append(topic.subscriptions, sub)
	topic.mu.Unlock()

	return nil
}

// Publish implements MessageBus interface
func (b *InMemoryBus) Publish(ctx context.Context, topicName string, message string) error {
	if b.closed {
		return fmt.Errorf("message bus is closed")
	}

	b.mu.RLock()
	topic, exists := b.topics[topicName]
	b.mu.RUnlock()

	if !exists {
		return fmt.Errorf("topic %q does not exist", topicName)
	}

	return b.PublishToTopic(ctx, topic, message)
}

// PublishToTopic implements MessagePublisher interface
func (b *InMemoryBus) PublishToTopic(ctx context.Context, topic *Topic, message string) error {
	topic.mu.RLock()
	subscriptions := make([]Subscription, len(topic.subscriptions))
	copy(subscriptions, topic.subscriptions)
	topic.mu.RUnlock()

	if len(subscriptions) == 0 {
		return nil // No subscribers, no error
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(subscriptions))

	for _, sub := range subscriptions {
		wg.Add(1)
		go func(handler MessageHandler) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
			default:
				if err := handler(message); err != nil {
					fmt.Printf("handler error: %v\n", err)
					errCh <- err
				}
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

// Unsubscribe implements MessageBus interface
func (b *InMemoryBus) Unsubscribe(topicName, subscriptionID string) error {
	return b.RemoveSubscription(topicName, subscriptionID)
}

// RemoveSubscription implements TopicManager interface
func (b *InMemoryBus) RemoveSubscription(topicName, subscriptionID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	topic, exists := b.topics[topicName]
	if !exists {
		return fmt.Errorf("topic %q does not exist", topicName)
	}

	topic.mu.Lock()
	defer topic.mu.Unlock()

	for i, sub := range topic.subscriptions {
		if sub.id == subscriptionID {
			topic.subscriptions = append(topic.subscriptions[:i], topic.subscriptions[i+1:]...)
			break
		}
	}

	// Clean up empty topics
	if len(topic.subscriptions) == 0 {
		delete(b.topics, topicName)
	}

	return nil
}

// GetSubscriptions implements TopicManager interface
func (b *InMemoryBus) GetSubscriptions(topicName string) ([]Subscription, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topic, exists := b.topics[topicName]
	if !exists {
		return nil, fmt.Errorf("topic %q does not exist", topicName)
	}

	topic.mu.RLock()
	defer topic.mu.RUnlock()

	// Return a copy to prevent external modification
	subscriptions := make([]Subscription, len(topic.subscriptions))
	copy(subscriptions, topic.subscriptions)
	return subscriptions, nil
}

// ListTopics implements TopicManager interface
func (b *InMemoryBus) ListTopics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topics := make([]string, 0, len(b.topics))
	for topicName := range b.topics {
		topics = append(topics, topicName)
	}
	return topics
}

// DeleteTopic implements TopicManager interface
func (b *InMemoryBus) DeleteTopic(topicName string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.topics[topicName]; !exists {
		return fmt.Errorf("topic %q does not exist", topicName)
	}

	delete(b.topics, topicName)
	return nil
}

// Store implements SubscriptionStore interface
func (b *InMemoryBus) Store(topicName string, subscription Subscription) error {
	return b.AddSubscription(topicName, subscription)
}

// Retrieve implements SubscriptionStore interface
func (b *InMemoryBus) Retrieve(topicName, subscriptionID string) (Subscription, error) {
	subscriptions, err := b.GetSubscriptions(topicName)
	if err != nil {
		return Subscription{}, err
	}

	for _, sub := range subscriptions {
		if sub.id == subscriptionID {
			return sub, nil
		}
	}

	return Subscription{}, fmt.Errorf("subscription %q not found in topic %q", subscriptionID, topicName)
}

// Remove implements SubscriptionStore interface
func (b *InMemoryBus) Remove(topicName, subscriptionID string) error {
	return b.RemoveSubscription(topicName, subscriptionID)
}

// List implements SubscriptionStore interface
func (b *InMemoryBus) List(topicName string) ([]Subscription, error) {
	return b.GetSubscriptions(topicName)
}

// Close implements MessageBus interface
func (b *InMemoryBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return fmt.Errorf("message bus already closed")
	}

	b.closed = true
	b.topics = nil
	return nil
}

func NewMessageBus() *InMemoryBus {
	bus := &InMemoryBus{
		mu:     &sync.RWMutex{},
		topics: make(map[string]*Topic, 8),
		closed: false,
	}
	bus.publisher = bus // Self-reference for publisher interface
	return bus
}
