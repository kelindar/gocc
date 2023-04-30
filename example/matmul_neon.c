// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

#include <arm_neon.h>
#include <stdint.h>

void f32_axpy(const float *x, float *y, const uint64_t size, const float alpha) {
    float32x4_t alpha_vec = vdupq_n_f32(alpha);
    for (uint64_t i = 0; (i + 3) < size; i += 4) {
        float32x4_t y_vec = vld1q_f32(y + i);
        float32x4_t x_vec = vld1q_f32(x + i);
        float32x4_t out = vmlaq_f32(y_vec, x_vec, alpha_vec);
        vst1q_f32(y+i, out);
    }

    // Process the tail of the vector if the size is not divisible by 4.
    uint64_t tail = size % 4;
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