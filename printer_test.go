package dapper_test

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/dogmatiq/dapper"
)

func ExamplePrinter() {
	type TreeNode struct {
		Name     string
		Value    any
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

	p := &Printer{}
	s := p.Format(v)
	fmt.Println(s)

	// output: github.com/dogmatiq/dapper_test.TreeNode{
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
	//             Value:    github.com/dogmatiq/dapper_test.NodeValue{}
	//             Children: nil
	//         }
	//     }
	// }
}

func ExamplePrinter_Config() {
	type TreeNode struct {
		Name     string
		Value    any
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

	p := &Printer{
		Config: Config{
			OmitPackagePaths: true,
		},
	}
	s := p.Format(v)
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

func ExamplePrint() {
	Print(123, 456.0)

	// output: int(123)
	// float64(456)
}

func ExampleFormat() {
	s := Format(123)
	fmt.Println(s)

	// output: int(123)
}

func ExampleWrite() {
	w := &bytes.Buffer{}

	if _, err := Write(w, 123); err != nil {
		panic(err)
	}

	w.WriteTo(os.Stdout)

	// output: int(123)
}
