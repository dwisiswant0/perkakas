package graphql

import (
	graphqlgo "github.com/graphql-go/graphql"
)

// NullableStringField field for nullable string
var NullableStringField = &graphqlgo.Field{
	Type:    graphqlgo.String,
	Resolve: ResolveNullable,
}

// NullableInt64Field field for nullable int64
var NullableInt64Field = &graphqlgo.Field{
	Type:    graphqlgo.Int,
	Resolve: ResolveNullable,
}

// NullableFloat64Field field for nullable float64
var NullableFloat64Field = &graphqlgo.Field{
	Type:    graphqlgo.Float,
	Resolve: ResolveNullable,
}

// NullableBoolField field for nullable bool
var NullableBoolField = &graphqlgo.Field{
	Type:    graphqlgo.Boolean,
	Resolve: ResolveNullable,
}

// NullableTimeField field for nullable datetime
var NullableDateTimeField = &graphqlgo.Field{
	Type:    graphqlgo.DateTime,
	Resolve: ResolveNullable,
}
