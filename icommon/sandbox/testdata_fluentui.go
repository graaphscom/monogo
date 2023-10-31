package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func main() {
	wd, _ := os.Getwd()
	wd = path.Join(wd, "testdata", "initial", "fluentui", "assets", "Next Modified", "SVG")
	dirEntries, _ := os.ReadDir(wd)
	replacer := strings.NewReplacer("accessibility", "next_modified")

	for _, dirEntry := range dirEntries {
		fmt.Println(os.Rename(path.Join(wd, dirEntry.Name()), path.Join(wd, replacer.Replace(dirEntry.Name()))))
	}
}
