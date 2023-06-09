//go:build !noasm && darwin && arm64
// AUTO-GENERATED BY GOCC -- DO NOT EDIT

TEXT ·f32_axpy(SB), $0-32
	MOVD x+0(FP), R0
	MOVD y+8(FP), R1
	MOVD size+16(FP), R2
	MOVD alpha+24(FP), R3
	WORD $0xa9bf7bfd      // stp	x29, x30, [sp, #-16]!           ; 16-byte Folded Spill
	WORD $0x910003fd      // mov	x29, sp
	WORD $0xf100105f      // cmp	x2, #4
	WORD $0x54000163      // b.lo	LBB0_3
	WORD $0x52800068      // mov	w8, #3
	WORD $0xaa0003e9      // mov	x9, x0
	WORD $0xaa0103ea      // mov	x10, x1

BB0_2:
	WORD $0x3cc10521 // ldr	q1, [x9], #16
	WORD $0x3dc00142 // ldr	q2, [x10]
	WORD $0x4f801022 // fmla.4s	v2, v1, v0[0]
	WORD $0x3c810542 // str	q2, [x10], #16
	WORD $0x91001108 // add	x8, x8, #4
	WORD $0xeb02011f // cmp	x8, x2
	WORD $0x54ffff43 // b.lo	LBB0_2

BB0_3:
	WORD $0xf240045f // tst	x2, #0x3
	WORD $0x540001e0 // b.eq	LBB0_7
	WORD $0x927ef448 // and	x8, x2, #0xfffffffffffffffc
	WORD $0xeb02011f // cmp	x8, x2
	WORD $0x54000182 // b.hs	LBB0_7
	WORD $0xcb080048 // sub	x8, x2, x8
	WORD $0xd37ef449 // lsl	x9, x2, #2
	WORD $0x927ced2a // and	x10, x9, #0xfffffffffffffff0
	WORD $0x8b0a0029 // add	x9, x1, x10
	WORD $0x8b0a000a // add	x10, x0, x10

BB0_6:
	WORD $0xbc404541 // ldr	s1, [x10], #4
	WORD $0xbd400122 // ldr	s2, [x9]
	WORD $0x1f000821 // fmadd	s1, s1, s0, s2
	WORD $0xbc004521 // str	s1, [x9], #4
	WORD $0xf1000508 // subs	x8, x8, #1
	WORD $0x54ffff61 // b.ne	LBB0_6

BB0_7:
	WORD $0xa8c17bfd // ldp	x29, x30, [sp], #16             ; 16-byte Folded Reload
	WORD $0xd65f03c0 // ret

TEXT ·f32_matmul(SB), $0-32
	MOVD dst+0(FP), R0
	MOVD m+8(FP), R1
	MOVD n+16(FP), R2
	MOVD dims+24(FP), R3
	WORD $0xa9bf7bfd     // stp	x29, x30, [sp, #-16]!           ; 16-byte Folded Spill
	WORD $0x910003fd     // mov	x29, sp
	WORD $0xf2403c68     // ands	x8, x3, #0xffff
	WORD $0x54000ce0     // b.eq	LBB1_25
	WORD $0xd3507c69     // ubfx	x9, x3, #16, #16
	WORD $0xd370fc6a     // lsr	x10, x3, #48
	WORD $0xf250047f     // tst	x3, #0x3000000000000
	WORD $0x927e354b     // and	x11, x10, #0xfffc
	WORD $0xfa4b1140     // ccmp	x10, x11, #0, ne
	WORD $0x1a9f97ec     // cset	w12, hi
	WORD $0xb4000c09     // cbz	x9, LBB1_25
	WORD $0xd372fc6d     // lsr	x13, x3, #50
	WORD $0xb40004ad     // cbz	x13, LBB1_12
	WORD $0x3600086c     // tbz	w12, #0, LBB1_19
	WORD $0xd280000c     // mov	x12, #0
	WORD $0xd37ef54d     // lsl	x13, x10, #2

BB1_5:
	WORD $0xd280000e // mov	x14, #0
	WORD $0x9b097d8f // mul	x15, x12, x9
	WORD $0xaa0203f0 // mov	x16, x2

BB1_6:
	WORD $0x8b0f01d1 // add	x17, x14, x15
	WORD $0xbc717820 // ldr	s0, [x1, x17, lsl #2]
	WORD $0xaa0003f1 // mov	x17, x0
	WORD $0xaa1003e3 // mov	x3, x16
	WORD $0x52800064 // mov	w4, #3

BB1_7:
	WORD $0x3cc10461 // ldr	q1, [x3], #16
	WORD $0x3dc00222 // ldr	q2, [x17]
	WORD $0x4f801022 // fmla.4s	v2, v1, v0[0]
	WORD $0x3c810622 // str	q2, [x17], #16
	WORD $0x91001084 // add	x4, x4, #4
	WORD $0xeb0a009f // cmp	x4, x10
	WORD $0x54ffff43 // b.lo	LBB1_7
	WORD $0xaa0b03f1 // mov	x17, x11

