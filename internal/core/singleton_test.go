package core

import (
	"testing"
)

func TestNewInMemoryBusSingleton(t *testing.T) {
	bus1 := NewInMemoryBus()
	bus2 := NewInMemoryBus()

	if bus1 != bus2 {
		t.Errorf("Should exist only one instance")
	}
}
