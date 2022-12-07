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

func getEntryAddr(ctx *Context) uint64 {
	for _, osec := range ctx.OutputSections {
		if osec.Name == ".text" {
			return osec.Shdr.Addr
		}
	}
	return 0
}

func getFlags(ctx *Context) uint32 {
	utils.Assert(len(ctx.Objs) > 0)
	flags := ctx.Objs[0].GetEhdr().Flags
	for _, obj := range ctx.Objs[1:] {
		if obj == ctx.InternalObj {
			continue
		}

		if obj.GetEhdr().Flags&EF_RISCV_RVC != 0 {
			flags |= EF_RISCV_RVC
			break
		}
	}

	return flags
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
	ehdr.Entry = getEntryAddr(ctx)
	ehdr.PhOff = ctx.Phdr.Shdr.Offset
	ehdr.ShOff = ctx.Shdr.Shdr.Offset
	ehdr.Flags = getFlags(ctx)
	ehdr.EhSize = uint16(EhdrSize)
	ehdr.PhEntSize = uint16(PhdrSize)
	ehdr.PhNum = uint16(ctx.Phdr.Shdr.Size) / uint16(PhdrSize)
	ehdr.ShEntSize = uint16(ShdrSize)
	ehdr.ShNum = uint16(ctx.Shdr.Shdr.Size) / uint16(ShdrSize)

	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, ehdr)
	utils.MustNo(err)
	copy(ctx.Buf[o.Shdr.Offset:], buf.Bytes())
}
