// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO(rsc): Rewrite all nn(SP) references into name+(nn-8)(FP)
// so that go vet can check that they are correct.

#include "textflag.h"
#include "funcdata.h"

//
// System call support for Harvey
//
//func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err ErrorString)
//func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err ErrorString)
//func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr)
//func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr)

TEXT	·Syscall(SB),NOSPLIT,$0-64
	CALL	runtime·entersyscall(SB)
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	trap+0(FP), AX
	SYSCALL
	MOVQ	AX, r1+32(FP)
	MOVQ	$0, r2+40(FP)
	LEAQ	runtime·emptystring(SB), SI
	CMPL	AX, $-1
	JNE	ok1

//	LEAQ	runtime·emptystring(SB), SI
//	CMPQ	AX, $0xfffffffffffff001
//	JLS	ok1

	SUBQ	$16, SP
	CALL	runtime·m_errstr(SB)
	MOVQ	SP, SI
	ADDQ	$16, SP
ok1:
	MOVQ	0(SI), AX
	MOVQ	8(SI), DX
	MOVQ	AX, err_base+48(FP)
	MOVQ	DX, err_len+56(FP)
	CALL	runtime·exitsyscall(SB)
	RET


TEXT	·Syscall6(SB),NOSPLIT,$0-88
	CALL	runtime·entersyscall(SB)
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	a4+32(FP), R10
	MOVQ	a5+40(FP), R8
	MOVQ	a6+48(FP), R9
	MOVQ	trap+0(FP), AX
	SYSCALL
	MOVQ	AX, r1+56(FP)
	MOVQ	$0, r2+64(FP)
	LEAQ	runtime·emptystring(SB), SI
	CMPL	AX, $-1
	JNE	ok2

//	LEAQ	runtime·emptystring(SB), SI
//	CMPQ	AX, $0xfffffffffffff001
//	JLS	ok2

	SUBQ	$16, SP
	CALL	runtime·m_errstr(SB)
	MOVQ	SP, SI
	ADDQ	$16, SP
ok2:
	MOVQ	0(SI), AX
	MOVQ	8(SI), DX
	MOVQ	AX, err_base+72(FP)
	MOVQ	DX, err_len+80(FP)
	CALL	runtime·exitsyscall(SB)
	RET

TEXT	·RawSyscall(SB),NOSPLIT,$0-56
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	trap+0(FP), AX
	SYSCALL
	MOVQ	AX, r1+32(FP)
	MOVQ	$0, r2+40(FP)
	MOVQ	$0, err+48(FP)
	RET

TEXT	·RawSyscall6(SB),NOSPLIT,$0-80
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	a4+32(FP), R10
	MOVQ	a5+40(FP), R8
	MOVQ	a6+48(FP), R9
	MOVQ	trap+0(FP), AX
	SYSCALL
	MOVQ	AX, r1+56(FP)
	MOVQ	$0, r2+64(FP)
	MOVQ	$0, err+72(FP)
	RET

TEXT	·Seekcall(SB),NOSPLIT,$0-64
	CALL	runtime·entersyscall(SB)
	LEAQ	r1+32(FP), DI
	MOVQ	a1+8(FP), SI
	MOVQ	a2+16(FP), DX
	MOVQ	a3+24(FP), R10
	MOVQ	trap+0(FP), AX
	SYSCALL
	MOVQ	$0, r2+40(FP)
	LEAQ	runtime·emptystring(SB), SI
	CMPL	AX, $-1
	JNE	ok3

//	CMPQ	AX, $0xfffffffffffff001
//	LEAQ	runtime·emptystring(SB), SI
//	JLS	ok3

	MOVQ	AX, r1+32(FP)

	SUBQ	$16, SP
	CALL	runtime·m_errstr(SB)
	MOVQ	SP, SI
	ADDQ	$16, SP
ok3:
	MOVQ	0(SI), AX
	MOVQ	8(SI), DX
	MOVQ	AX, err_base+48(FP)
	MOVQ	DX, err_len+56(FP)
	CALL	runtime·exitsyscall(SB)
	RET

//func exit(code int)
// Import runtime·exit for cleanly exiting.
TEXT ·exit(SB),NOSPLIT,$8-8
	NO_LOCAL_POINTERS
	MOVQ	code+0(FP), AX
	MOVQ	AX, 0(SP)
	CALL	runtime·exit(SB)
	RET
