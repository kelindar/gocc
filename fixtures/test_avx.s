	.text
	.intel_syntax noprefix
	.file	"test_avx.c"
	.section	.rodata.cst32,"aM",@progbits,32
	.p2align	5                               # -- Begin function uint8_mul
.LCPI0_0:
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.short	255                             # 0xff
	.text
	.globl	uint8_mul
	.p2align	4, 0x90
	.type	uint8_mul,@function
uint8_mul:                              # @uint8_mul
# %bb.0:
	push	rbp
	mov	rbp, rsp
	and	rsp, -8
	test	ecx, ecx
	jle	.LBB0_18
# %bb.1:
	mov	r8d, ecx
	cmp	r8, 16
	jae	.LBB0_3
# %bb.2:
	xor	r9d, r9d
.LBB0_14:
	sub	ecx, r9d
	mov	r10, r9
	not	r10
	add	r10, r8
	and	rcx, 3
	je	.LBB0_16
	.p2align	4, 0x90
.LBB0_15:                               # =>This Inner Loop Header: Depth=1
	movzx	eax, byte ptr [rsi + r9]
	mul	byte ptr [rdi + r9]
	mov	byte ptr [rdx + r9], al
	inc	r9
	dec	rcx
	jne	.LBB0_15
.LBB0_16:
	cmp	r10, 3
	jb	.LBB0_18
	.p2align	4, 0x90
.LBB0_17:                               # =>This Inner Loop Header: Depth=1
	movzx	eax, byte ptr [rsi + r9]
	mul	byte ptr [rdi + r9]
	mov	byte ptr [rdx + r9], al
	movzx	eax, byte ptr [rsi + r9 + 1]
	mul	byte ptr [rdi + r9 + 1]
	mov	byte ptr [rdx + r9 + 1], al
	movzx	eax, byte ptr [rsi + r9 + 2]
	mul	byte ptr [rdi + r9 + 2]
	mov	byte ptr [rdx + r9 + 2], al
	movzx	eax, byte ptr [rsi + r9 + 3]
	mul	byte ptr [rdi + r9 + 3]
	mov	byte ptr [rdx + r9 + 3], al
	add	r9, 4
	cmp	r8, r9
	jne	.LBB0_17
.LBB0_18:
	mov	rsp, rbp
	pop	rbp
	vzeroupper
	ret
.LBB0_3:
	mov	rax, rdx
	sub	rax, rdi
	xor	r9d, r9d
	cmp	rax, 128
	jb	.LBB0_14
# %bb.4:
	mov	rax, rdx
	sub	rax, rsi
	cmp	rax, 128
	jb	.LBB0_14
# %bb.5:
	cmp	r8d, 128
	jae	.LBB0_7
# %bb.6:
	xor	r9d, r9d
	jmp	.LBB0_11
.LBB0_7:
	mov	r10d, ecx
	and	r10d, 127
	mov	r9, r8
	sub	r9, r10
	xor	eax, eax
	vmovdqa	ymm0, ymmword ptr [rip + .LCPI0_0] # ymm0 = [255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255]
	.p2align	4, 0x90
