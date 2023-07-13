#include <common.hh>
#include <iostream>

using namespace std;

class avl_tree_node {
  public:
    avl_tree_node *parent;

    avl_tree_node *left_child;
    avl_tree_node *right_child;

    int val;
    int depth;

    avl_tree_node(int val)
        : val(val), depth(1), left_child(NULL), right_child(NULL),
          parent(NULL) {}

    // 插入节点
    avl_tree_node *insert(int val);

    // 是否是叶子节点
    bool is_leaf() {
        return this->left_child == NULL && this->right_child == NULL;
    }

    // 遍历计算 到达叶子结点的最长路径长度
    int find_depth();

    // 右旋 保持二叉排序树的性质
    void right_rotate();
    // 左旋 保持二叉排序树的性质
    void left_rotate();

    int left_child_depth() {
        return this->left_child ? this->left_child->depth : 0;
    }
    int right_child_depth() {
        return this->right_child ? this->right_child->depth : 0;
    }
};

void avl_tree_node::right_rotate() {
    if (!this->left_child)
        return;

    // a 是最小不平衡树
    // a 左边比较重 进行右旋
    // a 变成 b 右子树
    // b 从 a 的左子树分离
    // b 的右子树变成 a 的左子树
    avl_tree_node *b = this->left_child;
    avl_tree_node *a = this;

    avl_tree_node *b_right = b->right_child;
    avl_tree_node *a_parent = a->parent;
    b->right_child = a;
    a->parent = b;

    b->parent = a_parent;
    if (a_parent)
        a_parent->left_child = b;

    a->left_child = b_right;
    if (b_right)
        b_right->parent = a;

    // 重新计算 a 的深度
    a->depth = MAX(a->left_child_depth(), a->right_child_depth()) + 1;
    // 重新计算 b 的深度
    b->depth = MAX(b->left_child_depth(), b->right_child_depth()) + 1;
}

void avl_tree_node::left_rotate() {
    if (!this->right_child)
        return;

    // a 是最小不平衡树
    // a 右边比较重 进行左旋
    // a 变成 b 左子树
    // b 从 a 的右子树分离
    // b 的左子树变成 a 的右子树
    avl_tree_node *b = this->right_child;
    avl_tree_node *a = this;

    avl_tree_node *b_left = b->left_child;
    avl_tree_node *a_parent = a->parent;
    b->left_child = a;
    a->parent = b;

    b->parent = a_parent;
    if (a_parent)
        a_parent->right_child = b;

    a->right_child = b_left;
    if (b_left)
        b_left->parent = a;

    // 重新计算 a 的深度
    a->depth = MAX(a->left_child_depth(), a->right_child_depth()) + 1;
    // 重新计算 b 的深度
    b->depth = MAX(b->left_child_depth(), b->right_child_depth()) + 1;
}

avl_tree_node *avl_tree_node::insert(int val) {
    avl_tree_node *current = this;
    avl_tree_node **slot;
    avl_tree_node *tmp;

    // 查找合适的插入地点
    while (true) {
        if (current->val == val)
            return this;

        if (current->val < val) {
            if (current->right_child == NULL) {
                slot = &current->right_child;
                goto insert;
            }
            current = current->right_child;
        } else {
            if (current->left_child == NULL) {
                slot = &current->left_child;
                goto insert;
            }
            current = current->left_child;
        }
    }

insert:
        tmp = current;
        while (tmp) {
            tmp->depth = MAX(tmp->left_child_depth(), tmp->right_child_depth()) + 1;
            tmp = tmp->parent;
        }

    // 插入新的节点
    avl_tree_node *new_node = new avl_tree_node(val);
    new_node->parent = current;
    *slot = new_node;

    // 查找最小不平衡子树
    tmp = current->parent;

    while (tmp) {
        int left_depth = tmp->left_child_depth();
        int right_depth = tmp->right_child_depth();

        if (ABS(left_depth - right_depth) > 1)
            break;
        tmp = tmp->parent;
    }

    if (!tmp)
        return this;

    // 维持平衡
    int left_depth = tmp->left_child_depth();
    int right_depth = tmp->right_child_depth();

    if (left_depth > right_depth) {
        // ll
        if (tmp->left_child->left_child_depth() >
            tmp->left_child->right_child_depth()) {
        } else {
            // lr -> ll
            tmp->left_child->left_rotate();
        }
        avl_tree_node *new_root = this;
        if (tmp == this) {
            new_root = tmp->left_child;
        }
        tmp->right_rotate();
        return new_root;
    } else {
        // rr
        if (tmp->right_child->right_child_depth() >
            tmp->right_child->right_child_depth()) {
        } else {
            // rl -> rr
            tmp->right_child->right_rotate();
        }
        avl_tree_node *new_root = this;
        if (tmp == this) {
            new_root = tmp->right_child;
        }
        tmp->left_rotate();
        return new_root;
    }
}

int avl_tree_node::find_depth() {
    int left_depth = this->left_child_depth();
    int right_depth = this->right_child_depth();
    return 1 + MAX(left_depth, right_depth);
}
