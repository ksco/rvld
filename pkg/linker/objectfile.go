package linker

import (
	"debug/elf"
	"github.com/ksco/rvld/pkg/utils"
)

type ObjectFile struct {
	InputFile
	SymtabSec      *Shdr
	SymtabShndxSec []uint32
	Sections       []*InputSection
}

func NewObjectFile(file *File, isAlive bool) *ObjectFile {
	o := &ObjectFile{InputFile: NewInputFile(file)}
	o.IsAlive = isAlive
	return o
}

func (o *ObjectFile) Parse(ctx *Context) {
	o.SymtabSec = o.FindSection(uint32(elf.SHT_SYMTAB))
	if o.SymtabSec != nil {
		o.FirstGlobal = int(o.SymtabSec.Info)
		o.FillUpElfSyms(o.SymtabSec)
		o.SymbolStrtab = o.GetBytesFromIdx(int64(o.SymtabSec.Link))
	}

	o.InitializeSections()
	o.InitializeSymbols(ctx)
}

func (o *ObjectFile) InitializeSections() {
	o.Sections = make([]*InputSection, len(o.ElfSections))
	for i := 0; i < len(o.ElfSections); i++ {
		shdr := &o.ElfSections[i]
		switch elf.SectionType(shdr.Type) {
		case elf.SHT_GROUP, elf.SHT_SYMTAB, elf.SHT_STRTAB, elf.SHT_REL, elf.SHT_RELA,
			elf.SHT_NULL:
			break
		case elf.SHT_SYMTAB_SHNDX:
			o.FillUpSymtabShndxSec(shdr)
		default:
			o.Sections[i] = NewInputSection(o, uint32(i))
		}
	}
}

func (o *ObjectFile) FillUpSymtabShndxSec(s *Shdr) {
	bs := o.GetBytesFromShdr(s)
	nums := len(bs) / 4
	for nums > 0 {
		o.SymtabShndxSec = append(o.SymtabShndxSec, utils.Read[uint32](bs))
		bs = bs[4:]
		nums--
	}
}

func (o *ObjectFile) InitializeSymbols(ctx *Context) {
	if o.SymtabSec == nil {
		return
	}

	o.LocalSymbols = make([]Symbol, o.FirstGlobal)
	for i := 0; i < len(o.LocalSymbols); i++ {
		o.LocalSymbols[i] = *NewSymbol("")
	}
	o.LocalSymbols[0].File = o

	for i := 1; i < len(o.LocalSymbols); i++ {
		esym := &o.ElfSyms[i]
		sym := &o.LocalSymbols[i]
		sym.Name = ElfGetName(o.SymbolStrtab, esym.Name)
		sym.File = o
		sym.Value = esym.Val
		sym.SymIdx = i

		if !esym.IsAbs() {
			sym.SetInputSection(o.Sections[o.GetShndx(esym, i)])
		}
	}

	o.Symbols = make([]*Symbol, len(o.ElfSyms))
	for i := 0; i < len(o.LocalSymbols); i++ {
		o.Symbols[i] = &o.LocalSymbols[i]
	}

	for i := len(o.LocalSymbols); i < len(o.ElfSyms); i++ {
		esym := &o.ElfSyms[i]
		name := ElfGetName(o.SymbolStrtab, esym.Name)
		o.Symbols[i] = GetSymbolByName(ctx, name)
	}
}

func (o *ObjectFile) GetShndx(esym *Sym, idx int) int64 {
	utils.Assert(idx >= 0 && idx < len(o.ElfSyms))

	if esym.Shndx == uint16(elf.SHN_XINDEX) {
		return int64(o.SymtabShndxSec[idx])
	}
	return int64(esym.Shndx)
}

func (o *ObjectFile) ResolveSymbols() {
	for i := o.FirstGlobal; i < len(o.ElfSyms); i++ {
		sym := o.Symbols[i]
		esym := &o.ElfSyms[i]

		if esym.IsUndef() {
			continue
		}

		var isec *InputSection
		if !esym.IsAbs() {
			isec = o.GetSection(esym, i)
			if isec == nil {
				continue
			}
		}

		if sym.File == nil {
			sym.File = o
			sym.SetInputSection(isec)
			sym.Value = esym.Val
			sym.SymIdx = i
		}
	}
}

func (o *ObjectFile) GetSection(esym *Sym, idx int) *InputSection {
	return o.Sections[o.GetShndx(esym, idx)]
}

func (o *ObjectFile) MarkLiveObjects(ctx *Context, feeder func(*ObjectFile)) {
	utils.Assert(o.IsAlive)

	for i := o.FirstGlobal; i < len(o.ElfSyms); i++ {
		sym := o.Symbols[i]
		esym := &o.ElfSyms[i]

		if sym.File == nil {
			continue
		}

		if esym.IsUndef() && !sym.File.IsAlive {
			sym.File.IsAlive = true
			feeder(sym.File)
		}
	}
}

func (o *ObjectFile) ClearSymbols() {
	for _, sym := range o.Symbols[o.FirstGlobal:] {
		if sym.File == o {
			sym.Clear()
		}
	}
}
