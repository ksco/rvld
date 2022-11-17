package linker

type Chunker interface {
	GetShdr() *Shdr
	CopyBuf(ctx *Context)
}

type Chunk struct {
	Name string
	Shdr Shdr
}

func NewChunk() Chunk {
	return Chunk{Shdr: Shdr{AddrAlign: 1}}
}

func (c *Chunk) GetShdr() *Shdr {
	return &c.Shdr
}

func (c *Chunk) CopyBuf(ctx *Context) {}
