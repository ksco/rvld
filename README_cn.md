# rvld

[English](README.md) | 中文版

rvld 是一个针对 RV64GC 架构的最小链接器实现，主要是出于教学目的。rvld 开发过程中大量参考了 [rui314/mold](https://github.com/rui314/mold) 的源代码，是它的衍生作品，同样在 [GNU AGPL v3](LICENSE) 协议下进行分发。

rvld 可以静态链接一个简单的 C 程序（例如下面示例中的 Hello world），并生成可运行的二进制。

```bash
cat <<EOF | $CC -o a.o -c -xc -static -
#include <stdio.h>
int main() {
  printf("Hello, World.\n");
  return 0;
}
EOF

$CC -B. -s -static a.o -o out
qemu-riscv64 out

# Hello, World.
```

rvld 只有 2000 行左右的 Go 代码，没有除标准库外的外部依赖。PLCT 实验室基于此项目推出了一门“从零开始实现链接器”的公开课，课程官网：[https://ksco.cc/rvld](https://ksco.cc/rvld)。

## 注意

rvld 可能不适应于高版本工具链，请考虑使用 Docker 环境以便于复现。

```bash
docker run -u root --volume ~/Developer:/code -it golang:bullseye

# Inside the docker container:
apt update
apt install -y gcc-10-riscv64-linux-gnu qemu-user
ln -sf /usr/bin/riscv64-linux-gnu-gcc-10 /usr/bin/riscv64-linux-gnu-gcc
```
