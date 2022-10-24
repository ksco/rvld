package linker

import "bytes"

func CheckMagic(contents []byte) bool {
	return bytes.HasPrefix(contents, []byte("\177ELF"))
}
