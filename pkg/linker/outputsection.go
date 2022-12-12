package linker

import "debug/elf"

type OutputSection struct {
	Chunk
	Members []*InputSection
	Idx     uint32
}

func NewOutputSection(
	name string, typ uint32, flags uint64, idx uint32) *OutputSection {
	o := &OutputSection{Chunk: NewChunk()}
	o.Name = name
	o.Shdr.Type = typ
	o.Shdr.Flags = flags
	o.Idx = idx
	return o
}

func (o *OutputSection) CopyBuf(ctx *Context) {
	if o.Shdr.Type == uint32(elf.SHT_NOBITS) {
		return
	}

	base := ctx.Buf[o.Shdr.Offset:]
	for _, isec := range o.Members {
		isec.WriteTo(ctx, base[isec.Offset:])
	}
}

func GetOutputSection(
	ctx *Context, name string, typ, flags uint64) *OutputSection {
	name = GetOutputName(name, flags)
	flags = flags &^ uint64(elf.SHF_GROUP) &^
		uint64(elf.SHF_COMPRESSED) &^ uint64(elf.SHF_LINK_ORDER)

	find := func() *OutputSection {
		for _, osec := range ctx.OutputSections {
			if name == osec.Name && typ == uint64(osec.Shdr.Type) &&
				flags == osec.Shdr.Flags {
				return osec
			}
		}
		return nil
	}

	if osec := find(); osec != nil {
		return osec
	}

	osec := NewOutputSection(name, uint32(typ), flags,
		uint32(len(ctx.OutputSections)))
	ctx.OutputSections = append(ctx.OutputSections, osec)
	return osec
}
