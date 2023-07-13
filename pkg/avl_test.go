package pkg

import (
	"testing"
)

func TestAVLTree(t *testing.T) {
	root := &AVLNode{
		Key:   0,
		Depth: 1,
	}

	for i := 1; i < 32; i++ {
		root = root.Insert(i)
	}

}
