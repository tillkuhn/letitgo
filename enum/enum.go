package enum

type Rice int

// Rice type int backed enum
const (
	Jasmin  Rice = iota // 0
	Basmati             // 1
)

func (r Rice) String() string {
	// Typically, you would take the enum names, Str suffix just for illustration
	rices := [...]string{"JasminStr", "BasmatiStr", "PlainStr"}
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
