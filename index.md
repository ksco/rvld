# 从零开始实现链接器

《从零开始实现链接器》是 PLCT Lab 推出的一门公开课，本课程还在筹备中，计划在 2023 年一季度推出，敬请期待。

在本课程中，我们会从零开始使用 Go 语言实现一个 RV64GC 架构的链接器，可以正确地链接相对简单的 C 程序。通过学习本课程，我们可以掌握链接器最核心部分的工作原理。

本课程在 [GitHub](https://github.com/ksco/rvld) 上开源。为了确保课程的顺利进行，我们提前实现了本课程中最终会实现的参考代码，放在了 [main](https://github.com/ksco/rvld/tree/main) 分支中。在 [course](https://github.com/ksco/rvld/tree/course) 分支中则按照课程记录放有每节课的代码，每节课一个 commit。

## 第一课：搭建开发环境、初始化项目、开始读取 ELF 文件。

在本节课中，我们使用 Docker 搭建了 Go 语言的开发环境，并使用 Go Mod 对项目进行了初始化。

然后我们简单介绍了 ELF 文件的结构，并开始读取 ELF 文件。

[此处插入第一课的视频链接]

参考链接：

[Executable and Linkable Format - Wikipedia](https://en.wikipedia.org/wiki/Executable_and_Linkable_Format)

[ELF64 File Header](https://fasterthanli.me/content/series/making-our-own-executable-packer/part-1/assets/elf64-file-header.bfa657ccd8ab3a7d.svg)



## 第二课：继续读取 ELF 文件

在本节课中，我们继续读取并解析了 object file 中几个重要的 section 类型。

[此处插入第二课的视频链接]

参考链接：

[Executable and Linkable Format - Wikipedia](https://en.wikipedia.org/wiki/Executable_and_Linkable_Format)

[Sections - System V Application Binary Interface 2001](https://refspecs.linuxbase.org/elf/gabi4+/ch4.sheader.html)



## 第三课：解析链接器参数

在本节课中，我们完成了链接器参数的解析。

[此处插入第三课的视频链接]
