//go:build !noasm && amd64
// AUTO-GENERATED BY GOCC -- DO NOT EDIT

TEXT ·f32_axpy(SB), $0-32
	MOVQ x+0(FP), DI
	MOVQ y+8(FP), SI
	MOVQ size+16(FP), DX
	MOVQ alpha+24(FP), CX
	BYTE $0x55                   // push	rbp
	WORD $0x8948; BYTE $0xe5     // mov	rbp, rsp
	LONG $0xf8e48348             // and	rsp, -8
	LONG $0x08fa8348             // cmp	rdx, 8
	JB   LBB0_5
	LONG $0x187de2c4; BYTE $0xc8 // vbroadcastss	ymm1, xmm0
	LONG $0xf8428d48             // lea	rax, [rdx - 8]
	WORD $0x8949; BYTE $0xc0     // mov	r8, rax
	LONG $0x03e8c149             // shr	r8, 3
	WORD $0xff49; BYTE $0xc0     // inc	r8
	LONG $0x08f88348             // cmp	rax, 8
	JAE  LBB0_12
	WORD $0xc931                 // xor	ecx, ecx
	JMP  LBB0_3

LBB0_12:
	WORD $0x894c; BYTE $0xc0 // mov	rax, r8
	LONG $0xfee08348         // and	rax, -2
	WORD $0xc931             // xor	ecx, ecx

LBB0_13:
	LONG $0x1410fcc5; BYTE $0x8f               // vmovups	ymm2, ymmword ptr [rdi + 4*rcx]
	LONG $0xa875e2c4; WORD $0x8e14             // vfmadd213ps	ymm2, ymm1, ymmword ptr [rsi + 4*rcx]
	LONG $0x1411fcc5; BYTE $0x8e               // vmovups	ymmword ptr [rsi + 4*rcx], ymm2
	LONG $0x5410fcc5; WORD $0x208f             // vmovups	ymm2, ymmword ptr [rdi + 4*rcx + 32]
	LONG $0xa875e2c4; WORD $0x8e54; BYTE $0x20 // vfmadd213ps	ymm2, ymm1, ymmword ptr [rsi + 4*rcx + 32]
	LONG $0x5411fcc5; WORD $0x208e             // vmovups	ymmword ptr [rsi + 4*rcx + 32], ymm2
	LONG $0x10c18348                           // add	rcx, 16
	LONG $0xfec08348                           // add	rax, -2
	JNE  LBB0_13

LBB0_3:
	LONG $0x01c0f641               // test	r8b, 1
	JE   LBB0_5
	LONG $0x1410fcc5; BYTE $0x8f   // vmovups	ymm2, ymmword ptr [rdi + 4*rcx]
	LONG $0xa86de2c4; WORD $0x8e0c // vfmadd213ps	ymm1, ymm2, ymmword ptr [rsi + 4*rcx]
	LONG $0x0c11fcc5; BYTE $0x8e   // vmovups	ymmword ptr [rsi + 4*rcx], ymm1

LBB0_5:
	WORD $0xc2f6; BYTE $0x07       // test	dl, 7
	JE   LBB0_11
	WORD $0x8948; BYTE $0xd0       // mov	rax, rdx
	LONG $0xf8e08348               // and	rax, -8
	WORD $0x3948; BYTE $0xd0       // cmp	rax, rdx
	JAE  LBB0_11
	WORD $0x8948; BYTE $0xc1       // mov	rcx, rax
	WORD $0xf748; BYTE $0xd1       // not	rcx
	WORD $0xc2f6; BYTE $0x01       // test	dl, 1
	JE   LBB0_9
	LONG $0x0c10fac5; BYTE $0x87   // vmovss	xmm1, dword ptr [rdi + 4*rax]
	LONG $0xa979e2c4; WORD $0x860c // vfmadd213ss	xmm1, xmm0, dword ptr [rsi + 4*rax]
	LONG $0x0c11fac5; BYTE $0x86   // vmovss	dword ptr [rsi + 4*rax], xmm1
	LONG $0x01c88348               // or	rax, 1

LBB0_9:
	WORD $0x0148; BYTE $0xd1 // add	rcx, rdx
	JE   LBB0_11

