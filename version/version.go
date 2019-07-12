package version

import (
	"fmt"
	"os"
	"path"
)

var (
	gitCommit string
)

func init() {
	fmt.Printf("%s gitCommit: %s \n", path.Base(os.Args[0]), gitCommit)
}
