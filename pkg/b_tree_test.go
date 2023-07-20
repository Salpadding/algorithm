package pkg

import (
	"fmt"
	"testing"
)

func TestBTree(t *testing.T) {
	cmp := func(x, y interface{}) int {
		return x.(int) - y.(int)
	}
	btree := &BTree{
		cmp: cmp,
	}

	var i int
	for i = 0; i < 10240; i++ {
		if i == 8 {
			fmt.Println()
		}
		btree.Insert(i)
	}

	fmt.Println()
}
