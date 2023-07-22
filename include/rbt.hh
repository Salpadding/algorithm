#pragma once
#include <cstdio>
#include <queue>

using namespace std;

class rb_node {
    static int red() { return 0; }
    static int black() { return 1; }

  public:
    rb_node *parent;
    rb_node *left;
    rb_node *right;
    int color;
    int key;
    int value;
    int size;

    void rotate_left();
    void rotate_right();

    rb_node *find(int key, rb_node **last) {
        rb_node *cur = this;

        while (cur) {
            if (last)
                *last = cur;
            if (key == cur->key) {
                return cur;
            }
            if (key < cur->key) {
                cur = cur->left;
            } else {
                cur = cur->right;
            }
        }
        return cur;
    }

    void insert(int key, int value, rb_node **root) {}

    void print() {
        puts("digraph G { \n\
graph [ratio=.48]; node \n\
[style=filled, color=black, shape=circle, width=.6 \n\
fontname=Helvetica, fontweight=bold, \n\
fontcolor=white, fontsize=24, fixedsize=true];");
        rb_node *root = this;
        auto q = queue<rb_node *>();
        auto red_nodes = vector<rb_node *>();
        rb_node *cur;
        int nil_count = 0;
        int i;

        // 打印缓冲
        char buf[2][100];

        q.push(root);
#define fmt(x, buf)                                                            \
    {                                                                          \
        if (x) {                                                               \
            sprintf(buf, "%d", x->key);                                        \
        } else {                                                               \
            sprintf(buf, "n%d", nil_count++);                                  \
        }                                                                      \
    }

        while (q.size() > 0) {
            cur = q.front();
            q.pop();

            if (cur->color == rb_node::red()) {
                red_nodes.push_back(cur);
            }

            fmt(cur->left, buf[0]);
            fmt(cur->left, buf[1]);
            printf("%d -> %s, %s;\n", cur->key, buf[0], buf[1]);

            if (cur->left) {
                q.push(cur->left);
            }
            if (cur->right) {
                q.push(cur->right);
            }
        }

        if (nil_count > 0) {
            for (i = 0; i < nil_count; i++) {
                printf("n%d %s", i, i == nil_count - 1 ? "\n" : ",");
            }
            puts("[label=\"NIL\", shape=record, width=.4,height=.25, "
                 "fontsize=16];\n");
        }

        if (red_nodes.size() > 0) {
            for (i = 0; i < red_nodes.size(); i++) {
                printf("%d %s", red_nodes[i],
                       i == red_nodes.size() - 1 ? "\n" : ",");
            }
            puts("[fillcolor=red]\n");
        }
    }

    void fix_insert(rb_node **root) {
        rb_node *cur = this;
        rb_node *p = cur->parent;
        rb_node *uncle;
        rb_node *grandpa;

    retry:
        if (!p) {
            this->color = rb_node::black();
            return;
        }

        if (p->color == rb_node::black()) {
            return;
        }

        grandpa = p->parent;
        if (p == grandpa->left) {
            uncle = grandpa->right;
        } else {
            uncle = grandpa->left;
        }
        if (uncle && uncle->color == rb_node::red()) {
            p->color = rb_node::black();
            uncle->color = rb_node::black();
            grandpa->color = rb_node::red();
            cur = grandpa;
            goto retry;
        }

        if (grandpa == *root) {
            *root = p;
        }
        if (p == grandpa->left) {
            if (cur == p->right) {
                p->rotate_left();
                p = cur;
            }
            p->color = rb_node::black();
            grandpa->color = rb_node::red();
            if (grandpa == *root) {
                *root = p;
            }
            grandpa->rotate_right();
        } else {
            if (cur == p->left) {
                p->rotate_right();
                p = cur;
            }
            p->color = rb_node::black();
            grandpa->color = rb_node::red();
            grandpa->rotate_left();
        }
        return;
    }
};
