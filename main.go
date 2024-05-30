package main

import (
	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		src, _ := file.Open()

		fs := filesystem.NewFileSystemByRelative("./abc.txt")
		_, err := fs.WriteIoReader(src)
		if err != nil {
			panic(err)
		}
	})

}
