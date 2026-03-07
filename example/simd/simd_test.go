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
BenchmarkMatmul/4x4-std-24         	24242570	        49.69 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/4x4-asm-24         	26667140	        45.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/8x8-std-24         	 4545457	       265.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/8x8-asm-24         	21428494	        50.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/16x16-std-24       	 1000000	      1267 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/16x16-asm-24       	 7017567	       165.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/32x32-std-24       	  129031	      9893 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/32x32-asm-24       	 1854714	       623.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/64x64-std-24       	   18644	     64486 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatmul/64x64-asm-24       	  510646	      2408 ns/op	       0 B/op	       0 allocs/op
*/
func BenchmarkMatmul(b *testing.B) {
	for _, size := range []int{4, 8, 16, 32, 64, 128, 256, 512} {
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

func TestMatmulNative(t *testing.T) {
	x := []float32{1, 2, 3, 4}
	y := []float32{5, 6, 7, 8}
	o := make([]float32, 4)

	f32_matmul(
		unsafe.Pointer(&o[0]), unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]),
		dimensionsOf(2, 2, 2, 2),
	)

	assert.Equal(t, []float32{19, 22, 43, 50}, o)
}

func TestMatmul(t *testing.T) {
	x := Matrix{Rows: 2, Cols: 2, Data: []float32{1, 2, 3, 4}}
	y := Matrix{Rows: 2, Cols: 2, Data: []float32{5, 6, 7, 8}}
	o := Matrix{Rows: 2, Cols: 2, Data: make([]float32, 4)}

	Matmul(&o, &x, &y)
	assert.Equal(t, []float32{19, 22, 43, 50}, o.Data)
}

// newTestMatrix creates a new matrix
func newTestMatrix(r, c int) *Matrix {
	mx := NewMatrix(r, c, nil)
	for i := 0; i < len(mx.Data); i++ {
		mx.Data[i] = 2
	}
	return &mx
}

func TestUintMul(t *testing.T) {
	input1 := makeVector[uint8](70)
	input2 := makeVector[uint8](70)

	// generic implementation
	dst1 := make([]uint8, 70)
	mul(dst1, input1, input2)

	// asm implementation
	dst2 := make([]uint8, 70)
	_uint8_mul(unsafe.Pointer(&input1[0]), unsafe.Pointer(&input2[0]), unsafe.Pointer(&(dst2)[0]), uint64(len(dst2)))

	assert.Equal(t, dst1, dst2)
}

func TestIntMin(t *testing.T) {
	input := makeVector[int8](70)

	// generic implementation
	min1 := min(input)

	// asm implementation
	min2 := int8(0)
	_int8_min(unsafe.Pointer(&input[0]), unsafe.Pointer(&min2), uint64(len(input)))

	assert.EqualValues(t, min1, min2)
}

func TestIntDiv(t *testing.T) {
	input1 := makeVector[int8](70)
	input2 := makeVector[int8](70)

	// generic implementation
	dst1 := make([]int8, 70)
	div(dst1, input1, input2)

	// asm implementation
	dst2 := make([]int8, 70)
	_int8_div(unsafe.Pointer(&input1[0]), unsafe.Pointer(&input2[0]), unsafe.Pointer(&dst2[0]), uint64(len(dst2)))

	assert.InDeltaSlice(t, dst1, dst2, 0.01)
}

// Mul multiplies input1 by input2 and writes back the result into dst slice
func mul(dst, input1, input2 []uint8) []uint8 {
	for i, v := range input1 {
		dst[i] = v * input2[i]
	}
	return dst
}

// Min finds the minimum value in the input slice and writes back the result into dst slice
func min(input []int8) int8 {
	min := input[0]
	for _, v := range input[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Div divides input1 by input2 and writes back the result into dst slice
func div(dst, input1, input2 []int8) []int8 {
	for i, v := range input1 {
		dst[i] = v / input2[i]
	}
	return dst
}

// makeVector generates a test vector
func makeVector[T int8 | uint8](count int) []T {
	arr := make([]T, count)
	for i := 0; i < count; i++ {
		arr[i] = T((i % 100) + 1)
	}
	return arr
}
