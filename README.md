# Dapper

[![Build Status](https://github.com/dogmatiq/dapper/workflows/CI/badge.svg)](https://github.com/dogmatiq/dapper/actions?workflow=CI)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dapper/master.svg)](https://codecov.io/github/dogmatiq/dapper)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dapper.svg?label=semver)](https://semver.org)
[![GoDoc](https://godoc.org/github.com/dogmatiq/dapper?status.svg)](https://godoc.org/github.com/dogmatiq/dapper)
[![Go Report Card](https://goreportcard.com/badge/github.com/dogmatiq/dapper)](https://goreportcard.com/report/github.com/dogmatiq/dapper)

Dapper is a pretty-printer for Go values.

It is not intended to be used directly as a debugging tool, but as a library
for applications that need to describe Go values to humans, such as testing
frameworks.

Some features include:

- Concise formatting, without type ambiguity
- Deterministic output, useful for generating diffs using standard tools
- A filtering system for producing customized output on a per-value basis

## Example

This example renders a basic tree structure. Note that the output only includes
type names where the value's type can not be inferred from the context.

### Code

```go
package main

import (
	"fmt"
	"github.com/dogmatiq/dapper"
)

type TreeNode struct {
	Name     string
	Value    interface{}
	Children []*TreeNode
}

type NodeValue struct{}

func main() {
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

	dapper.Print(v)
}
```

### Output

```
main.TreeNode{
    Name:     "root"
    Value:    nil
    Children: {
        {
            Name:     "branch #1"
            Value:    int(100)
            Children: nil
        }
        {
            Name:     "branch #2"
            Value:    main.NodeValue{}
            Children: nil
        }
    }
}
```
