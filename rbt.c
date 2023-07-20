#include <stdio.h>
#include <stdlib.h>

#define RED 0
#define BLACK 1

struct rb_node {
    struct rb_node *parent;
    struct rb_node *left;
    struct rb_node *right;
    int color;
    int key;
    int value;
};

void print_tree(struct rb_node *root);

// 左旋
void rotate_left(struct rb_node *node, struct rb_node **root) {
    if (!node->right)
        return;
    struct rb_node *y = node->right;
    y->parent = node->parent;
    if (node->parent) {
        if (node == node->parent->left) {
            y->parent->left = y;
        } else {
            y->parent->right = y;
        }
    } else {
        *root = y;
    }
    struct rb_node *yl = y->left;
    node->parent = y;
    y->left = node;
    node->right = yl;
    if (yl) {
        yl->parent = node;
    }
}

// 右旋
void rotate_right(struct rb_node *node, struct rb_node **root) {
    if (!node->left)
        return;
    struct rb_node *y = node->left;
    y->parent = node->parent;
    if (node->parent) {
        if (node == node->parent->left) {
            y->parent->left = y;
        } else {
            y->parent->right = y;
        }
    } else {
        *root = y;
    }

    struct rb_node *yr = y->right;
    node->parent = y;
    y->right = node;
    node->left = yr;
    if (yr) {
        yr->parent = node;
    }
}

// 插入操作后 修复红黑树
void fix_insert(struct rb_node *node, struct rb_node **root) {
    struct rb_node *parent;
    struct rb_node *uncle;
retry:
    if (!node->parent) {
        node->color = BLACK;
        return;
    }

    if (node->parent->color == BLACK)
        return;

    parent = node->parent;

    if (parent == parent->parent->left) {
        uncle = parent->parent->right;
    } else {
        uncle = parent->parent->left;
    }

    // uncle 节点是红色向上递归
    if (uncle && uncle->color == RED) {
        parent->color = BLACK;
        uncle->color = BLACK;
        parent->parent->color = RED;
        node = parent->parent;
        goto retry;
    }

    if (parent == parent->parent->left) {
        // LR -> LL
        if (node == parent->right) {
            rotate_left(parent, root);
            parent = node;
        }
        if (parent->parent == *root) {
            *root = parent;
        }
        parent->color = BLACK;
        parent->parent->color = RED;
        rotate_right(parent->parent, root);
    } else {
        if (node == parent->left) {
            rotate_right(parent, root);
            parent = node;
        }
        if (parent->parent == *root) {
            *root = parent;
        }
        parent->color = BLACK;
        parent->parent->color = RED;
        rotate_left(parent->parent, root);
    }
}

// 查找
struct rb_node *find(struct rb_node *node, int key, struct rb_node **last) {
    while (node) {
        if (last)
            *last = node;
        if (key == node->key)
            return node;
        if (key < node->key)
            node = node->left;
        node = node->right;
    }
    return node;
}

