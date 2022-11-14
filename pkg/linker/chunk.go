package linker

type Chunk struct {
	Name string
	Shdr Shdr
}

func NewChunk() Chunk {
	return Chunk{Shdr: Shdr{AddrAlign: 1}}
}
