package filesystem

import (
	"fmt"
	"path"
)

func DoStuff() {
	fmt.Println(join("hase", "horst", "hubert"))
}

func join(pathParts ...string) string {
	return path.Join(pathParts...) // should use os.PathSeparator
}
