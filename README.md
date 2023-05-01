# GOCC: Compile C to Go Assembly

This utility transpiles C code to Go assembly. It uses the LLVM toolchain to compile C code to assembly and machine code and generates Go assembly from it, as well as the corresponding Go stubs. This is useful for certain features such as using intrinsics, which are not supported by the Go ecosystem. The example folder includes matrix multiplication using intrinsics compiled to ARM Linux, x86_x64 and Apple Silicon.

## Features

- Only requires `clang` and `objdump` to be installed in order to compile.
- Auto-detects the appropriate version of `clang` and `objdump` to use.
- Supports cross-compilation.
- Auto-generates Go stubs for the C functions by parsing C code.
- Automatically formats go assembly using `asmfmt`.

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
