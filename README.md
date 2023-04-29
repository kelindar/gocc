## Compiling

```bash
sudo apt install build-essential
bash -c "$(wget -O - https://apt.llvm.org/llvm.sh)"
```

## Cross-Compiling for ARM from x86_64

I'm using the following toolchain to cross-compile for ARM from my x86_64 machine (Ubuntu).

```bash
sudo apt install -qy gcc-arm-none-eabi
sudo apt install -qy gcc-arm-linux-gnueabihf
sudo apt install -qy binutils-aarch64-linux-gnu gcc-aarch64-linux-gnu g++-aarch64-linux-gnu
```

```bash
clang-15 --target=arm-linux-gnueabihf -march=armv7-a -mfpu=neon-vfpv4 -mfloat-abi=hard -S -c floats_neon.c -o output.asm
```
