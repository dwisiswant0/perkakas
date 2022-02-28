# Sanitize

This is [bluemonday](https://github.com/microcosm-cc/bluemonday)-wrapper to perform a deep-clean and/or sanitization of interface. It'll sanitize all fields in the structure, supports fields with (slice of, and/or just) `string` and `byte` types.

## Examples

```golang
import "github.com/kitabisa/perkakas/v2/sanitize"

t := &T{/* ... */} // Should be a pointer
s := true // Use strict policy? See [StrictPolicy](https://pkg.go.dev/github.com/microcosm-cc/bluemonday#StrictPolicy)

t = sanitize.Clean(t, s).(*T) // Converting interface to its type
```

### Workaround

The following is an example of how this module is implemented & works:

```golang
type Info struct {
	Fname string
	Lname string
	Phone int64
	Notes []string
	Story []byte
}

func main() {
	i := &Info{
		Fname: "Foo",
		Lname: "<b>Bar</b>",
		Phone: 911,
		Notes: []string{
			"Hello,",
			`world<script>alert("world")</script>!`,
		},
		Story: []byte(`<blockquote><h1>Lorem</h1> <p>ipsum.</p></blockquote>`),
	}
	i = sanitize.Clean(i, true).(*Info)

	fmt.Printf("%+v\n", i)
	// Output: &{Fname:Foo Lname:Bar Phone:911 Notes:[Hello, world!] Story:[76 111 114 101 109 32 105 112 115 117 109 46]}
	// which i.Strory produced: Lorem ipsum (as string).
}
```

## Limitations

Nested types are **NOT** supported, so make sure the type inside is sanitized first.