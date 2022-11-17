package linker

import "bytes"

func CheckMagic(contents []byte) bool {
	return bytes.HasPrefix(contents, []byte("\177ELF"))
}

func WriteMagic(contents []byte) {
	copy(contents, "\177ELF")
}
