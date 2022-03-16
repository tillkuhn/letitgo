package types

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// SOURCE: https://eli.thegreenplace.net/2020/embedding-in-go-part-1-structs-in-structs/
//
// Go doesn't support inheritance in the classical sense; instead, in encourages composition as a way to extend the
// functionality of types. This is not a notion peculiar to Go. Composition over inheritance is a known principle of
// OOP and is featured in the very first chapter of the Design Patterns book.
//
// Embedding is an important Go feature making composition more convenient and useful. While Go strives to be simple,
// embedding is one place where the essential complexity of the problem leaks somewhat. In this series of short posts,

// Params is the embedded struct
type Params struct {
	Name string
	ID   int
}

func (p Params) Describe() string {
	return fmt.Sprintf("Name: %s", p.Name)
}

// Report is the embedding struct
type Report struct {
	Params
	Payload string
	ID      int // Shadows embedded Params ID
}

func TestEmbedded(t *testing.T) {

	rep := Report{
		Params: Params{
			Name: "Hase",
			ID:   4711,
		},
		Payload: "Cost is 123 €",
		ID:      7,
	}
	rep.Name = "Hase2"
	assert.Equal(t, "Hase2", rep.Name) // can access Name Directly on embedding struct
	assert.Contains(t, rep.Payload, "123")
	assert.Equal(t, 7, rep.ID)                     // displays the "outer" id
	assert.Equal(t, 4711, rep.Params.ID)           // ID is ambiguous, it wil use the ID of embedding struct
	assert.Equal(t, "Name: Hase2", rep.Describe()) // Describe method from params
}
