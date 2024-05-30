package main

import (
	"github.com/jericho-yu/filesystem/filesystem"
)

func main() {
	var (
		e   error
		src *filesystem.FileSystem
	)

	src = filesystem.NewFileSystemByRelative("a")
	e = src.CopyDir("b", false)

	if e != nil {
		panic(e)
	}

}
