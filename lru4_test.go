// Author: dup2X
// Last modified: 2016-07-22 14:50
// Filename: lru4_test.go
package lru4

import (
	"fmt"
	"testing"
)

func TestLRU4(t *testing.T) {
	c := New(10)
	for i := 0; i < 14; i++ {
		c.Add(fmt.Sprintf("key-%d", i), i)
	}
	for i := 0; i < 14; i++ {
		val, ok := c.Get(fmt.Sprintf("key-%d", i))
		if ok {
			println(val.(int))
		} else {
			println("miss", i)
		}
	}
}
