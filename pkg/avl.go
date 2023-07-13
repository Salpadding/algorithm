package pkg

import "fmt"

// AVLNode AVL 树的节点
type AVLNode struct {
	Depth      int      // 节点深度 到叶子结点的最长路径
	Key        int      // 关键词 用于比较
	Parent     *AVLNode // 父节点
	LeftChild  *AVLNode // 左子树
	RightChild *AVLNode // 右子树
}

// LeftRotate 左旋
// 当右子树的右子树比较重的时候
// 把当前节点加入右子树的左子树 以此达到平衡效果
func (a *AVLNode) LeftRotate() {
	parent := a.Parent
	b := a.RightChild //
	if b == nil {
		return
	}

	left := b.LeftChild

	// a 变成 b 左子树
	b.LeftChild = a
	a.Parent = b
	b.Parent = parent
	if parent != nil {
		if a.Key > parent.Key {
			parent.RightChild = b
		} else {
			parent.LeftChild = b
		}
	}

	// b 原来的左子树变成 a 的右子树
	a.RightChild = left
	if left != nil {
		left.Parent = a
	}

	// 重新计算 a, b 的深度
	// 先更新子节点的深度
	a.updateDepth()
	b.updateDepth()
}

// RightRotate 右旋
// 当左子树的左子树比较重的时候
// 把当前节点加入左子树的右子树 以此达到平衡效果
func (a *AVLNode) RightRotate() {
	parent := a.Parent
	b := a.LeftChild //
	if b == nil {
		return
	}

	right := b.RightChild

	// a 变成 b 右子树
	b.RightChild = a
	a.Parent = b
	b.Parent = parent
	if parent != nil {
		if a.Key > parent.Key {
			parent.RightChild = b
		} else {
			parent.LeftChild = b
		}
	}

	// b 原来的右子树变成 a 的左子树
	a.LeftChild = right
	if right != nil {
		right.Parent = a
	}

	a.updateDepth()
	b.updateDepth()
}

func (a *AVLNode) updateDepth() {
	// 重新计算 a, b 的深度
	a.Depth = func() int {
		if a.leftChildDepth() < a.rightChildDepth() {
			return a.rightChildDepth()
		}
		return a.leftChildDepth()
	}() + 1
}

// leftChildDepth 左子树的深度
func (a *AVLNode) leftChildDepth() int {
	if a.LeftChild == nil {
		return 0
	}
	return a.LeftChild.Depth
}

// rightChildDepth 右子树的深度
func (a *AVLNode) rightChildDepth() int {
	if a.RightChild == nil {
		return 0
	}
	return a.RightChild.Depth
}

// insert 把 key 插入当前树
// 返回新插入的节点 不作平衡
// 已存在则 return nil
func (a *AVLNode) insert(key int) *AVLNode {
	if a.Key == key {
		return nil
	}

	// 查找树的定义就是左子树比跟节点小 所以插入左子树
	if key < a.Key {
		if a.LeftChild == nil {
			a.LeftChild = &AVLNode{
				Depth:  1,
				Key:    key,
				Parent: a,
			}
			return a.LeftChild
		}
		return a.LeftChild.insert(key)
	}
	if a.RightChild == nil {
		a.RightChild = &AVLNode{
			Depth:  1,
			Key:    key,
			Parent: a,
		}
		return a.RightChild
	}
	return a.RightChild.insert(key)

}

// balanceFactor 平衡因子
func (a *AVLNode) balanceFactor() int {
	b := a.leftChildDepth() - a.rightChildDepth()
	if b < 0 {
		return -b
	}
	return b
}

func (a *AVLNode) Insert(key int) (newRoot *AVLNode) {
	var leastUnbalanced *AVLNode
	newRoot = a
	leaf := a.insert(key)
	// key 已存在
	if leaf == nil {
		return
	}

	// 递归更新深度
	current := leaf

	// 只有新节点的祖先节点需要更新深度
	for current != nil {
		current.updateDepth()
		if leastUnbalanced == nil && current.balanceFactor() > 1 {
			leastUnbalanced = current
		}
		current = current.Parent
	}

	// 不需要调整
	if leastUnbalanced == nil {
		return
	}

	fmt.Printf("balancing least unbalanced = %d\n", leastUnbalanced.Key)
	if leastUnbalanced.leftChildDepth() > leastUnbalanced.rightChildDepth() {
		// lr 通过以此左旋 转化成 ll
		if leastUnbalanced.LeftChild.rightChildDepth() > leastUnbalanced.LeftChild.leftChildDepth() {
			leastUnbalanced.LeftChild.LeftRotate()
			fmt.Println("lr")
		} else {
			fmt.Println("ll")
		}
		if leastUnbalanced == a {
			newRoot = a.LeftChild
		}
		leastUnbalanced.RightRotate()
	} else {
		// rl 通过右旋 转化成 rr
		if leastUnbalanced.RightChild.leftChildDepth() > leastUnbalanced.RightChild.rightChildDepth() {
			fmt.Println("rl")
			leastUnbalanced.RightChild.RightRotate()
		} else {
			fmt.Println("rr")
		}
		if leastUnbalanced == a {
			newRoot = a.RightChild
		}
		leastUnbalanced.LeftRotate()
	}
	return
}
