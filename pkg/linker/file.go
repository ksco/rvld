package linker

import (
	"github.com/ksco/rvld/pkg/utils"
	"os"
)

type File struct {
	Name     string
	Contents []byte
}

func MustNewFile(filename string) *File {
	contents, err := os.ReadFile(filename)
	utils.MustNo(err)
	return &File{
		Name:     filename,
		Contents: contents,
	}
}
