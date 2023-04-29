#!/bin/bash

go run ../. floats_avx.c -O3 -mavx
# go run ../. floats_neon.c -O3