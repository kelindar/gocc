#!/bin/bash

go run ../. floats_avx.c --arch amd64 -O3 -mavx2 -o out
# go run ../. floats_neon.c --arch arm64 -O3