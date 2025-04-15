package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Sets in Go. Go does not have a type for sets.
// A better alternative is to use map[T]struct{} (a map with empty structs as values).
// https://emersion.fr/blog/2017/sets-in-go/

func TestSet(t *testing.T) {
	// Create a set
	set := make(map[string]struct{})

	// Add some values
	set["till"] = struct{}{}
	set["spock"] = struct{}{}

	// Check if some values are in the set
	if _, ok := set["till"]; ok {
		fmt.Println("Till is in the set")
	} else {
		assert.Fail(t, "we shouldn't be here")
	}
	if _, ok := set["hase"]; !ok {
		fmt.Println("hase is not in the set")
	} else {
		assert.Fail(t, "we shouldn't be here")
	}

	// Remove a value
	delete(set, "till")

	// List values
	for key := range set {
		assert.Equal(t, "spock", key)
	}
}
