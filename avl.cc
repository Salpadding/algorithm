#include <avl.hh>

int main() {
    auto root = new avl_tree_node(0);
    for (int i = 1; i < 10; i++) {
        printf("insert %d\n", i);
        root = root->insert(i);
    }

    printf("depth = %d left depth = %d right depth = %d\n", root->find_depth(), root->left_child_depth(), root->right_child_depth());
}
