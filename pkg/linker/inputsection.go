package linker

import (
	"debug/elf"
	"github.com/ksco/rvld/pkg/utils"
	"math/bits"
)

type InputSection struct {
	File     *ObjectFile
	Contents []byte
	Shndx    uint32
	ShSize   uint32
	IsAlive  bool
	P2Align  uint8

	Offset        uint32
	OutputSection *OutputSection
}

func NewInputSection(ctx *Context, name string, file *ObjectFile, shndx uint32) *InputSection {
	s := &InputSection{
		File:    file,
		Shndx:   shndx,
		IsAlive: true,
	}

	shdr := s.Shdr()
	s.Contents = file.File.Contents[shdr.Offset : shdr.Offset+shdr.Size]

	utils.Assert(shdr.Flags&uint64(elf.SHF_COMPRESSED) == 0)
	s.ShSize = uint32(shdr.Size)

	toP2Align := func(align uint64) uint8 {
		if align == 0 {
			return 0
		}
		return uint8(bits.TrailingZeros64(align))
	}
	s.P2Align = toP2Align(shdr.AddrAlign)

	s.OutputSection = GetOutputSection(
		ctx, name, uint64(shdr.Type), shdr.Flags)

	return s
}

func (i *InputSection) Shdr() *Shdr {
	utils.Assert(i.Shndx < uint32(len(i.File.ElfSections)))
	return &i.File.ElfSections[i.Shndx]
}

func (i *InputSection) Name() string {
	return ElfGetName(i.File.ShStrtab, i.Shdr().Name)
}

func (i *InputSection) WriteTo(buf []byte) {
	if i.Shdr().Type == uint32(elf.SHT_NOBITS) || i.ShSize == 0 {
		return
	}

	i.CopyContents(buf)
}

func (i *InputSection) CopyContents(buf []byte) {
	copy(buf, i.Contents)
}
