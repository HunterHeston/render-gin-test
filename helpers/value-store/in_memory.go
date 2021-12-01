package valuestore

import (
	"math/rand"
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
	return nil, nil
}

func (in InMemory) Save(value []byte) (string, error) {
	// make a copy of the value paseed in by pointer
	data := make([]byte, len(value))
	copy(data, value)

	var strKey string
	for i := 0; i < int(in.maxEntries); i++ {
		key := rand.Intn(int(in.maxEntries))
		strKey = strconv.FormatInt(int64(key), 10)
		_, inUse := in.store[strKey]
		if !inUse {
			in.store[strKey] = value

			break
		}
	}

	return strKey, nil
}

func keyUsed()