LBB0_10:
	LONG $0x0c10fac5; BYTE $0x87               // vmovss	xmm1, dword ptr [rdi + 4*rax]
	LONG $0xa979e2c4; WORD $0x860c             // vfmadd213ss	xmm1, xmm0, dword ptr [rsi + 4*rax]
	LONG $0x0c11fac5; BYTE $0x86               // vmovss	dword ptr [rsi + 4*rax], xmm1
	LONG $0x4c10fac5; WORD $0x0487             // vmovss	xmm1, dword ptr [rdi + 4*rax + 4]
	LONG $0xa979e2c4; WORD $0x864c; BYTE $0x04 // vfmadd213ss	xmm1, xmm0, dword ptr [rsi + 4*rax + 4]
	LONG $0x4c11fac5; WORD $0x0486             // vmovss	dword ptr [rsi + 4*rax + 4], xmm1
	LONG $0x02c08348                           // add	rax, 2
	WORD $0x3948; BYTE $0xd0                   // cmp	rax, rdx
	JB   LBB0_10

LBB0_11:
	WORD $0x8948; BYTE $0xec // mov	rsp, rbp
	BYTE $0x5d               // pop	rbp
	WORD $0xf8c5; BYTE $0x77 // vzeroupper
	BYTE $0xc3               // ret

