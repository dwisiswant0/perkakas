# Elastic Search
This package provides helper for common functionality using elastic search, such as
search, index (upsert), remove document, remove index and bulk index. You expected to
embed this elastic search client to your client and wrap the elastic search common
functionality with your specific needs.

## Embedding Client
You can just embed the client and create wrapper function like this:
```go
package myesclient

import ( 
    "github.com/kitabisa/perkakas/elastic"
    es "github.com/olivere/elastic/v7"
    "context"
    "fmt"
)

type MyClient struct {
    *elastic.ElasticClient
}

func NewMyClient(url string) *MyClient{
	c, err := elastic.NewClient(url)
	if err != nil {
		fmt.Println(err)
	}

	return &MyClient {
		ElasticClient: c,
	}
}

func (m *MyClient) createExampleTemplate() *elastic.DynamicTemplate {
    // campaigner dynamic template here
	mappings := elastic.NewMappings().
		AddDynamicTemplate("id_field", MatchConditions{
			Match:            "id",
			MatchMappingType: "string",
			Mapping: MatchMapping{
				Type: "text",
			},
		}).
		AddDynamicTemplate("name_field", MatchConditions{
			Match:            "name",
			MatchMappingType: "string",
			Mapping: MatchMapping{
				Type: "text",
				Fields: map[string]Field{
					"keyword": Field{
						Type:        "keyword",
						IgnoreAbove: 256,
					},
				},
			},
		}).
		AddDynamicTemplate("age", MatchConditions{
			Match:            "age",
			MatchMappingType: "double",
			Mapping: MatchMapping{
				Type: "double",
			},
		})

	return &elastic.DynamicTemplate{
		Settings: map[string]interface{}{
			"index.refresh_interval": "5s",
		},
		Mappings: mappings,
	}
}

// This will shadow Store function from perkakas
func (m *MyClient) Store(ctx context.Context, name string, doc interface{}) (res *es.IndexResponse, err error) {
    template := m.createExampleTemplate()
    return m.Client.Store(ctx context.Context, name string, doc interface{}, template)
}

// other wrap function here...
```

## Connecting To Elastic Search
```go
client, err := NewClient("http://localhost:9200")
if err != nil {
    fmt.Println(err)
}
```

## Search
```go
ctx := context.Background()
searchOption := SearchOption{
    Query: es.NewWildcardQuery("id", "*1234*"),
    Sort:  nil,
    Size:  10, // if you not set the size, it will be default to 10
}

result, err := client.Search(ctx, "campaign-test-111", searchOption)
if err != nil {
    // handle error
}
```

Search result is `[]byte` with json format
The query can be any query model for elastic search. For reference see here: https://github.com/olivere/elastic/wiki/Search

## Index Document (Upsert)
```go
person := Person{
    ID: "person-123",
    Name: "Mukidi",
    Age: 45,
}

ctx := context.Background()
mapping := elasticCreateIndex()

result, err := client.Store(ctx, "person-test-123", person, mapping)
if err != nil {
    // handle error here...
}
```
Index id is string and can have any value, but better if you give prefix for avoiding collision in the id

## Remove Document
```go
ctx := context.Background()
_, err := client.Remove(ctx, "person-test-123", "1234")
if err != nil {
    // handle error here...
}
```

## Remove Index
```go
ctx := context.Background()
_, err := client.RemoveIndex(ctx, "index1", "index2")
if err != nil {
    t.Log(err)
    t.FailNow()
}
```

## Bulk Index Document (Upsert)
```go
ctx := context.Background()
docs := []interface{}{
    Person {
        ID:   "abcde",
        Name: "Budi",
        Age: 17,
    },
    Person {
        ID:   "fghij",
        Name: "Badu",
        Age: 50,
    },
}

template := elasticCreateIndex()
err := client.BulkStore(ctx, "person-list", "default", docs, template)
if err != nil {
    // handle error here
}
```
As you can see when client.BulkStore() called, it require a bulk processor. You can ignore it by give an empty string
argument, and it will use default bulk processor.
You can also create your own bulk processor and add it to the client:
```go
bulkProcessor := BulkProcessor{
    Name:          "my-bulk-processor", // name of the bulk processor
    Workers:       10, // how many workers
    BulkActions:   100, // flush when reach 100 bulk requests
    BulkSize:      2 << 20, // flush when reach bulk request size 2 MB
    FlushInterval: 1 * time.Second, // flush every 1 second
}

err := client.AddBulkProcessor(bulkProcessor)
```

Or if you have add bulk processor earlier, you can get the bulk processor with `client.GetBulkProcessor("name")`.
It will returns the bulk processor and error if cannot find the specified bulk processor name.

