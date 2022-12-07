package linker

import "github.com/ksco/rvld/pkg/utils"

type Symbol struct {
	File   *ObjectFile
	Name   string
	Value  uint64
	SymIdx int

	InputSection    *InputSection
	SectionFragment *SectionFragment
}

func NewSymbol(name string) *Symbol {
	s := &Symbol{
		Name:   name,
		SymIdx: -1,
	}
	return s
}

func (s *Symbol) SetInputSection(isec *InputSection) {
	s.InputSection = isec
	s.SectionFragment = nil
}

func (s *Symbol) SetSectionFragment(frag *SectionFragment) {
	s.InputSection = nil
	s.SectionFragment = frag
}

func GetSymbolByName(ctx *Context, name string) *Symbol {
	if sym, ok := ctx.SymbolMap[name]; ok {
		return sym
	}
	ctx.SymbolMap[name] = NewSymbol(name)
	return ctx.SymbolMap[name]
}

func (s *Symbol) ElfSym() *Sym {
	utils.Assert(s.SymIdx < len(s.File.ElfSyms))
	return &s.File.ElfSyms[s.SymIdx]
}

func (s *Symbol) Clear() {
	s.File = nil
	s.InputSection = nil
	s.SymIdx = -1
}
