package services

import (
	"sync"

	"shifumi/internal/models"
)

type gameEventBroker struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan models.GameEvent]struct{}
}

func newGameEventBroker() *gameEventBroker {
	return &gameEventBroker{
		subscribers: make(map[string]map[chan models.GameEvent]struct{}),
	}
}

func (b *gameEventBroker) Subscribe(gameID string) chan models.GameEvent {
	ch := make(chan models.GameEvent, 8)

	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.subscribers[gameID]; !ok {
		b.subscribers[gameID] = make(map[chan models.GameEvent]struct{})
	}
	b.subscribers[gameID][ch] = struct{}{}

	return ch
}

func (b *gameEventBroker) Unsubscribe(gameID string, ch chan models.GameEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subscribers, ok := b.subscribers[gameID]
	if !ok {
		return
	}

	if _, exists := subscribers[ch]; exists {
		delete(subscribers, ch)
		close(ch)
	}

	if len(subscribers) == 0 {
		delete(b.subscribers, gameID)
	}
}

func (b *gameEventBroker) Publish(gameID string, event models.GameEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.subscribers[gameID] {
		select {
		case ch <- event:
		default:
		}
	}
}