TEXT ·f32_matmul(SB), $0-32
	MOVQ dst+0(FP), DI
	MOVQ m+8(FP), SI
	MOVQ n+16(FP), DX
	MOVQ dims+24(FP), CX
	BYTE $0x55                                 // push	rbp
	WORD $0x8948; BYTE $0xe5                   // mov	rbp, rsp
	WORD $0x5741                               // push	r15
	WORD $0x5641                               // push	r14
	WORD $0x5541                               // push	r13
	WORD $0x5441                               // push	r12
	BYTE $0x53                                 // push	rbx
	LONG $0xf8e48348                           // and	rsp, -8
	LONG $0x18ec8148; WORD $0x0001; BYTE $0x00 // sub	rsp, 280
	LONG $0x24548948; BYTE $0x18               // mov	qword ptr [rsp + 24], rdx
	WORD $0x8948; BYTE $0xc8                   // mov	rax, rcx
	LONG $0xffff2548; WORD $0x0000             // and	rax, 65535
	LONG $0x24448948; BYTE $0x40               // mov	qword ptr [rsp + 64], rax
	JE   LBB1_13
	WORD $0x8949; BYTE $0xc8                   // mov	r8, rcx
	WORD $0x8948; BYTE $0xc8                   // mov	rax, rcx
	LONG $0x30e8c148                           // shr	rax, 48
	LONG $0x24048948                           // mov	qword ptr [rsp], rax
	JE   LBB1_13
	WORD $0x8949; BYTE $0xf3                   // mov	r11, rsi
	WORD $0x894d; BYTE $0xc7                   // mov	r15, r8
	LONG $0x10efc149                           // shr	r15, 16
	LONG $0xf7b70f41                           // movzx	esi, r15w
	LONG $0xf8e78141; WORD $0x00ff; BYTE $0x00 // and	r15d, 65528
	LONG $0x240c8b48                           // mov	rcx, qword ptr [rsp]
	LONG $0x02f98348                           // cmp	rcx, 2
	LONG $0x000001ba; BYTE $0x00               // mov	edx, 1
	LONG $0xd1430f48                           // cmovae	rdx, rcx
	QUAD $0x000000c824948948                   // mov	qword ptr [rsp + 200], rdx
	WORD $0x8948; BYTE $0xf2                   // mov	rdx, rsi
	WORD $0x294c; BYTE $0xfa                   // sub	rdx, r15
	LONG $0x24548948; BYTE $0x28               // mov	qword ptr [rsp + 40], rdx
	LONG $0xf0e28348                           // and	rdx, -16
	LONG $0x075f8d49                           // lea	rbx, [r15 + 7]
	LONG $0xd9af0f48                           // imul	rbx, rcx
	QUAD $0x000000b0249c8948                   // mov	qword ptr [rsp + 176], rbx
	WORD $0x8948; BYTE $0xcb                   // mov	rbx, rcx
	LONG $0x06e3c148                           // shl	rbx, 6
	QUAD $0x000000a0249c8948                   // mov	qword ptr [rsp + 160], rbx
	LONG $0x065f8d49                           // lea	rbx, [r15 + 6]
	LONG $0xd9af0f48                           // imul	rbx, rcx
	QUAD $0x00000090249c8948                   // mov	qword ptr [rsp + 144], rbx
	LONG $0x05478d49                           // lea	rax, [r15 + 5]
	LONG $0xc1af0f48                           // imul	rax, rcx
	LONG $0x24448948; BYTE $0x50               // mov	qword ptr [rsp + 80], rax
	LONG $0x04478d49                           // lea	rax, [r15 + 4]
	LONG $0xc1af0f48                           // imul	rax, rcx
	QUAD $0x0000011024848948                   // mov	qword ptr [rsp + 272], rax
	LONG $0x03478d49                           // lea	rax, [r15 + 3]
	LONG $0xc1af0f48                           // imul	rax, rcx
	QUAD $0x000000f824848948                   // mov	qword ptr [rsp + 248], rax
	LONG $0x026f8d4d                           // lea	r13, [r15 + 2]
	LONG $0xe9af0f4c                           // imul	r13, rcx
	LONG $0x01678d4d                           // lea	r12, [r15 + 1]
	LONG $0xe1af0f4c                           // imul	r12, rcx
	LONG $0x13e8c141                           // shr	r8d, 19
	WORD $0x8949; BYTE $0xc9                   // mov	r9, rcx
	LONG $0xc8af0f4d                           // imul	r9, r8
	LONG $0x05e0c149                           // shl	r8, 5
	LONG $0x181c8d4b                           // lea	rbx, [r8 + r11]
	LONG $0x20c38348                           // add	rbx, 32
	LONG $0x245c8948; BYTE $0x10               // mov	qword ptr [rsp + 16], rbx
	LONG $0x05e1c149                           // shl	r9, 5
	QUAD $0x00000088248c894c                   // mov	qword ptr [rsp + 136], r9
	LONG $0x0f5f8d49                           // lea	rbx, [r15 + 15]
	LONG $0xd9af0f48                           // imul	rbx, rcx
	LONG $0x245c8948; BYTE $0x78               // mov	qword ptr [rsp + 120], rbx
	LONG $0x0e778d4d                           // lea	r14, [r15 + 14]
	LONG $0xf1af0f4c                           // imul	r14, rcx
	LONG $0x0d578d4d                           // lea	r10, [r15 + 13]
	LONG $0xd1af0f4c                           // imul	r10, rcx
	LONG $0x0c478d4d                           // lea	r8, [r15 + 12]
	LONG $0xc1af0f4c                           // imul	r8, rcx
	LONG $0x0b478d49                           // lea	rax, [r15 + 11]
	LONG $0xc1af0f48                           // imul	rax, rcx
	LONG $0x0a5f8d49                           // lea	rbx, [r15 + 10]
	LONG $0xd9af0f48                           // imul	rbx, rcx
	LONG $0x245c8948; BYTE $0x70               // mov	qword ptr [rsp + 112], rbx
	LONG $0x095f8d49                           // lea	rbx, [r15 + 9]
	LONG $0xd9af0f48                           // imul	rbx, rcx
	LONG $0x245c8948; BYTE $0x68               // mov	qword ptr [rsp + 104], rbx
	LONG $0x085f8d49                           // lea	rbx, [r15 + 8]
	LONG $0xd9af0f48                           // imul	rbx, rcx
	LONG $0x245c8948; BYTE $0x60               // mov	qword ptr [rsp + 96], rbx
	LONG $0x24548948; BYTE $0x20               // mov	qword ptr [rsp + 32], rdx
	WORD $0x014c; BYTE $0xfa                   // add	rdx, r15
	LONG $0x24548948; BYTE $0x58               // mov	qword ptr [rsp + 88], rdx
	QUAD $0x00000000b5148d48                   // lea	rdx, [4*rsi]
	LONG $0x24548948; BYTE $0x38               // mov	qword ptr [rsp + 56], rdx
	QUAD $0x000000008d0c8d48                   // lea	rcx, [4*rcx]
	QUAD $0x000000c0248c8948                   // mov	qword ptr [rsp + 192], rcx
	WORD $0xd231                               // xor	edx, edx
	QUAD $0x000000a824bc894c                   // mov	qword ptr [rsp + 168], r15
	QUAD $0x0000009824bc8948                   // mov	qword ptr [rsp + 152], rdi
	QUAD $0x0000008024b48948                   // mov	qword ptr [rsp + 128], rsi
	QUAD $0x0000010824a4894c                   // mov	qword ptr [rsp + 264], r12
	QUAD $0x0000010024ac894c                   // mov	qword ptr [rsp + 256], r13
	QUAD $0x000000f024848948                   // mov	qword ptr [rsp + 240], rax
	QUAD $0x000000e82484894c                   // mov	qword ptr [rsp + 232], r8
	QUAD $0x000000e02494894c                   // mov	qword ptr [rsp + 224], r10
	QUAD $0x000000d824b4894c                   // mov	qword ptr [rsp + 216], r14
	QUAD $0x000000b0248c8b4c                   // mov	r9, qword ptr [rsp + 176]
	LONG $0x246c8b4c; BYTE $0x50               // mov	r13, qword ptr [rsp + 80]
	LONG $0x24748b4c; BYTE $0x78               // mov	r14, qword ptr [rsp + 120]
	LONG $0x24648b4c; BYTE $0x70               // mov	r12, qword ptr [rsp + 112]
	LONG $0x24448b48; BYTE $0x68               // mov	rax, qword ptr [rsp + 104]
	LONG $0x24448b4c; BYTE $0x60               // mov	r8, qword ptr [rsp + 96]
	QUAD $0x0000010824948b4c                   // mov	r10, qword ptr [rsp + 264]
	JMP  LBB1_3

