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
	objFile := linker.NewObjectFile(file)
	objFile.Parse()
	utils.Assert(len(objFile.ElfSections) == 11)
	utils.Assert(objFile.FirstGlobal == 10)
	utils.Assert(len(objFile.ElfSyms) == 12)

	for _, sym := range objFile.ElfSyms {
		println(linker.ElfGetName(objFile.SymbolStrtab, sym.Name))
	}
}
