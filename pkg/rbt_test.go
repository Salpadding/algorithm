package pkg

import (
	"fmt"
	"testing"
)

func TestRBT(t *testing.T) {
	tree := NewRBTree(func(x, y interface{}) int {
		return x.(int) - y.(int)
	})

	var i int
	for i = 0; i < 10240; i++ {
		tree.Set(i, fmt.Sprintf("%d", i))
	}

	fmt.Println()
}
