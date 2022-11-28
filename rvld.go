package main

import (
	"fmt"
	"github.com/ksco/rvld/pkg/linker"
	"github.com/ksco/rvld/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

var version string

func main() {
	ctx := linker.NewContext()
	remaining := parseArgs(ctx)

	if ctx.Args.Emulation == linker.MachineTypeNone {
		for _, filename := range remaining {
			if strings.HasPrefix(filename, "-") {
				continue
			}

			file := linker.MustNewFile(filename)
			ctx.Args.Emulation =
				linker.GetMachineTypeFromContents(file.Contents)
			if ctx.Args.Emulation != linker.MachineTypeNone {
				break
			}
		}
	}

	if ctx.Args.Emulation != linker.MachineTypeRISCV64 {
		utils.Fatal("unknown emulation type")
	}

	linker.ReadInputFiles(ctx, remaining)
	linker.CreateInternalFile(ctx)
	linker.ResolveSymbols(ctx)
	linker.RegisterSectionPieces(ctx)
	linker.CreateSyntheticSections(ctx)
	linker.BinSections(ctx)
	ctx.Chunks = append(ctx.Chunks, linker.CollectOutputSections(ctx)...)
	linker.ComputeSectionSizes(ctx)
	linker.SortOutputSections(ctx)

	for _, chunk := range ctx.Chunks {
		chunk.UpdateShdr(ctx)
	}

	fileSize := linker.SetOutputSectionOffsets(ctx)
	println(fileSize)

	ctx.Buf = make([]byte, fileSize)

	file, err := os.OpenFile(ctx.Args.Output, os.O_RDWR|os.O_CREATE, 0777)
	utils.MustNo(err)

	for _, chunk := range ctx.Chunks {
		chunk.CopyBuf(ctx)
	}

	_, err = file.Write(ctx.Buf)
	utils.MustNo(err)
}

func parseArgs(ctx *linker.Context) []string {
	args := os.Args[1:]

	dashes := func(name string) []string {
		if len(name) == 1 {
			return []string{"-" + name}
		}
		return []string{"-" + name, "--" + name}
	}

	arg := ""
	readArg := func(name string) bool {
		for _, opt := range dashes(name) {
			if args[0] == opt {
				if len(args) == 1 {
					utils.Fatal(fmt.Sprintf("option -%s: argument missing", name))
				}

				arg = args[1]
				args = args[2:]
				return true
			}

			prefix := opt
			if len(name) > 1 {
				prefix += "="
			}
			if strings.HasPrefix(args[0], prefix) {
				arg = args[0][len(prefix):]
				args = args[1:]
				return true
			}
		}

		return false
	}

	readFlag := func(name string) bool {
		for _, opt := range dashes(name) {
			if args[0] == opt {
				args = args[1:]
				return true
			}
		}

		return false
	}

	remaining := make([]string, 0)
	for len(args) > 0 {
		if readFlag("help") {
			fmt.Printf("usage: %s [options] file...\n", os.Args[0])
			os.Exit(0)
		}

		if readArg("o") || readArg("output") {
			ctx.Args.Output = arg
		} else if readFlag("v") || readFlag("version") {
			fmt.Printf("rvld %s\n", version)
			os.Exit(0)
		} else if readArg("m") {
			if arg == "elf64lriscv" {
				ctx.Args.Emulation = linker.MachineTypeRISCV64
			} else {
				utils.Fatal(fmt.Sprintf("unknown -m argument: %s", arg))
			}
		} else if readArg("L") {
			ctx.Args.LibraryPaths = append(ctx.Args.LibraryPaths, arg)
		} else if readArg("l") {
			remaining = append(remaining, "-l"+arg)
		} else if readArg("sysroot") ||
			readFlag("static") ||
			readArg("plugin") ||
			readArg("plugin-opt") ||
			readFlag("as-needed") ||
			readFlag("start-group") ||
			readFlag("end-group") ||
			readArg("hash-style") ||
			readArg("build-id") ||
			readFlag("s") ||
			readFlag("no-relax") {
			// Ignored
		} else {
			if args[0][0] == '-' {
				utils.Fatal(fmt.Sprintf(
					"unknown command line option: %s", args[0]))
			}
			remaining = append(remaining, args[0])
			args = args[1:]
		}
	}

	for i, path := range ctx.Args.LibraryPaths {
		ctx.Args.LibraryPaths[i] = filepath.Clean(path)
	}

	return remaining
}
