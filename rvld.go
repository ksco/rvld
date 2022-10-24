package main

import (
	"github.com/ksco/rvld/pkg/linker"
	"github.com/ksco/rvld/pkg/utils"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		utils.Fatal("wrong args")
	}

	file := linker.MustNewFile(os.Args[1])

	inputFile := linker.NewInputFile(file)
	utils.Assert(len(inputFile.ElfSections) == 11)
}