BB1_9:
	WORD $0xd37ef623 // lsl	x3, x17, #2
	WORD $0xbc636a01 // ldr	s1, [x16, x3]
	WORD $0xbc636802 // ldr	s2, [x0, x3]
	WORD $0x1f000821 // fmadd	s1, s1, s0, s2
	WORD $0xbc236801 // str	s1, [x0, x3]
	WORD $0x91000631 // add	x17, x17, #1
	WORD $0xeb11015f // cmp	x10, x17
	WORD $0x54ffff21 // b.ne	LBB1_9
	WORD $0x910005ce // add	x14, x14, #1
	WORD $0x8b0d0210 // add	x16, x16, x13
	WORD $0xeb0901df // cmp	x14, x9
	WORD $0x54fffd01 // b.ne	LBB1_6
	WORD $0x9100058c // add	x12, x12, #1
	WORD $0x8b0d0000 // add	x0, x0, x13
	WORD $0xeb08019f // cmp	x12, x8
	WORD $0x54fffc21 // b.ne	LBB1_5
	WORD $0x1400003a // b	LBB1_25

BB1_12:
	WORD $0x3400072c // cbz	w12, LBB1_25
	WORD $0xd280000c // mov	x12, #0
	WORD $0xcb0b014b // sub	x11, x10, x11
	WORD $0xd36efc6d // lsr	x13, x3, #46
	WORD $0x927c35ae // and	x14, x13, #0x3fff0
	WORD $0x8b0e000d // add	x13, x0, x14
	WORD $0xd37ef54a // lsl	x10, x10, #2
	WORD $0x8b0e004e // add	x14, x2, x14

BB1_14:
	WORD $0xd280000f // mov	x15, #0
	WORD $0x9b097d90 // mul	x16, x12, x9
	WORD $0xaa0e03f1 // mov	x17, x14

BB1_15:
	WORD $0x8b1001e0 // add	x0, x15, x16
	WORD $0xbc607820 // ldr	s0, [x1, x0, lsl #2]
	WORD $0xaa1103e0 // mov	x0, x17
	WORD $0xaa0d03e2 // mov	x2, x13
	WORD $0xaa0b03e3 // mov	x3, x11

BB1_16:
	WORD $0xbc404401 // ldr	s1, [x0], #4
	WORD $0xbd400042 // ldr	s2, [x2]
	WORD $0x1f000821 // fmadd	s1, s1, s0, s2
	WORD $0xbc004441 // str	s1, [x2], #4
	WORD $0xf1000463 // subs	x3, x3, #1
	WORD $0x54ffff61 // b.ne	LBB1_16
	WORD $0x910005ef // add	x15, x15, #1
	WORD $0x8b0a0231 // add	x17, x17, x10
	WORD $0xeb0901ff // cmp	x15, x9
	WORD $0x54fffe41 // b.ne	LBB1_15
	WORD $0x9100058c // add	x12, x12, #1
	WORD $0x8b0a01ad // add	x13, x13, x10
	WORD $0xeb08019f // cmp	x12, x8
	WORD $0x54fffd61 // b.ne	LBB1_14
	WORD $0x1400001b // b	LBB1_25

BB1_19:
	WORD $0xd280000b // mov	x11, #0
	WORD $0xd37ef54c // lsl	x12, x10, #2

BB1_20:
	WORD $0xd280000d // mov	x13, #0
	WORD $0x9b097d6e // mul	x14, x11, x9
	WORD $0xaa0203ef // mov	x15, x2

BB1_21:
	WORD $0x8b0e01b0 // add	x16, x13, x14
	WORD $0x8b100830 // add	x16, x1, x16, lsl #2
	WORD $0x4d40ca00 // ld1r.4s	{ v0 }, [x16]
	WORD $0xaa0003f0 // mov	x16, x0
	WORD $0xaa0f03f1 // mov	x17, x15
	WORD $0x52800063 // mov	w3, #3

BB1_22:
	WORD $0x3cc10621 // ldr	q1, [x17], #16
	WORD $0x3dc00202 // ldr	q2, [x16]
	WORD $0x4e20cc22 // fmla.4s	v2, v1, v0
	WORD $0x3c810602 // str	q2, [x16], #16
	WORD $0x91001063 // add	x3, x3, #4
	WORD $0xeb0a007f // cmp	x3, x10
	WORD $0x54ffff43 // b.lo	LBB1_22
	WORD $0x910005ad // add	x13, x13, #1
	WORD $0x8b0c01ef // add	x15, x15, x12
	WORD $0xeb0901bf // cmp	x13, x9
	WORD $0x54fffe01 // b.ne	LBB1_21
	WORD $0x9100056b // add	x11, x11, #1
	WORD $0x8b0c0000 // add	x0, x0, x12
	WORD $0xeb08017f // cmp	x11, x8
	WORD $0x54fffd21 // b.ne	LBB1_20

BB1_25:
	WORD $0xa8c17bfd // ldp	x29, x30, [sp], #16             ; 16-byte Folded Reload
	WORD $0xd65f03c0 // ret
