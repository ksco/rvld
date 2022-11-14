package linker

import (
	"debug/elf"
	"strings"
)

var prefixes = []string{
	".text.", ".data.rel.ro.", ".data.", ".rodata.", ".bss.rel.ro.", ".bss.",
	".init_array.", ".fini_array.", ".tbss.", ".tdata.", ".gcc_except_table.",
	".ctors.", ".dtors.",
}

func GetOutputName(name string, flags uint64) string {
	if (name == ".rodata" || strings.HasPrefix(name, ".rodata.")) &&
		flags&uint64(elf.SHF_MERGE) != 0 {
		if flags&uint64(elf.SHF_STRINGS) != 0 {
			return ".rodata.str"
		} else {
			return ".rodata.cst"
		}
	}

	for _, prefix := range prefixes {
		stem := prefix[:len(prefix)-1]
		if name == stem || strings.HasPrefix(name, prefix) {
			return stem
		}
	}

	return name
}
