// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package simd

// matmul function, generic
func matmul(dst, m, n []float32, mr, mc, nr, nc int) {
	for i := 0; i < mr; i++ {
		y := dst[i*nc : (i+1)*nc]
		for l, a := range m[i*mc : (i+1)*mc] {
			axpy(n[l*nc:(l+1)*nc], y, a)
		}
	}
}

// axpy function, generic
func axpy(x, y []float32, alpha float32) {
	_ = y[len(x)-1] // remove bounds checks
	for i, v := range x {
		y[i] += alpha * v
	}
}
