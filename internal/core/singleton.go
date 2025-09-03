package core

import (
	"sync"
)

var instance *InMemoryBus
var once sync.Once

func NewInMemoryBus() *InMemoryBus {
	once.Do(func() {
		instance = newInMemoryBus()
	})
	return instance
}
