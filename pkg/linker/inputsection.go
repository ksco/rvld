package linker

import "github.com/ksco/rvld/pkg/utils"

type InputSection struct {
	File     *ObjectFile
	Contents []byte
	Shndx    uint32
}

func NewInputSection(file *ObjectFile, shndx uint32) *InputSection {
	s := &InputSection{File: file, Shndx: shndx}

	shdr := s.Shdr()
	s.Contents = file.File.Contents[shdr.Offset : shdr.Offset+shdr.Size]

	return s
}

func (i *InputSection) Shdr() *Shdr {
	utils.Assert(i.Shndx < uint32(len(i.File.ElfSections)))
	return &i.File.ElfSections[i.Shndx]
}

func (i *InputSection) Name() string {
	return ElfGetName(i.File.ShStrtab, i.Shdr().Name)
}
