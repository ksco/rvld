package linker

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"github.com/ksco/rvld/pkg/utils"
)

type OutputEhdr struct {
	Chunk
}

func NewOutputEhdr() *OutputEhdr {
	return &OutputEhdr{Chunk{
		Shdr: Shdr{
			Flags:     uint64(elf.SHF_ALLOC),
			Size:      uint64(EhdrSize),
			AddrAlign: 8,
		},
	}}
}

func (o *OutputEhdr) CopyBuf(ctx *Context) {
	ehdr := &Ehdr{}
	WriteMagic(ehdr.Ident[:])
	ehdr.Ident[elf.EI_CLASS] = uint8(elf.ELFCLASS64)
	ehdr.Ident[elf.EI_DATA] = uint8(elf.ELFDATA2LSB)
	ehdr.Ident[elf.EI_VERSION] = uint8(elf.EV_CURRENT)
	ehdr.Ident[elf.EI_OSABI] = 0
	ehdr.Ident[elf.EI_ABIVERSION] = 0
	ehdr.Type = uint16(elf.ET_EXEC)
	ehdr.Machine = uint16(elf.EM_RISCV)
	ehdr.Version = uint32(elf.EV_CURRENT)
	// TODO
	ehdr.EhSize = uint16(EhdrSize)
	ehdr.PhEntSize = uint16(PhdrSize)
	// TODO
	ehdr.ShEntSize = uint16(ShdrSize)
	// TODO

	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, ehdr)
	utils.MustNo(err)
	copy(ctx.Buf[o.Shdr.Offset:], buf.Bytes())
}
