// Copyright 2021 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package garment

import (
	"errors"
	"testing"

	"github.com/franela/goblin"
)

type DatabaseClient struct {
	Connection int
}

func (d *DatabaseClient) Execute() int {
	return d.Connection
}

func (d *DatabaseClient) Terminate() {
	d.Connection = 0
}

func (d *DatabaseClient) Close() {
	d.Connection = 0
}

func (d *DatabaseClient) Reconnect() {
	d.Connection = 1
}

func (d *DatabaseClient) Ping() bool {
	if d.Connection == 1 {
		return true
	}

	return false
}

// TestUnitPool test cases
func TestUnitPool(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("#Pool", func() {
		g.It("It should satisfy all provided test cases", func() {
			pool1 := NewPool()
			pool2 := NewPool()

			g.Assert(pool1.Count()).Equal(0)
			g.Assert(pool2.Count()).Equal(0)

			ping := func(con interface{}) error {
				if con.(*DatabaseClient).Ping() {
					return nil
				}

				return errors.New("DB connection is lost")
			}

			close := func(con interface{}) error {
				con.(*DatabaseClient).Close()

				return nil
			}

			reconnect := func(con interface{}) error {
				con.(*DatabaseClient).Reconnect()

				return nil
			}

			pool1.Set("db", &DatabaseClient{Connection: 1}, ping, close, reconnect)

			// Concurrent Access
			go func() {
				for i := 0; i < 100; i++ {
					if pool1.Has("db") && pool1.Ping("db") == nil {
						g.Assert(pool1.Get("db").(*DatabaseClient).Execute()).Equal(1)
					}
				}
			}()

			go func() {
				for i := 0; i < 100; i++ {
					if pool1.Has("db") && pool1.Ping("db") == nil {
						g.Assert(pool1.Get("db").(*DatabaseClient).Execute()).Equal(1)
					}
				}
			}()

			g.Assert(pool1.Count()).Equal(1)
			g.Assert(pool2.Count()).Equal(1)

			g.Assert(pool1.Has("db")).Equal(true)
			g.Assert(pool2.Has("db")).Equal(true)

			g.Assert(pool1.Has("db1")).Equal(false)
			g.Assert(pool2.Has("db1")).Equal(false)

			g.Assert(pool1.Get("db").(*DatabaseClient).Execute()).Equal(1)
			g.Assert(pool2.Get("db").(*DatabaseClient).Execute()).Equal(1)

			// Terminate Connection
			pool2.Get("db").(*DatabaseClient).Terminate()
			g.Assert(pool1.Get("db").(*DatabaseClient).Execute()).Equal(0)
			pool1.Reconnect("db")
			g.Assert(pool1.Get("db").(*DatabaseClient).Execute()).Equal(1)
			pool2.Remove("db")
			pool1.Reconnect("db")
			pool2.Remove("db")
			pool2.Ping("db")

			g.Assert(pool1.Has("db")).Equal(false)
			g.Assert(pool2.Has("db")).Equal(false)

			g.Assert(pool1.Count()).Equal(0)
			g.Assert(pool2.Count()).Equal(0)

			pool2.Set("db", &DatabaseClient{Connection: 1}, ping, close, reconnect)

			g.Assert(pool1.Count()).Equal(1)
			g.Assert(pool2.Count()).Equal(1)

			g.Assert(pool1.Ping("db") == nil).Equal(true)
			g.Assert(pool2.Ping("db") == nil).Equal(true)

			pool1.Close("db")
			pool1.Close("db1")

			g.Assert(pool1.Ping("db") == nil).Equal(false)
			g.Assert(pool2.Ping("db") == nil).Equal(false)

			g.Assert(pool1.Count()).Equal(1)
			g.Assert(pool2.Count()).Equal(1)
		})
	})
}
