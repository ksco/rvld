package linker

import "debug/elf"

type ObjectFile struct {
	InputFile

	SymtabSec *Shdr
}

func NewObjectFile(file *File) *ObjectFile {
	o := &ObjectFile{InputFile: NewInputFile(file)}
	return o
}

func (o *ObjectFile) Parse() {
	o.SymtabSec = o.FindSection(uint32(elf.SHT_SYMTAB))
	if o.SymtabSec != nil {
		o.FirstGlobal = int64(o.SymtabSec.Info)
		o.FillUpElfSyms(o.SymtabSec)
		o.SymbolStrtab = o.GetBytesFromIdx(int64(o.SymtabSec.Link))
	}
}
