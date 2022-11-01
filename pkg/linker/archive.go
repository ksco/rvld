package linker

import "github.com/ksco/rvld/pkg/utils"

func ReadArchiveMembers(file *File) []*File {
	utils.Assert(GetFileType(file.Contents) == FileTypeArchive)

	pos := 8
	var strTab []byte
	var files []*File
	for len(file.Contents)-pos > 1 {
		if pos%2 == 1 {
			pos++
		}
		hdr := utils.Read[ArHdr](file.Contents[pos:])
		dataStart := pos + ArHdrSize
		pos = dataStart + hdr.GetSize()
		dataEnd := pos
		contents := file.Contents[dataStart:dataEnd]

		if hdr.IsSymtab() {
			continue
		} else if hdr.IsStrtab() {
			strTab = contents
			continue
		}

		files = append(files, &File{
			Name:     hdr.ReadName(strTab),
			Contents: contents,
			Parent:   file,
		})
	}

	return files
}
