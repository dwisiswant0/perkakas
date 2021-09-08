package graphql

import "sync"

var (
	gqlInst *Gql
	once    = new(sync.Once)
)

func init() {
	once.Do(func() {
		gqlInst = new(Gql)
	})
}

// Gql defines the variable needed to initialize this package.
type Gql struct {
	IsSnakeCase bool
}

// Options define the functional options for this package.
type Options func(*Gql)

// Init initialize the gql package.
func Init(opts ...Options) {
	for _, opt := range opts {
		opt(gqlInst)
	}
}

// WithSnakeCase set if this package use snake case convention.
func WithSnakeCase(isSnakeCase bool) Options {
	return func(g *Gql) {
		g.IsSnakeCase = isSnakeCase
	}
}