LBB1_12:
	LONG $0x24548b48; BYTE $0x48 // mov	rdx, qword ptr [rsp + 72]
	WORD $0xff48; BYTE $0xc2     // inc	rdx
	LONG $0x244c8b48; BYTE $0x38 // mov	rcx, qword ptr [rsp + 56]
	LONG $0x244c0148; BYTE $0x10 // add	qword ptr [rsp + 16], rcx
	WORD $0x0149; BYTE $0xcb     // add	r11, rcx
	LONG $0x24543b48; BYTE $0x40 // cmp	rdx, qword ptr [rsp + 64]
	JE   LBB1_13

LBB1_3:
	LONG $0x24548948; BYTE $0x48 // mov	qword ptr [rsp + 72], rdx
	LONG $0x14af0f48; BYTE $0x24 // imul	rdx, qword ptr [rsp]
	QUAD $0x000000d024948948     // mov	qword ptr [rsp + 208], rdx
	LONG $0x244c8b48; BYTE $0x18 // mov	rcx, qword ptr [rsp + 24]
	LONG $0x244c8948; BYTE $0x30 // mov	qword ptr [rsp + 48], rcx
	WORD $0xc931                 // xor	ecx, ecx
	LONG $0x244c8948; BYTE $0x08 // mov	qword ptr [rsp + 8], rcx
	QUAD $0x000000b8249c894c     // mov	qword ptr [rsp + 184], r11
	JMP  LBB1_6

LBB1_5:
	QUAD $0x000000d0248c8b48       // mov	rcx, qword ptr [rsp + 208]
	LONG $0x245c8b48; BYTE $0x08   // mov	rbx, qword ptr [rsp + 8]
	LONG $0x0b148d48               // lea	rdx, [rbx + rcx]
	LONG $0x0411fac5; BYTE $0x97   // vmovss	dword ptr [rdi + 4*rdx], xmm0
	WORD $0xff48; BYTE $0xc3       // inc	rbx
	LONG $0x24448348; WORD $0x0430 // add	qword ptr [rsp + 48], 4
	WORD $0x8948; BYTE $0xd9       // mov	rcx, rbx
	LONG $0x245c8948; BYTE $0x08   // mov	qword ptr [rsp + 8], rbx
	QUAD $0x000000c8249c3b48       // cmp	rbx, qword ptr [rsp + 200]
	JE   LBB1_12

