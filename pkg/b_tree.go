package pkg

const (
	TreeOrder = 3
)

// BTree B树 TreeOrder=3时是2-3树
// 2-3 树的性质
// 1. 以顺序方式插入 2^k-1个节点 结果是一个满二叉树
type BTree struct {
	root *BNode
	cmp  Cmp
}

func (b *BNode) minKeys() int {
	return (TreeOrder+1)/2 - 1
}

func (b *BTree) Insert(key interface{}) {
	if b.root == nil {
		b.root = b.root.New()
		b.root.addKey(key, 0, b.cmp)
		return
	}
	b.root = b.root.insert(key, b.cmp)
}

// BNode b树定义
// m 阶 b树定义
// 1. 根节点所有子树的高度相同
// 2. 如果不是根节点 每个节点至少有 (m+1)/2 - 1个关键字 至多 m-1 个关键字 m 个子树
// 通过 1,2 可以证明 b 树的任意子树也是 b 树
type BNode struct {
	Parent *BNode
	// 长度是 m 因为可能插入需要临时加入一个节点
	Keys []interface{}
	// 长度是 m+1
	Children []*BNode
	Len      int
}

func (b *BNode) addKey(key interface{}, idx int, cmp Cmp) int {
	// idx < 0 时候说明是下面的节点分裂 然后插进来的
	// 需要调整子树的位置
	if idx < 0 {
		idx = b.findSlot(key, cmp) + 1
		for i := b.Len; i >= idx+1; i-- {
			b.Children[i] = b.Children[i-1]
		}
		b.Children[idx] = nil
	}
	for i := b.Len - 1; i >= idx+1; i-- {
		// 往右移动
		b.Keys[i] = b.Keys[i-1]
	}
	b.Keys[idx] = key
	b.Len++
	return idx
}

func (b *BNode) insert(key interface{}, cmp Cmp) (root *BNode) {
	root = b
	node, target := b.find(key, cmp)
	if target < node.Len && cmp(node.Keys[target], key) == 0 {
		return
	}
	node.addKey(key, target, cmp)
	node.insertFix(&root, cmp)
	return
}

func (b *BNode) insertFix(root **BNode, cmp Cmp) {
	if b.Len < TreeOrder {
		return
	}
	// 进行分裂

	parent := b.Parent
	if parent == nil {
		parent = b.New()
		*root = parent
	}
	b.split(parent, cmp)
	parent.insertFix(root, cmp)
}

func (b *BNode) New() *BNode {
	return &BNode{
		Keys:     make([]interface{}, TreeOrder),
		Children: make([]*BNode, TreeOrder+1),
		Len:      0,
	}
}

func (b *BNode) delete(root **RBNode, idx int) {
	if b.Children[idx] != nil {
		cur := b.Children[idx]
		// 找前驱 转化为删除前驱
		for cur.Children[cur.Len] != nil {
			cur = cur.Children[cur.Len]
		}
		b.Keys[idx] = cur.Keys[cur.Len-1]
		cur.delete(root, cur.Len-1)
		return
	}

	for i := idx; i < b.Len-1; i++ {
		b.Keys[i] = b.Keys[i+1]
	}
	b.Len--

	// 根节点
	if b.Parent == nil {
		return
	}

	if b.Len >= b.minKeys() {
		return
	}
}

func (b *BNode) rotate() {

}

func (b *BNode) fixDelete() {
	if b.Parent == nil || b.Len >= b.minKeys() {
		return
	}

	// 找最合适的 sib
	var i int
	// TODO: 二分查找
	for i = 0; i <= b.Parent.Len; i++ {
		if b.Children[i] == b {
			break
		}
	}

}

// split 分裂
// b 树保持性质的核心操作
func (b *BNode) split(parent *BNode, cmp Cmp) {
	idx := parent.addKey(b.Keys[TreeOrder/2], -1, cmp)
	left := b.New()
	left.Len = TreeOrder / 2
	copy(left.Keys[:left.Len], b.Keys[:left.Len])
	copy(left.Children[:TreeOrder/2+1], b.Children[:TreeOrder/2+1])
	for i := range left.Children[:TreeOrder/2+1] {
		if left.Children[i] != nil {
			left.Children[i].Parent = left
		}
	}

	right := b.New()
	right.Len = TreeOrder - left.Len - 1
	copy(right.Keys[:right.Len], b.Keys[left.Len+1:TreeOrder])
	copy(right.Children[:right.Len+1], b.Children[left.Len+1:])
	for i := range right.Children[:TreeOrder/2+1] {
		if right.Children[i] != nil {
			right.Children[i].Parent = right
		}
	}

	left.Parent = parent
	right.Parent = parent
	parent.Children[idx] = left
	parent.Children[idx+1] = right
}

// findSlot 二分查找
// 默认 b.keys[b.Len] = 正无穷
// b.keys[i] <= key < b.keys[i+1]
func (b *BNode) findSlot(key interface{}, cmp Cmp) int {
	i := -1
	j := b.Len - 1
	for {
		mid := (i + j + 1) / 2
		// 找到了
		if mid > 0 && mid < b.Len && cmp(key, b.Keys[mid]) == 0 {
			return mid
		}

		if i == j {
			return i
		}

		// 比当前小
		if mid >= 0 && (mid == b.Len || cmp(key, b.Keys[mid]) < 0) {
			j = mid - 1
		} else {
			i = mid
		}
	}
}

// find 递归查找应该插入的位置
func (b *BNode) find(key interface{}, cmp Cmp) (*BNode, int) {
	idx := b.findSlot(key, cmp)
	if idx >= 0 && cmp(b.Keys[idx], key) == 0 {
		return b, idx
	}

	if b.Children[idx+1] != nil {
		return b.Children[idx+1].find(key, cmp)
	}
	return b, idx + 1
}
