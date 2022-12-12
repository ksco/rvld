package linker

import (
	"bytes"
	"debug/elf"
	"github.com/ksco/rvld/pkg/utils"
	"math"
)

type ObjectFile struct {
	InputFile
	SymtabSec         *Shdr
	SymtabShndxSec    []uint32
	Sections          []*InputSection
	MergeableSections []*MergeableSection
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

	o.InitializeSections(ctx)
	o.InitializeSymbols(ctx)
	o.InitializeMergeableSections(ctx)
	o.SkipEhframeSections()
}

func (o *ObjectFile) InitializeSections(ctx *Context) {
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
			name := ElfGetName(o.InputFile.ShStrtab, shdr.Name)
			o.Sections[i] = NewInputSection(ctx, name, o, uint32(i))
		}
	}

	for i := 0; i < len(o.ElfSections); i++ {
		shdr := &o.InputFile.ElfSections[i]
		if shdr.Type != uint32(elf.SHT_RELA) {
			continue
		}

		utils.Assert(shdr.Info < uint32(len(o.Sections)))
		if target := o.Sections[shdr.Info]; target != nil {
			utils.Assert(target.RelsecIdx == math.MaxUint32)
			target.RelsecIdx = uint32(i)
		}
	}
}

func (o *ObjectFile) FillUpSymtabShndxSec(s *Shdr) {
	bs := o.GetBytesFromShdr(s)
	o.SymtabShndxSec = utils.ReadSlice[uint32](bs, 4)
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

func (o *ObjectFile) MarkLiveObjects(feeder func(*ObjectFile)) {
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

func (o *ObjectFile) InitializeMergeableSections(ctx *Context) {
	o.MergeableSections = make([]*MergeableSection, len(o.Sections))
	for i := 0; i < len(o.Sections); i++ {
		isec := o.Sections[i]
		if isec != nil && isec.IsAlive &&
			isec.Shdr().Flags&uint64(elf.SHF_MERGE) != 0 {
			o.MergeableSections[i] = splitSection(ctx, isec)
			isec.IsAlive = false
		}
	}
}

func findNull(data []byte, entSize int) int {
	if entSize == 1 {
		return bytes.Index(data, []byte{0})
	}

	for i := 0; i <= len(data)-entSize; i += entSize {
		bs := data[i : i+entSize]
		if utils.AllZeros(bs) {
			return i
		}
	}

	return -1
}

func splitSection(ctx *Context, isec *InputSection) *MergeableSection {
	m := &MergeableSection{}
	shdr := isec.Shdr()

	m.Parent = GetMergedSectionInstance(ctx, isec.Name(), shdr.Type,
		shdr.Flags)
	m.P2Align = isec.P2Align

	data := isec.Contents
	offset := uint64(0)
	if shdr.Flags&uint64(elf.SHF_STRINGS) != 0 {
		for len(data) > 0 {
			end := findNull(data, int(shdr.EntSize))
			if end == -1 {
				utils.Fatal("string is not null terminated")
			}

			sz := uint64(end) + shdr.EntSize
			substr := data[:sz]
			data = data[sz:]
			m.Strs = append(m.Strs, string(substr))
			m.FragOffsets = append(m.FragOffsets, uint32(offset))
			offset += sz
		}
	} else {
		if uint64(len(data))%shdr.EntSize != 0 {
			utils.Fatal("section size is not multiple of entsize")
		}

		for len(data) > 0 {
			substr := data[:shdr.EntSize]
			data = data[shdr.EntSize:]
			m.Strs = append(m.Strs, string(substr))
			m.FragOffsets = append(m.FragOffsets, uint32(offset))
			offset += shdr.EntSize
		}
	}

	return m
}

func (o *ObjectFile) RegisterSectionPieces() {
	for _, m := range o.MergeableSections {
		if m == nil {
			continue
		}

		m.Fragments = make([]*SectionFragment, 0, len(m.Strs))
		for i := 0; i < len(m.Strs); i++ {
			m.Fragments = append(m.Fragments,
				m.Parent.Insert(m.Strs[i], uint32(m.P2Align)))
		}
	}

	for i := 1; i < len(o.ElfSyms); i++ {
		sym := o.Symbols[i]
		esym := &o.ElfSyms[i]

		if esym.IsAbs() || esym.IsUndef() || esym.IsCommon() {
			continue
		}

		m := o.MergeableSections[o.GetShndx(esym, i)]
		if m == nil {
			continue
		}

		frag, fragOffset := m.GetFragment(uint32(esym.Val))
		if frag == nil {
			utils.Fatal("bad symbol value")
		}
		sym.SetSectionFragment(frag)
		sym.Value = uint64(fragOffset)
	}
}

func (o *ObjectFile) SkipEhframeSections() {
	for _, isec := range o.Sections {
		if isec != nil && isec.IsAlive && isec.Name() == ".eh_frame" {
			isec.IsAlive = false
		}
	}
}

func (o *ObjectFile) ScanRelocations() {
	for _, isec := range o.Sections {
		if isec != nil && isec.IsAlive &&
			isec.Shdr().Flags&uint64(elf.SHF_ALLOC) != 0 {
			isec.ScanRelocations()
		}
	}
}