LBB1_6:
	LONG $0xc057f8c5               // vxorps	xmm0, xmm0, xmm0
	WORD $0x394c; BYTE $0xfe       // cmp	rsi, r15
	JBE  LBB1_5
	LONG $0xc057f8c5               // vxorps	xmm0, xmm0, xmm0
	WORD $0x894c; BYTE $0xfa       // mov	rdx, r15
	LONG $0x247c8348; WORD $0x1028 // cmp	qword ptr [rsp + 40], 16
	JB   LBB1_11
	LONG $0xc057f8c5               // vxorps	xmm0, xmm0, xmm0
	LONG $0x247c8b4c; BYTE $0x10   // mov	r15, qword ptr [rsp + 16]
	LONG $0x24548b48; BYTE $0x30   // mov	rdx, qword ptr [rsp + 48]
	LONG $0x247c8b48; BYTE $0x20   // mov	rdi, qword ptr [rsp + 32]
	LONG $0xc957f0c5               // vxorps	xmm1, xmm1, xmm1
	QUAD $0x000000a0249c8b48       // mov	rbx, qword ptr [rsp + 160]
	QUAD $0x00000090249c8b4c       // mov	r11, qword ptr [rsp + 144]
	QUAD $0x0000008824b48b48       // mov	rsi, qword ptr [rsp + 136]

LBB1_9:
	QUAD $0x00000110248c8b48                   // mov	rcx, qword ptr [rsp + 272]
	LONG $0x1410fac5; BYTE $0x8a               // vmovss	xmm2, dword ptr [rdx + 4*rcx]
	LONG $0x2169a3c4; WORD $0xaa14; BYTE $0x10 // vinsertps	xmm2, xmm2, dword ptr [rdx + 4*r13], 16
	LONG $0x2169a3c4; WORD $0x9a14; BYTE $0x20 // vinsertps	xmm2, xmm2, dword ptr [rdx + 4*r11], 32
	LONG $0x2169a3c4; WORD $0x8a14; BYTE $0x30 // vinsertps	xmm2, xmm2, dword ptr [rdx + 4*r9], 48
	LONG $0x1c10fac5; BYTE $0x32               // vmovss	xmm3, dword ptr [rdx + rsi]
	LONG $0x2161a3c4; WORD $0x921c; BYTE $0x10 // vinsertps	xmm3, xmm3, dword ptr [rdx + 4*r10], 16
	QUAD $0x00000100248c8b48                   // mov	rcx, qword ptr [rsp + 256]
	LONG $0x2161e3c4; WORD $0x8a1c; BYTE $0x20 // vinsertps	xmm3, xmm3, dword ptr [rdx + 4*rcx], 32
	QUAD $0x000000f8248c8b48                   // mov	rcx, qword ptr [rsp + 248]
	LONG $0x2161e3c4; WORD $0x8a1c; BYTE $0x30 // vinsertps	xmm3, xmm3, dword ptr [rdx + 4*rcx], 48
	QUAD $0x000000e8248c8b48                   // mov	rcx, qword ptr [rsp + 232]
	LONG $0x2410fac5; BYTE $0x8a               // vmovss	xmm4, dword ptr [rdx + 4*rcx]
	QUAD $0x000000e0248c8b48                   // mov	rcx, qword ptr [rsp + 224]
	LONG $0x2159e3c4; WORD $0x8a24; BYTE $0x10 // vinsertps	xmm4, xmm4, dword ptr [rdx + 4*rcx], 16
	QUAD $0x000000d8248c8b48                   // mov	rcx, qword ptr [rsp + 216]
	LONG $0x2159e3c4; WORD $0x8a24; BYTE $0x20 // vinsertps	xmm4, xmm4, dword ptr [rdx + 4*rcx], 32
	LONG $0x2159a3c4; WORD $0xb224; BYTE $0x30 // vinsertps	xmm4, xmm4, dword ptr [rdx + 4*r14], 48
	LONG $0x107aa1c4; WORD $0x822c             // vmovss	xmm5, dword ptr [rdx + 4*r8]
	LONG $0x2151e3c4; WORD $0x822c; BYTE $0x10 // vinsertps	xmm5, xmm5, dword ptr [rdx + 4*rax], 16
	LONG $0x2151a3c4; WORD $0xa22c; BYTE $0x20 // vinsertps	xmm5, xmm5, dword ptr [rdx + 4*r12], 32
	LONG $0x1865e3c4; WORD $0x01d2             // vinsertf128	ymm2, ymm3, xmm2, 1
	QUAD $0x000000f0248c8b48                   // mov	rcx, qword ptr [rsp + 240]
	LONG $0x2151e3c4; WORD $0x8a1c; BYTE $0x30 // vinsertps	xmm3, xmm5, dword ptr [rdx + 4*rcx], 48
	LONG $0x1865e3c4; WORD $0x01dc             // vinsertf128	ymm3, ymm3, xmm4, 1
	LONG $0xb86dc2c4; WORD $0xe047             // vfmadd231ps	ymm0, ymm2, ymmword ptr [r15 - 32]
	LONG $0xb865c2c4; BYTE $0x0f               // vfmadd231ps	ymm1, ymm3, ymmword ptr [r15]
	WORD $0x0148; BYTE $0xda                   // add	rdx, rbx
	LONG $0x40c78349                           // add	r15, 64
	LONG $0xf0c78348                           // add	rdi, -16
	JNE  LBB1_9
	LONG $0xc058f4c5                           // vaddps	ymm0, ymm1, ymm0
	LONG $0x197de3c4; WORD $0x01c1             // vextractf128	xmm1, ymm0, 1
	LONG $0xc158f8c5                           // vaddps	xmm0, xmm0, xmm1
	LONG $0x0579e3c4; WORD $0x01c8             // vpermilpd	xmm1, xmm0, 1
	LONG $0xc158f8c5                           // vaddps	xmm0, xmm0, xmm1
	LONG $0xc816fac5                           // vmovshdup	xmm1, xmm0
	LONG $0xc158fac5                           // vaddss	xmm0, xmm0, xmm1
	LONG $0x24548b48; BYTE $0x58               // mov	rdx, qword ptr [rsp + 88]
	LONG $0x247c8b48; BYTE $0x20               // mov	rdi, qword ptr [rsp + 32]
	LONG $0x247c3948; BYTE $0x28               // cmp	qword ptr [rsp + 40], rdi
	QUAD $0x0000009824bc8b48                   // mov	rdi, qword ptr [rsp + 152]
	QUAD $0x000000a824bc8b4c                   // mov	r15, qword ptr [rsp + 168]
	QUAD $0x000000b8249c8b4c                   // mov	r11, qword ptr [rsp + 184]
	QUAD $0x0000008024b48b48                   // mov	rsi, qword ptr [rsp + 128]
	JE   LBB1_5

