package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	appName := path.Base(os.Args[0])
	fmt.Printf("%s not implemented yet.\n", appName)
}
