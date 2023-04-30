#!/bin/bash

# go run ../. floats_avx.c --arch amd64 -O3 -mavx2 -o simd
# go run ../. floats_neon.c --arch arm64 -O3  -o simd

go run ../. matmul_avx.c --arch amd64 -O3 -mavx2 -mfma -masm=intel -o simd
go run ../. matmul_neon.c --arch arm64 -O3 -o simd