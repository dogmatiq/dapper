package dapper

// Annotator is a function that annotates a value with additional information.
//
// If it returns a non-empty string, the string is rendered after the value,
// regardless of whether the value is rendered by a filter.
type Annotator func(Value) string
