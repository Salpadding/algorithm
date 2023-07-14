package pkg

type Cmp func(x, y interface{}) int

// RBTree
// 红黑树可以这样定义(简洁版):
// 1. 红黑树是二叉排序树 根节点是黑节点 nil节点也是黑节点
// 2. 红节点的父节点是黑节点
// 3. 对于任意两个不同的 nil节点, 从nil节点到根节点的路径上出现的黑节点数量相同
// 通过定义可以证明
// 1. 红黑树中任意一个黑节点以及它的左右子树也构成红黑树
// 2. 如果根节点是黑色 根节点的左右是红色 我把左右都染成黑色 这依然是一颗红黑树
// 3. 红节点一定有父亲节点
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

// rightRotate 右旋
// 向右移动节点同时保持二叉排序树的性质
func (y *RBNode) rightRotate() {
	x := y.L
	if x == nil {
		return
	}
	beta := x.R

	x.P = y.P
	if x.P != nil {
		if y == x.P.L {
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

// leftRotate 右旋
// 向左移动节点同时保持二叉排序树的性质
// 如果对一个红节点左旋 不会改变黑高度
func (y *RBNode) leftRotate() {
	x := y.R
	if x == nil {
		return
	}
	beta := x.L

	x.P = y.P
	if x.P != nil {
		if y == x.P.L {
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

	leaf.insertFix(&newRoot)
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

// simpleDelete 删除当前节点
// 删除黑色节点可能会造成黑高度-1 需要进行平衡
//
// 首先要保证二叉排序树的性质
// 1. 左右子树有且至少有一个是 nil  把那个比较高的左/右子树接到 parent 上面
// 2. (一次递归) 左右子树都不是 nil 把前驱节点的数值移植到当前节点 转化成删前驱节点 前驱节点的右子树一定是nil 所以满足1
//
// 在满足1的条件下 继续保证红黑树的性质
// case1. 删除的节点是红色 那么删除不会改变黑高度 不需要修复
// case2. 删除节点是黑色 但是这个节点一边是 nil 一边是 红色 因为这个红色节点会接替删除的位置 只要把这个红色节点染黑就行了
// case3. 删除节点的 uncle 节点是红色 把这个红色的节点移过来再染黑
func (y *RBNode) delete(root **RBNode) {
	// 左右均不为 nil
	// 转化为删前驱
	if y.L != nil && y.R != nil {
		cur := y.L
		for cur.R != nil {
			cur = cur.R
		}

		y.Key = cur.Key
		y.Value = cur.Value
		cur.delete(root)
		return
	}

	// 删除红色节点
	// 不会出现红-红
	if y.Color == y.red() {
		if y.L != nil {
			y.replaceBy(y.L)
		} else {
			y.replaceBy(y.R)
		}
		return
	}

	// 有红色左子
	// 染黑后 不可能出现红-红
	if y.L != nil && y.L.Color == y.red() {
		y.replaceBy(y.L)
		y.L.Color = y.black()
		if y == *root {
			*root = y.L
		}
		return
	}

	// 有红色右子
	// 染黑后 不可能出现红-红
	if y.R != nil && y.R.Color == y.red() {
		y.replaceBy(y.R)
		y.R.Color = y.black()
		if y == *root {
			*root = y.R
		}
		return
	}

	// 删的是黑节点 而且不是根节点
	// 这种情况先不考虑
	if y.P == nil {
		panic("delete black child of root")
	}

	uncle := y.P.L
	if y == uncle {
		uncle = y.P.R
	}

	// 下面讨论 uncle, uncle左, uncle右 = ()

	//  (黑,红,黑) (黑,红,红) (黑,黑,红)

}

func (x *RBNode) replaceBy(child *RBNode) {
	child.P = x.P
	if x.P != nil {
		if x == x.P.L {
			x.P.L = child
		} else {
			x.P.R = child
		}
	}
}

// insertFix 插入完成后修复红黑树
// 红黑树插入后可能出现 红-红 的情况
//
// 如何修复?
// case1. 如果另一条分支上是黑-黑 那么我们可以把红色节点移动到这两个黑色节点的中间
// 为了在移动的同时保持二叉排序树的性质 我们会采用 旋转 + 染色的操作
//
// case2. 另一条分支上是黑-红
// 我们 parent + uncle 从红的变成黑的 把 uncle 的 parent 从黑的变成红的
// 这样这两条path的 黑高度都是不变的, 并且缩小了问题规模
// 顺着往上递归 一定能
// 1. 到达根节点 把根节点染黑 根节点黑高度 + 1 结束递归
// 2. 遇到红-黑 结束递归
// 3. 跳到 case1 结束递归
func (y *RBNode) insertFix(root **RBNode) {
	// y.P == nil 递归到根节点了 染黑
	// 根节点 以及根节点左右都黑了 根节点黑高度+1
	if y.P == nil {
		y.Color = y.black()
		return
	}
	// 红-黑
	if y.P.Color == y.black() {
		return
	}

	// 红-红
	// 在左边
	if y.P == y.P.P.L {
		uncle := y.P.P.R
		// case2
		if uncle != nil && uncle.Color == y.red() {
			y.P.Color = y.black()
			uncle.Color = y.black()
			y.P.P.Color = y.red()
			y.P.P.insertFix(root)
			return
		}
		// case1
		grandpa := y.P.P
		// LR -> LL
		if y == y.P.R {
			y.P.leftRotate()
		}

		if *root == grandpa {
			*root = grandpa.L
		}
		grandpa.rightRotate()
		grandpa.Color = y.red()
		grandpa.P.Color = y.black()
		return
	}

	// 在右边
	uncle := y.P.P.L
	if uncle != nil && uncle.Color == y.red() {
		y.P.Color = y.black()
		uncle.Color = y.black()
		y.P.P.Color = y.red()
		y.P.P.insertFix(root)
		return
	}
	grandpa := y.P.P
	// RL -> RR
	if y == y.P.L {
		y.P.rightRotate()
	}

	if *root == grandpa {
		*root = grandpa.R
	}
	grandpa.leftRotate()
	grandpa.Color = y.red()
	grandpa.P.Color = y.black()
}