void insert(struct rb_node **root, int key, int value) {
    struct rb_node *parent;
    struct rb_node *found;
    if (!(*root)) {
        *root = calloc(1, sizeof(struct rb_node));
        (*root)->color = BLACK;
        (*root)->key = key;
        (*root)->value = value;
        return;
    }
    found = find(*root, key, &parent);
    if (found) {
        found->value = value;
        return;
    }

#define new_node found
    new_node = calloc(1, sizeof(struct rb_node));
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

void replace_by(struct rb_node *x, struct rb_node *y, struct rb_node **root) {
    if (!(x->parent)) {
        *root = y;
        if (y)
            y->parent = NULL;
        return;
    }
    if (x == x->parent->left) {
        x->parent->left = y;
        if (y)
            y->parent = x->parent;
    } else {
        x->parent->right = y;
        if (y)
            y->parent = x->parent;
    }
}

void delete_node(struct rb_node *n, struct rb_node **root) {
    struct rb_node *x;
    struct rb_node *w;

    if (n->left && n->right) {
        x = n->left;
        while (x->right)
            x = x->right;

        n->key = x->key;
        n->value = x->value;
        n = x;
    }

    if (n->color == RED) {
        if (n->left) {
            replace_by(n, n->left, root);
        } else {
            replace_by(n, n->right, root);
        }
        return;
    }

    if (n->left && n->left->color == RED) {
        replace_by(n, n->left, root);
        n->left->color = BLACK;
        return;
    }

    if (n->right && n->right->color == RED) {
        replace_by(n, n->right, root);
        n->right->color = BLACK;
        return;
    }
    struct rb_node **xx;
    struct rb_node *p;

    p = n->parent;
    if (!p) {
        *root = NULL;
        return;
    }
    if (n == n->parent->left) {
        xx = &p->left;
    } else {
        xx = &p->right;
    }
    replace_by(n, n->left, root);

#define x xx
#define do_direct(DIR1, DIR2, LABEL)                                           \
    w = p->DIR2;                                                               \
    if (w && w->color == RED) {                                                \
        p->color = RED;                                                        \
        w->color = BLACK;                                                      \
        rotate_##DIR1(p, root);                                                \
        w = p->DIR2;                                                           \
    }                                                                          \
    LABEL:                                                                     \
    if (w && w->DIR2 && w->DIR2->color == RED) {                               \
        w->color = p->color;                                                   \
        w->DIR2->color = BLACK;                                                \
        p->color = BLACK;                                                      \
        rotate_##DIR1(p, root);                                                \
        return;                                                                \
    }                                                                          \
    if (w && w->DIR1 && w->DIR1->color == RED) {                               \
        w->color = RED;                                                        \
        w->DIR1->color = BLACK;                                                \
        rotate_##DIR2(w, root);                                                \
        w = w->parent;                                                         \
        goto LABEL;                                                            \
    }                                                                          \
    if (p->color == RED) {                                                     \
        p->color = BLACK;                                                      \
        w->color = RED;                                                        \
        return;                                                                \
    }                                                                          \
    w->color = RED;                                                            \
    if (!p->parent) {                                                          \
        p->color = BLACK;                                                      \
        return;                                                                \
    }                                                                          \
    if (p == p->parent->left)                                                  \
        x = &p->left;                                                          \
    else                                                                       \
        x = &p->right;                                                         \
    p = p->parent;                                                             \
    goto fix_up;

fix_up:
    // 走到根节点了
    if (x == &p->left) {
        do_direct(left, right, RR)
    } else {
        do_direct(right, left, LL)
    }
}

void delete(int key, struct rb_node **root) {
    if (!(*root))
        return;
    struct rb_node *found = find(*root, key, NULL);
    if (!found)
        return;
    delete_node(found, root);
}

/*
digraph G {
    graph [ratio=.48];
    node [style=filled, color=black, shape=circle, width=.6
        fontname=Helvetica, fontweight=bold, fontcolor=white,
        fontsize=24, fixedsize=true];

    6, 8, 17, 22, 27
    [fillcolor=red];

    n1, n2, n3, n4, n5, n6, n7, n8, n9, n10, n11
    [label="NIL", shape=record, width=.4,height=.25, fontsize=16];

    13 -> 8, 17;
    8 -> 1 [weight=6];
    8 -> 11 [weight=5];
    17 -> 15 [weight=4];
    17 -> 25 [weight=5];
    1 -> n1 [weight=7];
    1 -> 6;
    11 -> n4 [weight=10];
    11 -> n5 [weight=14];
    6 -> n2, n3;
    15 -> n6 [weight=14];
    15 -> n7 [weight=10];
    25 -> 22;
    25 -> 27 [weight=6];
    22 -> n8 [weight=5];
    22 -> n9 [weight=3];
    27 -> n10 [weight=3];
    27 -> n11 [weight=5];
}
 */

#define N 16

// 打印红黑树 格式= Graphviz
void print_tree(struct rb_node *root) {
    puts("digraph G { \n\
graph [ratio=.48]; node [style=filled, color=black, shape=circle, width=.6 \n\
fontname=Helvetica, fontweight=bold, fontcolor=white, fontsize=24, fixedsize=true];");
    // 广度优先遍历
    struct rb_node **queue = calloc(N * 2, sizeof(void *));
    int *red_nodes = calloc(N * 2, sizeof(int));
    int red_cnt = 0;
    int head = 0;
    int tail = 0;
    int nil_cnt = 0;

#define size (tail - head)
#define append(x)                                                              \
    do {                                                                       \
        queue[tail] = x;                                                       \
        tail++;                                                                \
    } while (0);
#define dequeue()                                                              \
    ({                                                                         \
        struct rb_node *first = queue[head];                                   \
        head++;                                                                \
        first;                                                                 \
    })

    char buf1[128];
    char buf2[128];

#define fmt(x, buf)                                                            \
    {                                                                          \
        if (x) {                                                               \
            sprintf(buf, "%d", x->key);                                        \
        } else {                                                               \
            sprintf(buf, "n%d", nil_cnt++);                                    \
        }                                                                      \
    }
    if (!root)
        return;

    append(root);

    while (size > 0) {
        struct rb_node *x = dequeue();
        if (x->color == RED) {
            red_nodes[red_cnt++] = x->key;
        }
        fmt(x->left, buf1);
        fmt(x->right, buf2);

        printf("%d -> %s, %s;\n", x->key, buf1, buf2);
        if (x->left)
            append(x->left);
        if (x->right)
            append(x->right);
    }

    int i;
    for (i = 0; i < nil_cnt; i++) {
        printf("n%d %s", i, i == nil_cnt - 1 ? "\n" : ",");
    }
    puts("[label=\"NIL\", shape=record, width=.4,height=.25, fontsize=16];\n");

    if (red_cnt > 0) {
        for (i = 0; i < red_cnt; i++) {
            printf("%d %s", red_nodes[i], i == red_cnt - 1 ? "\n" : ",");
        }
        puts("[fillcolor=red]\n");
    }

    puts("}\n");
}

/**
 * Note: The returned array must be malloced, assume caller calls free().
 */
int *twoSum(int *nums, int numsSize, int target, int *returnSize) {
    struct rb_node *root;
    struct rb_node *found;
    for (int i = 0; i < numsSize; i++) {
        found = find(root, nums[i], NULL);
        if (found) {
            int *ret = calloc(2, sizeof(int));
            ret[0] = found->value;
            ret[1] = i;
            *returnSize = 2;
            return ret;
        }
        insert(&root, target - nums[i], i);
    }
    return NULL;
}

int main() {
    int i = 0;
    struct rb_node *root = NULL;
    for (i = 0; i < N; i++) {
        insert(&root, i, i);
    }
    for (i = N - 1; i >= 0; i--) {
        delete (i, &root);
        print_tree(root);
        fflush(stdout);
    }
    return 0;
}
