package linker

import (
	"debug/elf"
	"github.com/ksco/rvld/pkg/utils"
)

type MachineType = uint8

const (
	MachineTypeNone    MachineType = iota
	MachineTypeRISCV64 MachineType = iota
)

func GetMachineTypeFromContents(contents []byte) MachineType {
	ft := GetFileType(contents)

	switch ft {
	case FileTypeObject:
		machine := utils.Read[uint16](contents[18:])
		if machine == uint16(elf.EM_RISCV) {
			class := elf.Class(contents[4])
			switch class {
			case elf.ELFCLASS64:
				return MachineTypeRISCV64
			}
		}
	}

	return MachineTypeNone
}

type MachineTypeStringer struct {
	MachineType
}

func (m MachineTypeStringer) String() string {
	switch m.MachineType {
	case MachineTypeRISCV64:
		return "riscv64"
	}

	utils.Assert(m.MachineType == MachineTypeNone)
	return "none"
}
