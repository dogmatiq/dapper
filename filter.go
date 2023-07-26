package dapper

// Filter is a function that provides custom formatting logic for a specific
// type or value.
//
// The filter uses r to render v. If r is unused v is rendered using the default
// formatting logic.
type Filter func(r Renderer, v Value)
