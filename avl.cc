#include <avl.hh>

int main() {
  auto root = new avl_tree_node(0);
  root = root->insert(100);
  root = root->insert(200);
  root = root->insert(300);
  return 0;
}
