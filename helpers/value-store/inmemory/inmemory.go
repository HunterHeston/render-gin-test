package inmemory

import (
	"fmt"

	stringgeneration "github.com/hunterheston/gin-server/helpers/string-generation"
)

type InMemory struct {
	store map[string][]byte
}

func NewInMemory() InMemory {
	return InMemory{
		store: make(map[string][]byte),
	}
}

func (im InMemory) Lookup(id string) ([]byte, error) {
	v, ok := im.store[id]
	if !ok {
		return nil, fmt.Errorf("value for key %q  does not exist", id)
	}
	return v, nil
}

// Generate an ID and save the value to an in memory store.
func (in InMemory) Save(value []byte) (string, error) {
	// make a copy of the value paseed in by pointer.
	data := make([]byte, len(value))
	copy(data, value)

	// Generate random strings of 6 chars until one does not exist.
	randomID := stringgeneration.RandStringBytesRmndr(6)
	for _, exists := in.store[randomID]; exists; {
		randomID = stringgeneration.RandStringBytesRmndr(6)
	}

	// save the string that ended up not existing.
	in.store[randomID] = value

	// return the id(string).
	return randomID, nil
}
