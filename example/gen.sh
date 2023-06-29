#!/bin/bash

go run ../cmd/gocc/. matmul_avx.c --arch avx2 -O3 -o simd --package simd
go run ../cmd/gocc/. matmul_neon.c --arch neon -O3 -o simd --package simd
go run ../cmd/gocc/. matmul_apple.c --arch apple -O3 -o simd --package simd