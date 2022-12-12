package linker

import (
	"debug/elf"
	"github.com/ksco/rvld/pkg/utils"
)

type GotSection struct {
	Chunk
	GotTpSyms []*Symbol
}

func NewGotSection() *GotSection {
	g := &GotSection{Chunk: NewChunk()}
	g.Name = ".got"
	g.Shdr.Type = uint32(elf.SHT_PROGBITS)
	g.Shdr.Flags = uint64(elf.SHF_ALLOC | elf.SHF_WRITE)
	g.Shdr.AddrAlign = 0
	return g
}

func (g *GotSection) AddGotTpSymbol(sym *Symbol) {
	sym.GotTpIdx = int32(g.Shdr.Size / 8)
	g.Shdr.Size += 8
	g.GotTpSyms = append(g.GotTpSyms, sym)
}

type GotEntry struct {
	Idx int64
	Val uint64
}

func (g *GotSection) GetEntries(ctx *Context) []GotEntry {
	entries := make([]GotEntry, 0)
	for _, sym := range g.GotTpSyms {
		idx := sym.GotTpIdx
		entries = append(entries,
			GotEntry{Idx: int64(idx), Val: sym.GetAddr() - ctx.TpAddr})
	}

	return entries
}

func (g *GotSection) CopyBuf(ctx *Context) {
	base := ctx.Buf[g.Shdr.Offset:]
	for _, ent := range g.GetEntries(ctx) {
		utils.Write[uint64](base[ent.Idx*8:], ent.Val)
	}
}
