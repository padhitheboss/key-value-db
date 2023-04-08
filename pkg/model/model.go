package model

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Command represents a parsed command received via the REST API.
type Command struct {
	Name string
	Args []string
}

// Datastore represents an in-memory key-value datastore.
type Datastore struct {
	mutex sync.RWMutex
	data  map[string]Data
}

// Data represents the value and metadata associated with a key.
type Data struct {
	value     string
	expiry    time.Time
	condition string
}

type Response struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
	Value  string `json:"value,omitempty"`
}

func (store *Datastore) Set(key string, value string, expiryTime time.Time, condition string) (string, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	currentValue, exists := store.data[key]

	if exists && !currentValue.expiry.IsZero() && currentValue.expiry.Before(time.Now()) {
		delete(store.data, key)
	}
	if exists && currentValue.condition == "NX" {
		return "", fmt.Errorf("key already exists")
	}

	if !exists && currentValue.condition == "XX" {
		return "", fmt.Errorf("key does not exist")
	}

	if !expiryTime.IsZero() && expiryTime.Before(time.Now()) {
		return "", fmt.Errorf("expiry time cannot be in the past")
	}

	store.data[key] = Data{
		value:     value,
		expiry:    expiryTime,
		condition: condition,
	}

	return "successful", nil
}

func (store *Datastore) Get(key string) (string, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	data, exists := store.data[key]

	if !exists {
		return "", fmt.Errorf("key not found")
	}

	if !data.expiry.IsZero() && data.expiry.Before(time.Now()) {
		delete(store.data, key)
		return "", fmt.Errorf("key not found")
	}

	return data.value, nil
}

func (store *Datastore) QPush(key string, values []string) (string, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	queue, exists := store.data[key]

	if !exists {
		queue = Data{value: ""}
		queue.value += strings.Join(values, ",")
	} else {
		queue.value += "," + strings.Join(values, ",")
	}

	store.data[key] = queue
	return "successful", nil
}

func (store *Datastore) QPop(key string) (string, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	queue, exists := store.data[key]
	if !exists {
		return "", fmt.Errorf("queue is empty")
	}

	values := strings.Split(queue.value, ",")

	if len(values) == 0 {
		return "", fmt.Errorf("queue is empty")
	}

	value := values[0]

	if len(values) > 1 {
		queue.value = strings.Join(values[1:], ",")
		store.data[key] = queue
	} else {
		delete(store.data, key)
	}

	return value, nil
}

var DB = Datastore{
	mutex: sync.RWMutex{},
	data:  make(map[string]Data),
}
