package pkg

const (
	TreeOrder = 5
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

// BNode b树
// b树所有子树的高度相同
// 每个节点至少有 (m+1)/2 个关键字
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
		// 查找插入位置
		for idx = 0; idx < b.Len; idx++ {
			if cmp(b.Keys[idx], key) > 0 {
				break
			}
		}
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

// find 递归查找应该插入的位置
func (b *BNode) find(key interface{}, cmp Cmp) (*BNode, int) {
	// 第一个一定是非空的
	i := 0
	// 假设keys数组是连续而且有序的的
	for {
		// 找到了
		if cmp(key, b.Keys[i]) == 0 {
			return b, i
		}
		// 比当前小
		if cmp(key, b.Keys[i]) < 0 {
			if b.Children[i] != nil {
				return b.Children[i].find(key, cmp)
			}
			return b, i
		}

		if i+1 < b.Len {
			i++
			continue
		}

		if b.Children[i+1] == nil {
			return b, i + 1
		}
		return b.Children[i+1].find(key, cmp)
	}
}