.LBB0_8:                                # =>This Inner Loop Header: Depth=1
	vmovdqu	ymm3, ymmword ptr [rdi + rax]
	vmovdqu	ymm4, ymmword ptr [rdi + rax + 32]
	vmovdqu	ymm5, ymmword ptr [rdi + rax + 64]
	vmovdqu	ymm1, ymmword ptr [rdi + rax + 96]
	vmovdqu	ymm6, ymmword ptr [rsi + rax]
	vmovdqu	ymm7, ymmword ptr [rsi + rax + 32]
	vmovdqu	ymm8, ymmword ptr [rsi + rax + 64]
	vmovdqu	ymm2, ymmword ptr [rsi + rax + 96]
	vpunpckhbw	ymm9, ymm3, ymm3        # ymm9 = ymm3[8,8,9,9,10,10,11,11,12,12,13,13,14,14,15,15,24,24,25,25,26,26,27,27,28,28,29,29,30,30,31,31]
	vpunpckhbw	ymm10, ymm6, ymm6       # ymm10 = ymm6[8,8,9,9,10,10,11,11,12,12,13,13,14,14,15,15,24,24,25,25,26,26,27,27,28,28,29,29,30,30,31,31]
	vpmullw	ymm9, ymm10, ymm9
	vpand	ymm9, ymm9, ymm0
	vpunpcklbw	ymm3, ymm3, ymm3        # ymm3 = ymm3[0,0,1,1,2,2,3,3,4,4,5,5,6,6,7,7,16,16,17,17,18,18,19,19,20,20,21,21,22,22,23,23]
	vpunpcklbw	ymm6, ymm6, ymm6        # ymm6 = ymm6[0,0,1,1,2,2,3,3,4,4,5,5,6,6,7,7,16,16,17,17,18,18,19,19,20,20,21,21,22,22,23,23]
	vpmullw	ymm3, ymm6, ymm3
	vpand	ymm3, ymm3, ymm0
	vpackuswb	ymm3, ymm3, ymm9
	vpunpckhbw	ymm6, ymm4, ymm4        # ymm6 = ymm4[8,8,9,9,10,10,11,11,12,12,13,13,14,14,15,15,24,24,25,25,26,26,27,27,28,28,29,29,30,30,31,31]
	vpunpckhbw	ymm9, ymm7, ymm7        # ymm9 = ymm7[8,8,9,9,10,10,11,11,12,12,13,13,14,14,15,15,24,24,25,25,26,26,27,27,28,28,29,29,30,30,31,31]
	vpmullw	ymm6, ymm9, ymm6
	vpand	ymm6, ymm6, ymm0
	vpunpcklbw	ymm4, ymm4, ymm4        # ymm4 = ymm4[0,0,1,1,2,2,3,3,4,4,5,5,6,6,7,7,16,16,17,17,18,18,19,19,20,20,21,21,22,22,23,23]
	vpunpcklbw	ymm7, ymm7, ymm7        # ymm7 = ymm7[0,0,1,1,2,2,3,3,4,4,5,5,6,6,7,7,16,16,17,17,18,18,19,19,20,20,21,21,22,22,23,23]
	vpmullw	ymm4, ymm7, ymm4
	vpand	ymm4, ymm4, ymm0
	vpackuswb	ymm4, ymm4, ymm6
	vpunpckhbw	ymm6, ymm5, ymm5        # ymm6 = ymm5[8,8,9,9,10,10,11,11,12,12,13,13,14,14,15,15,24,24,25,25,26,26,27,27,28,28,29,29,30,30,31,31]
	vpunpckhbw	ymm7, ymm8, ymm8        # ymm7 = ymm8[8,8,9,9,10,10,11,11,12,12,13,13,14,14,15,15,24,24,25,25,26,26,27,27,28,28,29,29,30,30,31,31]
	vpmullw	ymm6, ymm7, ymm6
	vpand	ymm6, ymm6, ymm0
	vpunpcklbw	ymm5, ymm5, ymm5        # ymm5 = ymm5[0,0,1,1,2,2,3,3,4,4,5,5,6,6,7,7,16,16,17,17,18,18,19,19,20,20,21,21,22,22,23,23]
	vpunpcklbw	ymm7, ymm8, ymm8        # ymm7 = ymm8[0,0,1,1,2,2,3,3,4,4,5,5,6,6,7,7,16,16,17,17,18,18,19,19,20,20,21,21,22,22,23,23]
	vpmullw	ymm5, ymm7, ymm5
	vpand	ymm5, ymm5, ymm0
	vpackuswb	ymm5, ymm5, ymm6
	vpunpckhbw	ymm6, ymm1, ymm1        # ymm6 = ymm1[8,8,9,9,10,10,11,11,12,12,13,13,14,14,15,15,24,24,25,25,26,26,27,27,28,28,29,29,30,30,31,31]
	vpunpckhbw	ymm7, ymm2, ymm2        # ymm7 = ymm2[8,8,9,9,10,10,11,11,12,12,13,13,14,14,15,15,24,24,25,25,26,26,27,27,28,28,29,29,30,30,31,31]
	vpmullw	ymm6, ymm7, ymm6
	vpand	ymm6, ymm6, ymm0
	vpunpcklbw	ymm1, ymm1, ymm1        # ymm1 = ymm1[0,0,1,1,2,2,3,3,4,4,5,5,6,6,7,7,16,16,17,17,18,18,19,19,20,20,21,21,22,22,23,23]
	vpunpcklbw	ymm2, ymm2, ymm2        # ymm2 = ymm2[0,0,1,1,2,2,3,3,4,4,5,5,6,6,7,7,16,16,17,17,18,18,19,19,20,20,21,21,22,22,23,23]
	vpmullw	ymm1, ymm2, ymm1
	vpand	ymm1, ymm1, ymm0
	vpackuswb	ymm1, ymm1, ymm6
	vmovdqu	ymmword ptr [rdx + rax], ymm3
	vmovdqu	ymmword ptr [rdx + rax + 32], ymm4
	vmovdqu	ymmword ptr [rdx + rax + 64], ymm5
	vmovdqu	ymmword ptr [rdx + rax + 96], ymm1
	sub	rax, -128
	cmp	r9, rax
	jne	.LBB0_8
# %bb.9:
	test	r10, r10
	je	.LBB0_18
# %bb.10:
	cmp	r10d, 16
	jb	.LBB0_14
.LBB0_11:
	mov	rax, r9
	mov	r10d, ecx
	and	r10d, 15
	mov	r9, r8
	sub	r9, r10
	vmovdqa	ymm0, ymmword ptr [rip + .LCPI0_0] # ymm0 = [255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255]
	.p2align	4, 0x90
.LBB0_12:                               # =>This Inner Loop Header: Depth=1
	vpmovzxbw	ymm1, xmmword ptr [rdi + rax] # ymm1 = mem[0],zero,mem[1],zero,mem[2],zero,mem[3],zero,mem[4],zero,mem[5],zero,mem[6],zero,mem[7],zero,mem[8],zero,mem[9],zero,mem[10],zero,mem[11],zero,mem[12],zero,mem[13],zero,mem[14],zero,mem[15],zero
	vpmovzxbw	ymm2, xmmword ptr [rsi + rax] # ymm2 = mem[0],zero,mem[1],zero,mem[2],zero,mem[3],zero,mem[4],zero,mem[5],zero,mem[6],zero,mem[7],zero,mem[8],zero,mem[9],zero,mem[10],zero,mem[11],zero,mem[12],zero,mem[13],zero,mem[14],zero,mem[15],zero
	vpmullw	ymm1, ymm2, ymm1
	vpand	ymm1, ymm1, ymm0
	vextracti128	xmm2, ymm1, 1
	vpackuswb	xmm1, xmm1, xmm2
	vmovdqu	xmmword ptr [rdx + rax], xmm1
	add	rax, 16
	cmp	r9, rax
	jne	.LBB0_12
# %bb.13:
	test	r10, r10
	jne	.LBB0_14
	jmp	.LBB0_18
.Lfunc_end0:
	.size	uint8_mul, .Lfunc_end0-uint8_mul
                                        # -- End function
	.ident	"Ubuntu clang version 15.0.7"
	.section	".note.GNU-stack","",@progbits
	.addrsig
