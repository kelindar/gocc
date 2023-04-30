package simd

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

/*
cpu: 13th Gen Intel(R) Core(TM) i7-13700K
BenchmarkAXPY/std-24         	287770197	         4.048 ns/op	       0 B/op	       0 allocs/op
BenchmarkAXPY/asm-24         	422536102	         2.870 ns/op	       0 B/op	       0 allocs/op
*/
func BenchmarkAXPY(b *testing.B) {
	x := []float32{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	y := []float32{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}

	b.Run("std", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			axpy(x, y, 3)
		}
	})

	b.Run("asm", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f32_axpy(
				unsafe.Pointer(&x[0]),
				unsafe.Pointer(&y[0]),
				4, 3.0,
			)
		}
	})
}

/*
cpu: 13th Gen Intel(R) Core(TM) i7-13700K
BenchmarkMatmul/4x4-std-24         	28915801	        41.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/4x4-asm-24         	48979990	        24.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/8x8-std-24         	 5381164	       223.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/8x8-asm-24         	38095237	        31.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/16x16-std-24       	 1000000	      1206 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/16x16-asm-24       	 3680982	       329.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/32x32-std-24       	  155844	      7597 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/32x32-asm-24       	  648651	      1802 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/64x64-std-24       	   15513	     77741 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/64x64-asm-24       	  143712	      8420 ns/op	       0 B/op	       0 allocs/op
*/
func BenchmarkMatmul(b *testing.B) {
	for _, size := range []int{4, 8, 16, 32, 64} {
		m := newTestMatrix(size, size)
		n := newTestMatrix(size, size)
		o := newTestMatrix(size, size)

		b.Run(fmt.Sprintf("%dx%d-std", size, size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				matmul(o.Data, m.Data, n.Data, m.Rows, m.Cols, n.Rows, n.Cols)
			}
		})

		b.Run(fmt.Sprintf("%dx%d-asm", size, size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f32_matmul(
					unsafe.Pointer(&o.Data[0]), unsafe.Pointer(&m.Data[0]), unsafe.Pointer(&n.Data[0]),
					dimensionsOf(m.Rows, m.Cols, n.Rows, n.Cols),
				)
			}
		})
	}
}

func TestGenericMatmul(t *testing.T) {
	x := []float32{1, 2, 3, 4}
	y := []float32{5, 6, 7, 8}
	o := make([]float32, 4)

	matmul(o, x, y, 2, 2, 2, 2)
	assert.Equal(t, []float32{19, 22, 43, 50}, o)
}

func TestMatmul(t *testing.T) {
	x := []float32{1, 2, 3, 4}
	y := []float32{5, 6, 7, 8}
	o := make([]float32, 4)

	f32_matmul(
		unsafe.Pointer(&o[0]), unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]),
		dimensionsOf(2, 2, 2, 2),
	)

	assert.Equal(t, []float32{19, 22, 43, 50}, o)
}

// newTestMatrix creates a new matrix
func newTestMatrix(r, c int) *Matrix {
	mx := NewMatrix(r, c, nil)
	for i := 0; i < len(mx.Data); i++ {
		mx.Data[i] = 2
	}
	return &mx
}
