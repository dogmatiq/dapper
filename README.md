# Dapper

[![Build Status](http://img.shields.io/travis/com/dogmatiq/dapper/master.svg)](https://travis-ci.com/dogmatiq/dapper)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dapper/master.svg)](https://codecov.io/github/dogmatiq/dapper)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dapper.svg?label=semver)](https://semver.org)
[![GoDoc](https://godoc.org/github.com/dogmatiq/dapper?status.svg)](https://godoc.org/github.com/dogmatiq/dapper)
[![Go Report Card](https://goreportcard.com/badge/github.com/dogmatiq/dapper)](https://goreportcard.com/report/github.com/dogmatiq/dapper)

Dapper is a pretty-printer for Go values that aims to produce the shortest
possible output without ambiguity.

Output is designed to be deterministic, to allow for the generation of useful
diffs using standard tools.

## Example

 This example renders a basic tree structure. Note that the output only includes
 type names where the value's type can not be inferred from the context.

### Code
```go
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
```

### Output

```
dapper_test.TreeNode{
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
            Value:    dapper_test.SpecialValue{}
            Children: nil
        }
    }
}
```
