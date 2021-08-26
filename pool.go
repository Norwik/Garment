// Copyright 2021 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package garment

import (
	"fmt"
	"sync"
)

var (
	// pool var
	pool *Pool
	// once var
	once sync.Once
)

// PingCallback type
type PingCallback func(interface{}) error

// CloseCallback type
type CloseCallback func(interface{}) error

// ReconnectCallback type
type ReconnectCallback func(interface{}) error

// Connection type
type Connection struct {
	Value     interface{}
	Ping      PingCallback
	Close     CloseCallback
	Reconnect ReconnectCallback
}

// Pool type
type Pool struct {
	sync.RWMutex

	items map[string]Connection
}

// NewPool creates a new instance of Pool
func NewPool() *Pool {
	once.Do(func() {
		pool = &Pool{
			items: make(map[string]Connection),
		}
	})

	return pool
}

// Get gets a connection with a key
func (c *Pool) Get(key string) interface{} {
	c.Lock()
	defer c.Unlock()

	connection, ok := c.items[key]

	if !ok {
		return nil
	}

	return connection.Value
}

// Has checks if connection exists
func (c *Pool) Has(key string) bool {
	if val := c.Get(key); val != nil {
		return true
	}

	return false
}

// Set sets a connection with a key
func (c *Pool) Set(key string, value interface{}, ping PingCallback, close CloseCallback, reconnect ReconnectCallback) {
	c.Lock()
	defer c.Unlock()

	c.items[key] = Connection{
		Value:     value,
		Ping:      ping,
		Close:     close,
		Reconnect: reconnect,
	}
}

// Close close the connection and removes the key
func (c *Pool) Close(key string) error {
	c.Lock()
	defer c.Unlock()

	connection, ok := c.items[key]

	if !ok {
		return fmt.Errorf("Unable to find %s", key)
	}

	return connection.Close(connection.Value)
}

// Ping checks the connection status
func (c *Pool) Ping(key string) error {
	c.Lock()
	defer c.Unlock()

	connection, ok := c.items[key]

	if !ok {
		return fmt.Errorf("Unable to find %s", key)
	}

	return connection.Ping(connection.Value)
}

// Reconnect reconnects again
func (c *Pool) Reconnect(key string) error {
	c.Lock()
	defer c.Unlock()

	connection, ok := c.items[key]

	if !ok {
		return fmt.Errorf("Unable to find %s", key)
	}

	return connection.Reconnect(connection.Value)
}

// Remove removes a connection
func (c *Pool) Remove(key string) {
	c.Lock()
	defer c.Unlock()

	_, ok := c.items[key]

	if !ok {
		return
	}

	delete(c.items, key)
}

// Count counts the number of managed connections
func (c *Pool) Count() int {
	c.Lock()
	defer c.Unlock()

	return len(c.items)
}
