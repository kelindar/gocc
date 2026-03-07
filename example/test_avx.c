// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

#include <stdint.h>



// ---------------------------------- Uint8 ----------------------------------

void _uint8_sum(uint8_t *input, uint8_t *result, uint64_t size) {
    uint8_t sum = 0.0;
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; ++i) {
        sum += input[i];
    }
    *result = sum;
}

void _uint8_min(uint8_t *input, uint8_t *result, uint64_t size) {
    uint8_t min = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] < min) {
            min = input[i];
        }
    }
    *result = min;
}

void _uint8_max(uint8_t *input, uint8_t *result, uint64_t size) {
    uint8_t max = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] > max) {
            max = input[i];
        }
    }
    *result = max;
}

void _uint8_add(uint8_t *input1, uint8_t *input2, uint8_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] + input2[i];
    }
}

void _uint8_sub(uint8_t *input1, uint8_t *input2, uint8_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] - input2[i];
    }
}

void _uint8_mul(uint8_t *input1, uint8_t *input2, uint8_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] * input2[i];
    }
}

void _uint8_div(uint8_t *input1, uint8_t *input2, uint8_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] / input2[i];
    }
}

// ---------------------------------- Uint16 ----------------------------------

void _uint16_sum(uint16_t *input, uint16_t *result, uint64_t size) {
    uint16_t sum = 0.0;
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; ++i) {
        sum += input[i];
    }
    *result = sum;
}

void _uint16_min(uint16_t *input, uint16_t *result, uint64_t size) {
    uint16_t min = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] < min) {
            min = input[i];
        }
    }
    *result = min;
}

void _uint16_max(uint16_t *input, uint16_t *result, uint64_t size) {
    uint16_t max = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] > max) {
            max = input[i];
        }
    }
    *result = max;
}

void _uint16_add(uint16_t *input1, uint16_t *input2, uint16_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] + input2[i];
    }
}

void _uint16_sub(uint16_t *input1, uint16_t *input2, uint16_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] - input2[i];
    }
}

void _uint16_mul(uint16_t *input1, uint16_t *input2, uint16_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] * input2[i];
    }
}

void _uint16_div(uint16_t *input1, uint16_t *input2, uint16_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] / input2[i];
    }
}

// ---------------------------------- Uint32 ----------------------------------

void _uint32_sum(uint32_t *input, uint32_t *result, uint64_t size) {
    uint32_t sum = 0.0;
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; ++i) {
        sum += input[i];
    }
    *result = sum;
}

void _uint32_min(uint32_t *input, uint32_t *result, uint64_t size) {
    uint32_t min = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] < min) {
            min = input[i];
        }
    }
    *result = min;
}

void _uint32_max(uint32_t *input, uint32_t *result, uint64_t size) {
    uint32_t max = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] > max) {
            max = input[i];
        }
    }
    *result = max;
}

void _uint32_add(uint32_t *input1, uint32_t *input2, uint32_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] + input2[i];
    }
}

void _uint32_sub(uint32_t *input1, uint32_t *input2, uint32_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] - input2[i];
    }
}

void _uint32_mul(uint32_t *input1, uint32_t *input2, uint32_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] * input2[i];
    }
}

void _uint32_div(uint32_t *input1, uint32_t *input2, uint32_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] / input2[i];
    }
}

// ---------------------------------- Uint64 ----------------------------------

void _uint64_sum(uint64_t *input, uint64_t *result, uint64_t size) {
    uint64_t sum = 0.0;
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; ++i) {
        sum += input[i];
    }
    *result = sum;
}

void _uint64_min(uint64_t *input, uint64_t *result, uint64_t size) {
    uint64_t min = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] < min) {
            min = input[i];
        }
    }
    *result = min;
}

void _uint64_max(uint64_t *input, uint64_t *result, uint64_t size) {
    uint64_t max = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] > max) {
            max = input[i];
        }
    }
    *result = max;
}

void _uint64_add(uint64_t *input1, uint64_t *input2, uint64_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] + input2[i];
    }
}

void _uint64_sub(uint64_t *input1, uint64_t *input2, uint64_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] - input2[i];
    }
}

void _uint64_mul(uint64_t *input1, uint64_t *input2, uint64_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] * input2[i];
    }
}

void _uint64_div(uint64_t *input1, uint64_t *input2, uint64_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] / input2[i];
    }
}

// ---------------------------------- Int8 ----------------------------------

void _int8_sum(int8_t *input, int8_t *result, uint64_t size) {
    int8_t sum = 0.0;
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; ++i) {
        sum += input[i];
    }
    *result = sum;
}

void _int8_min(int8_t *input, int8_t *result, uint64_t size) {
    int8_t min = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] < min) {
            min = input[i];
        }
    }
    *result = min;
}

void _int8_max(int8_t *input, int8_t *result, uint64_t size) {
    int8_t max = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] > max) {
            max = input[i];
        }
    }
    *result = max;
}

void _int8_add(int8_t *input1, int8_t *input2, int8_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] + input2[i];
    }
}

void _int8_sub(int8_t *input1, int8_t *input2, int8_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] - input2[i];
    }
}

