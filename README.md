<div align="center">

# Dapper

Dapper is a pretty-printer for Go values.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/dapper)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dapper.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/dapper/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/dapper/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/dapper/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dapper/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/dapper)

</div>

Dapper is not intended to be used directly as a debugging tool, but as a library
for applications that need to describe Go values to humans, such as testing
frameworks.

Some features and design goals include:

- Concise formatting, without type ambiguity
- Deterministic output, useful for generating diffs using standard tools
- A filter system for producing customized output on a per-value basis

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
	Value    any
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
