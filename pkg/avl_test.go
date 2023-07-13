package pkg

import (
	"fmt"
	"testing"
)

func TestAVLTree(t *testing.T) {
	root := &AVLNode{
		Key:   0,
		Depth: 1,
	}

	for i := 1; i < 1024; i++ {
		root = root.Insert(i)
	}

	fmt.Printf("left depth = %d right depth = %d", root.leftChildDepth(), root.rightChildDepth())
}
