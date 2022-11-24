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
	n := uint64(0)
	for _, chunk := range ctx.Chunks {
		if chunk.GetShndx() > 0 {
			n = uint64(chunk.GetShndx())
		}
	}

	o.Shdr.Size = (n + 1) * uint64(ShdrSize)
}

func (o *OutputShdr) CopyBuf(ctx *Context) {
	base := ctx.Buf[o.Shdr.Offset:]
	utils.Write[Shdr](base, Shdr{})

	for _, chunk := range ctx.Chunks {
		if chunk.GetShndx() > 0 {
			utils.Write[Shdr](base[chunk.GetShndx()*int64(ShdrSize):],
				*chunk.GetShdr())
		}
	}
}
