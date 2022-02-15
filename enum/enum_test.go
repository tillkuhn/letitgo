package enum

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntEnum(t *testing.T) {

	// The verb %v ('v' for 'value') always formats the argument in its default form, just how Print or Println would show
	// The special verb %T ('T' for 'Type') prints the type of the argument rather than its value.
	//  Implicitly Calls String() and prints "Rice: BasmatiStr (enum.Rice)"
	fmt.Printf("Rice: %v (%T)\n", Basmati, Basmati)
	fmt.Println(Jasmin)                       // also prints JasminStr
	assert.NotEqual(t, "BasmatiStr", Basmati) // string("BasmatiStr") vs enum.Rice(1)
	// isEqual1 := Basmati == "Basmati" // compile error mismatched types Rice and untyped string
	rice := Basmati
	assert.Equal(t, Basmati, rice)
	assert.NotEqual(t, 1, Basmati)  // not equal enum.Rice(1) and int(1)
	assert.False(t, rice == Jasmin) // sure it's false, but no compile error if you compare
	p := fmt.Sprintf("%v", TomKa)   // prints 1 as there is no String()
	assert.Equal(t, "1", p)
	assert.Equal(t, 0, int(Pho)) // works int == int

}

func TestStringEnum(t *testing.T) {
	assert.Equal(t, "JPG", fmt.Sprintf("%v", Jpeg))
	assert.True(t, Jpeg == "JPG") // yes you can do that
	assert.False(t, Jpeg == "NIX")
	strVal := string(Png)
	assert.Equal(t, "PNG", strVal)
}
