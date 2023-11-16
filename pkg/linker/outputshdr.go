package linker

import "github.com/ksco/rvld/pkg/utils"

type OutputShdr struct {
	Chunk
}

func NewOutputShdr() *OutputShdr {
	o := &OutputShdr{Chunk: NewChunk()}
	o.Shdr.AddrAlign = 8
	return o
}

func (o *OutputShdr) UpdateShdr(ctx *Context) {
	o.Shdr.Size = 1 * uint64(ShdrSize)
}

func (o *OutputShdr) CopyBuf(ctx *Context) {
	base := ctx.Buf[o.Shdr.Offset:]
	utils.Write[Shdr](base, Shdr{})
}
