package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ExampleData struct {
	Field1 string
	Field2 []string
	Field3 []byte
	Field4 int64
	Field5 interface{}
}

func GetExampleData() *ExampleData {
	return &ExampleData{
		Field1: "Foo <b>Bar</b>",
		Field2: []string{
			"Hello,",
			`world<script>alert("world")</script>!`,
		},
		Field3: []byte("<blockquote><h1>Lorem</h1> <p>ipsum.</p></blockquote>"),
		Field4: 911,
		Field5: nil,
	}
}

func TestClean(t *testing.T) {
	e := &ExampleData{
		Field1: "Foo <b>Bar</b>",
		Field2: []string{"Hello,", "world!"},
		Field3: []byte("<blockquote><h1>Lorem</h1> <p>ipsum.</p></blockquote>"),
		Field4: 911,
		Field5: nil,
	}

	assert.Exactly(t, e, Clean(GetExampleData(), false))
}

func TestCleanWithStrict(t *testing.T) {
	e := &ExampleData{
		Field1: "Foo Bar",
		Field2: []string{"Hello,", "world!"},
		Field3: []byte("Lorem ipsum."),
		Field4: 911,
		Field5: nil,
	}

	assert.Exactly(t, e, Clean(GetExampleData(), true))
}