LBB1_11:
	LONG $0x240c8b48             // mov	rcx, qword ptr [rsp]
	LONG $0xcaaf0f48             // imul	rcx, rdx
	LONG $0x244c0348; BYTE $0x08 // add	rcx, qword ptr [rsp + 8]
	LONG $0x245c8b48; BYTE $0x18 // mov	rbx, qword ptr [rsp + 24]
	LONG $0x8b1c8d48             // lea	rbx, [rbx + 4*rcx]
	QUAD $0x000000c0248c8b48     // mov	rcx, qword ptr [rsp + 192]

LBB1_4:
	LONG $0x0b10fac5               // vmovss	xmm1, dword ptr [rbx]
	LONG $0xb971c2c4; WORD $0x9304 // vfmadd231ss	xmm0, xmm1, dword ptr [r11 + 4*rdx]
	WORD $0xff48; BYTE $0xc2       // inc	rdx
	WORD $0x0148; BYTE $0xcb       // add	rbx, rcx
	WORD $0x3948; BYTE $0xd6       // cmp	rsi, rdx
	JNE  LBB1_4
	JMP  LBB1_5

LBB1_13:
	LONG $0xd8658d48         // lea	rsp, [rbp - 40]
	BYTE $0x5b               // pop	rbx
	WORD $0x5c41             // pop	r12
	WORD $0x5d41             // pop	r13
	WORD $0x5e41             // pop	r14
	WORD $0x5f41             // pop	r15
	BYTE $0x5d               // pop	rbp
	WORD $0xf8c5; BYTE $0x77 // vzeroupper
	BYTE $0xc3               // ret
