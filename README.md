<p align="center">
<img width="250" height="110" src=".github/logo.png" border="0" alt="kelindar/gocc">
<br>
<img src="https://img.shields.io/github/go-mod/go-version/kelindar/gocc" alt="Go Version">
<a href="https://opensource.org/license/apache-2-0/"><img src="https://img.shields.io/badge/License-Apache-blue.svg" alt="License"></a>
</p>

## GOCC: Compile C to Go Assembly

This utility transpiles C code to Go assembly. It uses the LLVM toolchain to compile C code to assembly and machine code and generates Go assembly from it, as well as the corresponding Go stubs. This is useful for certain features such as using intrinsics, which are not supported by the Go ecosystem. The example folder includes matrix multiplication using intrinsics compiled to ARM Linux, x86_x64 and Apple Silicon.

## Features

- Remote compilation using a docker container with all toolchains (including Apple Silicon) pre-installed.
- Only requires `clang` and `objdump` to be installed in order to compile.
- Auto-detects the appropriate version of `clang` and `objdump` to use.
- Supports cross-compilation.
- Auto-generates Go stubs for the C functions by parsing C code.
- Automatically formats go assembly using `asmfmt`.

## Using Remotely (default)

The easiest way to use this tool is to use it remotely. This will use a docker container with all the toolchains pre-installed. This is the default mode of operation and requires no additional setup. The only requirement is to have `go` installed on your machine.

First, we need to install `gocc` command-line tool. This will install the `gocc` command in your `$GOPATH/bin` folder, or `$GOBIN`. Make sure to add it to your `$PATH` variable as well.

```
go install github.com/kelindar/gocc/cmd/gocc@latest
```

Next, you can use it to compile your C code to Go assembly. For example, to compile the `matmul_avx.c` file, you can run the following command:

```bash
gocc matmul_avx.c --arch avx2
```

## Setting up locally

Before you use it, you need to install the LLVM toolchain. On Ubuntu, you can do it with the following commands:

```bash
sudo apt install build-essential
bash -c "$(wget -O - https://apt.llvm.org/llvm.sh)"
```

For cross-compilation you will also need to install the appropriate toolchain. For example, for ARM Linux:

```bash
sudo apt install -qy binutils-aarch64-linux-gnu gcc-aarch64-linux-gnu g++-aarch64-linux-gnu
```

## Limitations

This tool does not support most of the C features, it's not a replacement for C/Go. If you are using this for production code, make sure to test the generated code thoroughly. Also, this is not meant to be general-purpose tool, but rather a tool for solving my own problems of speeding up certain routines.

- Only supports C code that can be compiled by `clang`.
- Does not support C++ code or templates for now.
- Does not support call statements, thus require you to inline your C functions
- Currently limited to 4 arguments per function and must be 64-bit.

## Resources

The ideas of this are built on top of others who have done similar things, such as

- [c2goasm](https://github.com/minio/c2goasm) and [asm2plan9s](https://github.com/minio/asm2plan9s) by Minio
- [gorse/goat by Gorse](https://github.com/gorse-io/gorse/tree/master/cmd/goat)
- [A Primer on Go Assembly](https://github.com/teh-cmc/go-internals/blob/master/chapter1_assembly_primer/README.md)
- [Go Function in Assembly](https://github.com/golang/go/files/447163/GoFunctionsInAssembly.pdf)
- [Stack frame layout on x86-64](http://eli.thegreenplace.net/2011/09/06/stack-frame-layout-on-x86-64)
- [Compiler Explorer (interactive)](https://go.godbolt.org/)