void _int8_mul(int8_t *input1, int8_t *input2, int8_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] * input2[i];
    }
}

void _int8_div(int8_t *input1, int8_t *input2, int8_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] / input2[i];
    }
}

// ---------------------------------- Int16 ----------------------------------

void _int16_sum(int16_t *input, int16_t *result, uint64_t size) {
    int16_t sum = 0.0;
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; ++i) {
        sum += input[i];
    }
    *result = sum;
}

void _int16_min(int16_t *input, int16_t *result, uint64_t size) {
    int16_t min = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] < min) {
            min = input[i];
        }
    }
    *result = min;
}

void _int16_max(int16_t *input, int16_t *result, uint64_t size) {
    int16_t max = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] > max) {
            max = input[i];
        }
    }
    *result = max;
}

void _int16_add(int16_t *input1, int16_t *input2, int16_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] + input2[i];
    }
}

void _int16_sub(int16_t *input1, int16_t *input2, int16_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] - input2[i];
    }
}

void _int16_mul(int16_t *input1, int16_t *input2, int16_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] * input2[i];
    }
}

void _int16_div(int16_t *input1, int16_t *input2, int16_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] / input2[i];
    }
}

// ---------------------------------- Int32 ----------------------------------

void _int32_sum(int32_t *input, int32_t *result, uint64_t size) {
    int32_t sum = 0.0;
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; ++i) {
        sum += input[i];
    }
    *result = sum;
}

void _int32_min(int32_t *input, int32_t *result, uint64_t size) {
    int32_t min = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] < min) {
            min = input[i];
        }
    }
    *result = min;
}

void _int32_max(int32_t *input, int32_t *result, uint64_t size) {
    int32_t max = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] > max) {
            max = input[i];
        }
    }
    *result = max;
}

void _int32_add(int32_t *input1, int32_t *input2, int32_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] + input2[i];
    }
}

void _int32_sub(int32_t *input1, int32_t *input2, int32_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] - input2[i];
    }
}

void _int32_mul(int32_t *input1, int32_t *input2, int32_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] * input2[i];
    }
}

void _int32_div(int32_t *input1, int32_t *input2, int32_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] / input2[i];
    }
}

// ---------------------------------- Int64 ----------------------------------

void _int64_sum(int64_t *input, int64_t *result, uint64_t size) {
    int64_t sum = 0.0;
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; ++i) {
        sum += input[i];
    }
    *result = sum;
}

void _int64_min(int64_t *input, int64_t *result, uint64_t size) {
    int64_t min = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] < min) {
            min = input[i];
        }
    }
    *result = min;
}

void _int64_max(int64_t *input, int64_t *result, uint64_t size) {
    int64_t max = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] > max) {
            max = input[i];
        }
    }
    *result = max;
}

void _int64_add(int64_t *input1, int64_t *input2, int64_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] + input2[i];
    }
}

void _int64_sub(int64_t *input1, int64_t *input2, int64_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] - input2[i];
    }
}

void _int64_mul(int64_t *input1, int64_t *input2, int64_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] * input2[i];
    }
}

void _int64_div(int64_t *input1, int64_t *input2, int64_t *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] / input2[i];
    }
}

// ---------------------------------- Float32 ----------------------------------

void _float32_sum(float *input, float *result, uint64_t size) {
    float sum = 0.0;
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; ++i) {
        sum += input[i];
    }
    *result = sum;
}

void _float32_min(float *input, float *result, uint64_t size) {
    float min = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] < min) {
            min = input[i];
        }
    }
    *result = min;
}

void _float32_max(float *input, float *result, uint64_t size) {
    float max = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] > max) {
            max = input[i];
        }
    }
    *result = max;
}

void _float32_add(float *input1, float *input2, float *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] + input2[i];
    }
}

void _float32_sub(float *input1, float *input2, float *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] - input2[i];
    }
}

void _float32_mul(float *input1, float *input2, float *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] * input2[i];
    }
}

void _float32_div(float *input1, float *input2, float *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] / input2[i];
    }
}

// ---------------------------------- Float64 ----------------------------------

void _float64_sum(double *input, double *result, uint64_t size) {
    double sum = 0.0;
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; ++i) {
        sum += input[i];
    }
    *result = sum;
}

void _float64_min(double *input, double *result, uint64_t size) {
    double min = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] < min) {
            min = input[i];
        }
    }
    *result = min;
}

void _float64_max(double *input, double *result, uint64_t size) {
    double max = input[0];
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        if (input[i] > max) {
            max = input[i];
        }
    }
    *result = max;
}

void _float64_add(double *input1, double *input2, double *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] + input2[i];
    }
}

void _float64_sub(double *input1, double *input2, double *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] - input2[i];
    }
}

void _float64_mul(double *input1, double *input2, double *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] * input2[i];
    }
}

void _float64_div(double *input1, double *input2, double *output, uint64_t size) {
    #pragma clang loop vectorize(enable) interleave(enable)
    for (int i = 0; i < (int)size; i++) {
        output[i] = input1[i] / input2[i];
    }
}

