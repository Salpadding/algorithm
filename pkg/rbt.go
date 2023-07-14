package pkg

type Cmp func(x, y interface{}) int

type RBTree struct {
	compareFn func(x, y interface{}) int
	root      *RBNode
}

func NewRBTree(cmp Cmp) *RBTree {
	return &RBTree{
		compareFn: cmp,
	}
}

func (t *RBTree) Set(key, value interface{}) {
	if t.root == nil {
		t.root = &RBNode{
			Color: 1,
			Key:   key,
			Value: value,
		}
		return
	}
	t.root = t.root.Insert(key, value, t.compareFn)
}

func (t *RBTree) Get(key interface{}) interface{} {
	node := t.root.find(key, t.compareFn)
	if node == nil {
		return nil
	}
	return node.Value
}

type RBNode struct {
	P     *RBNode
	L     *RBNode
	R     *RBNode
	Color int
	Key   interface{}
	Value interface{}
}

func (n *RBNode) red() int {
	return 0
}

func (n *RBNode) black() int {
	return 1
}

func (y *RBNode) rightRotate(cmp Cmp) {
	x := y.L
	if x == nil {
		return
	}
	beta := x.R

	x.P = y.P
	if x.P != nil {
		if cmp(y.Key, x.Key) < 0 {
			x.P.L = x
		} else {
			x.P.R = x
		}
	}

	y.P = x
	x.R = y

	y.L = beta
	if beta != nil {
		beta.P = y
	}
}

func (y *RBNode) leftRotate(cmp Cmp) {
	x := y.R
	if x == nil {
		return
	}
	beta := x.L

	x.P = y.P
	if x.P != nil {
		if cmp(y.Key, x.Key) < 0 {
			x.P.L = x
		} else {
			x.P.R = x
		}
	}

	y.P = x
	x.L = y

	y.R = beta
	if beta != nil {
		beta.P = y
	}
}

func (y *RBNode) find(key interface{}, cmp Cmp) *RBNode {
	if cmp(y.Key, key) == 0 {
		return y
	}
	if cmp(key, y.Key) < 0 {
		if y.L == nil {
			return nil
		}
		return y.L.find(key, cmp)
	}
	if y.R == nil {
		return nil
	}
	return y.R.find(key, cmp)
}

func (y *RBNode) Insert(key interface{}, value interface{}, cmp Cmp) (newRoot *RBNode) {
	newRoot = y
	leaf := y.insert(key, value, cmp)
	if leaf == nil {
		return
	}

	leaf.fixup(&newRoot, cmp)
	newRoot.Color = y.black()
	return
}

func (y *RBNode) insert(key interface{}, value interface{}, cmp Cmp) *RBNode {
	if cmp(y.Key, key) == 0 {
		y.Value = value
		return nil
	}
	if cmp(key, y.Key) < 0 {
		if y.L == nil {
			y.L = &RBNode{
				P:     y,
				Key:   key,
				Value: value,
				Color: y.red(),
			}
			return y.L
		}
		return y.L.insert(key, value, cmp)
	}
	if y.R == nil {
		y.R = &RBNode{
			P:     y,
			Key:   key,
			Value: value,
			Color: y.red(),
		}
		return y.R
	}
	return y.R.insert(key, value, cmp)
}

func (y *RBNode) fixup(root **RBNode, cmp Cmp) {
	if y.P == nil || y.P.Color == y.black() {
		return
	}

	// 在左边
	if cmp(y.P.Key, y.P.P.Key) < 0 {
		uncle := y.P.P.R
		if uncle != nil && uncle.Color == y.red() {
			y.P.Color = y.black()
			uncle.Color = y.black()
			y.P.P.Color = y.red()
			y.P.P.fixup(root, cmp)
			return
		}
		grandpa := y.P.P
		// LR -> LL
		if cmp(y.Key, y.P.Key) > 0 {
			y.P.leftRotate(cmp)
		}

		if *root == grandpa {
			*root = grandpa.L
		}
		grandpa.rightRotate(cmp)
		grandpa.Color = y.red()
		grandpa.P.Color = y.black()
		return
	}

	uncle := y.P.P.L
	if uncle != nil && uncle.Color == y.red() {
		y.P.Color = y.black()
		uncle.Color = y.black()
		y.P.P.Color = y.red()
		y.P.P.fixup(root, cmp)
		return
	}
	grandpa := y.P.P
	// RL -> RR
	if cmp(y.Key, y.P.Key) < 0 {
		y.P.rightRotate(cmp)
	}

	if *root == grandpa {
		*root = grandpa.R
	}
	grandpa.leftRotate(cmp)
	grandpa.Color = y.red()
	grandpa.P.Color = y.black()
}
