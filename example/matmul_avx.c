// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

#include <stdint.h>
#include <immintrin.h>

void f32_axpy(const float *x, float *y, const uint64_t size, const float alpha) {
    __m256 a = _mm256_set1_ps(alpha);
    for (uint64_t i = 0; (i + 7) < size; i += 8) {
        __m256 y_vec = _mm256_loadu_ps(y + i);
        __m256 x_vec = _mm256_loadu_ps(x + i);
        __m256 out = _mm256_fmadd_ps(x_vec, a, y_vec);
        _mm256_storeu_ps(y+i, out);
    }

    // Process the tail of the vector if the size is not divisible by 8.
    uint64_t tail = size % 8;
    if (tail > 0) {
        for (uint64_t i = size - tail; i < size; i++) {
            y[i] += alpha * x[i];
        }
    }
}

void f32_matmul(float *dst, float *m, float *n, uint64_t dims) {
    uint64_t mr = dims & 0xFFFF;
    uint64_t mc = (dims >> 16) & 0xFFFF;
    uint64_t nr = (dims >> 32) & 0xFFFF;
    uint64_t nc = (dims >> 48) & 0xFFFF;

    for (uint64_t i = 0; i < mr; i++) {
        for (uint64_t k = 0; k < mc; k++) {
            f32_axpy(n + k*nc, dst + i*nc, nc, m[i*mc+k]);
        }
    }
}