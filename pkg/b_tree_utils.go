package pkg

import "io"

/*
digraph g {
node [shape = record,height=.1];
node0[label = "<f0> |10|<f1> |20|<f2> |30|<f3>"];
node1[label = "<f0> |1|<f1> |2|<f2>"];
"node0":f0 -> "node1"
node2[label = "<f0> |11|<f1> |12|<f2>"];
"node0":f1 -> "node2"
node3[label = "<f0> |21|<f1> |22|<f2>"];
"node0":f2 -> "node3"
node4[label = "<f0> |31|<f1> |32|<f2>"];
"node0":f3 -> "node4"

}
*/
func (root *BNode) print(io.Writer) {

}
