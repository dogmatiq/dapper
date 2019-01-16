package dapper_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dogmatiq/dapper"
)

func ExampleFormat() {
	type TreeNode struct {
		Name     string
		Value    interface{}
		Children []*TreeNode
	}

	type NodeValue struct{}

	v := TreeNode{
		Name: "root",
		Children: []*TreeNode{
			{
				Name:  "branch #1",
				Value: 100,
			},
			{
				Name:  "branch #2",
				Value: NodeValue{},
			},
		},
	}

	s := Format(v)
	fmt.Println(s)

	// output: dapper_test.TreeNode{
	//     Name:     "root"
	//     Value:    nil
	//     Children: {
	//         {
	//             Name:     "branch #1"
	//             Value:    int(100)
	//             Children: nil
	//         }
	//         {
	//             Name:     "branch #2"
	//             Value:    dapper_test.NodeValue{}
	//             Children: nil
	//         }
	//     }
	// }
}

func test(
	t *testing.T,
	n string,
	v interface{},
	lines ...string,
) {
	x := strings.Join(lines, "\n")

	t.Run(
		n,
		func(t *testing.T) {
			p := Format(v)

			t.Log("expected:\n\n" + x + "\n")

			if p != x {
				t.Fatal("actual:\n\n" + p + "\n")
			}
		},
	)
}
