#pragma once
#define RED 0
#define BLACK 1

struct rb_node {
    struct rb_node* parent;
    struct rb_node* left;
    struct rb_node* right;
    int color;
    int key;
    int value;
};

// 左旋
void rotate_left(struct rb_node* node) {
    if (! node->right) return;
    struct rb_node* y = node->right;
    y->parent = node->parent;
    if (node->parent) {
        if (node == node->parent->left) {
            y->parent->left = y;
        } else {
            y->parent->right = y;
        }
    }
    struct rb_node* yl = y->left;
    node->parent = y;
    y->left = node;
    node->right = yl;
    if (yl) {
        yl->parent = node;
    }
}

// 右旋
void rotate_right(struct rb_node* node) {
    if (!node->left) return;
    struct rb_node* y = node->left;
    y->parent = node->parent;
    if (node->parent) {
        if (node == node->parent->left) {
            y->parent->left = y;
        } else {
            y->parent->right = y;
        }
    }

    struct rb_node* yr = y->right;
    node->parent = y;
    y->right = node;
    node->left = yr;
    if (yr) {
        yr->parent = node;
    }
}

// 修复红黑树
void fix_insert(struct rb_node* node, struct rb_node** root) {
    if (!node->parent) {
        node->color = BLACK;
        return;
    }

    if (node->parent->color == BLACK) return;

    struct rb_node* parent = node->parent;

    struct rb_node* uncle;
    if (parent == parent->parent->left) {
        uncle = parent->parent->right;
    } else {
        uncle = parent->parent->left;
    }

    if (uncle && uncle->color == RED) {
        parent->color = BLACK;
        uncle->color = BLACK;
        parent->parent->color = RED;
        return fix_insert(parent->parent, root);
    }

    if (parent == parent->parent->left) {
        // LR -> LL
        if (node == parent->right) {
            rotate_left(parent);
            parent = node;
        }
        if (parent->parent == *root) {
            *root = parent;
        }
        parent->color = BLACK;
        parent->parent->color = RED;
        rotate_right(parent->parent);
    } else {
        if (node == parent->left) {
            rotate_right(parent);
            parent = node;
        }
        if (parent->parent == *root) {
            *root = parent;
        }        
        parent->color = BLACK;
        parent->parent->color = RED;
        rotate_left(parent->parent);
    }
}

// 查找
struct rb_node* find(struct rb_node* node, int key, struct rb_node** last) {
    if (!node) return NULL;
    if (last) *last = node;
    if (key == node->key) return node;
    if (key < node->key) {
        return find(node->left, key, last);
    }
    return find(node->right, key, last);
}

void insert(struct rb_node** root, int key, int value) {
    if (!(*root)) {
        *root = calloc(1, sizeof(struct rb_node));
        (*root)->color = BLACK;
        (*root)->key = key;
        (*root)->value = value;
        return;
    }
    struct rb_node* parent;
    struct rb_node* found = find(*root, key, &parent);
    if (found) {
        found->value = value;
        return;
    }

    struct rb_node* new_node = calloc(1, sizeof(struct rb_node));
    new_node->key = key;
    new_node->value = value;
    new_node->parent = parent;
    if (key < parent->key) {
        parent->left = new_node;
    } else {
        parent->right = new_node;
    }
    fix_insert(new_node, root);
}

void print_tree(struct rb_node* n) {
    if (!n) {
        return;
    }
    print_tree(n->left);
    printf("%d %d\n", n->key, n->value);
    print_tree(n->right);
}

/**
 * Note: The returned array must be malloced, assume caller calls free().
 */
int* twoSum(int* nums, int numsSize, int target, int* returnSize){
    struct rb_node* root;
    struct rb_node* found;
    for(int i = 0; i < numsSize; i++) {
        found = find(root, nums[i], NULL);
        if (found) {
            int* ret = calloc(2, sizeof(int));
            ret[0] = found->value;
            ret[1] = i;
            *returnSize= 2;
            return ret;
        }
        insert(&root, target - nums[i], i);
    }
    return NULL;
}


