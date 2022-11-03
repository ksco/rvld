package linker

type ContextArgs struct {
	Output       string
	Emulation    MachineType
	LibraryPaths []string
}

type Context struct {
	Args      ContextArgs
	Objs      []*ObjectFile
	SymbolMap map[string]*Symbol
}

func NewContext() *Context {
	return &Context{
		Args: ContextArgs{
			Output:    "a.out",
			Emulation: MachineTypeNone,
		},
		SymbolMap: make(map[string]*Symbol),
	}
}
