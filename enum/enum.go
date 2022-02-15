package enum

type Rice int

// Rice type int backed enum
// The iota keyword represents successive integer constants 0, 1, 2,…
// It resets to 0 whenever the word const appears in the source code, and increments after each const specification.
const (
	Jasmin  Rice = iota // 0
	Basmati             // 1
)

// String() will be returned when you call string(r) or printf("%s",r) etc.
// One of the most ubiquitous interfaces is Stringer defined by the fmt package.
// The fmt package (and many others) look for this interface to print values.
func (r Rice) String() string {
	// Typically, you would take the enum names, Str suffix just for illustration
	rices := [...]string{"JasminStr", "BasmatiStr", "PlainStr"} // index must follow the order in const !!!
	if len(rices) < int(r) {
		return ""
	}

	return rices[r]
}

// Soup type in backed but no String()

type Soup int

const (
	Pho Soup = iota
	TomKa
)

// Extension File Extension String backed
type Extension string

const (
	Jpeg Extension = "JPG"
	Png  Extension = "PNG"
	// Gif            = "GIF" // do not omit extension SA9004: only the first constant in this group has an explicit type
)

//func (e Extension) String() string {
//	extensions := [...]string{"JPG", "PNG", "GIF", "BMP"}
//
//	x := string(e)
//	for _, v := range extensions {
//		if v == x {
//			return x
//		}
//	}
//
//	return ""
//}
