package inmemory

import (
	"fmt"
	"strconv"
)

type InMemory struct {
	store      map[string][]byte
	maxEntries int64
}

func NewInMemory(maxEntries int64) InMemory {
	return InMemory{
		store:      make(map[string][]byte),
		maxEntries: maxEntries,
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
	// make a copy of the value paseed in by pointer
	data := make([]byte, len(value))
	copy(data, value)

	var strKey string
	for i := 0; i < int(in.maxEntries); i++ {
		strKey = strconv.FormatInt(int64(i), 10)
		_, inUse := in.store[strKey]
		if !inUse {
			in.store[strKey] = value
			break
		}
	}

	return strKey, nil
}
