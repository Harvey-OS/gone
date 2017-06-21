// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "funcdata.h"

//
// // System call support for HarveyPlan 9System calls for AMD64, Linux
//

// Trap # in AX, args in DI SI DX R10 R8 R9, return in AX DX
// Note that this differs from "standard" ABI convention, which
// would pass 4th arg in CX, not R10.

//func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err string)
//func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err string)
//func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr)
//func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr)
TEXT	·Syscall(SB),NOSPLIT,$0-56
	CALL	runtime·entersyscall(SB)
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	$0, R10
	MOVQ	$0, R8
	MOVQ	$0, R9
	MOVQ	trap+0(FP), AX	// syscall entry
	SYSCALL
	MOVQ	AX, r1+32(FP)
	MOVQ	$0, r2+40(FP)
	CMPL	AX, $-1
	JNE	ok3

	SUBQ	$16, SP
	CALL	runtime·errstr(SB)
	MOVQ	SP, SI
	ADDQ	$16, SP
	JMP	copyresult3
	
ok3:
	LEAQ	runtime·emptystring(SB), SI	
	
copyresult3:
	LEAQ	err+48(FP), DI

	CLD
	MOVSQ
	MOVSQ

	CALL	runtime·exitsyscall(SB)
	RET

// func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr)
TEXT ·Syscall6(SB),NOSPLIT,$0-80
	CALL	runtime·entersyscall(SB)
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	a4+32(FP), R10
	MOVQ	a5+40(FP), R8
	MOVQ	a6+48(FP), R9
	MOVQ	trap+0(FP), AX	// syscall entry
	SYSCALL
	MOVQ	AX, r1+56(FP)
	MOVQ	$0, r2+64(FP)
	CMPL	AX, $-1
	JNE	ok4
	
	SUBQ	$16, SP
	CALL	runtime·errstr(SB)
	MOVQ	SP, SI
	ADDQ	$16, SP
	JMP	copyresult4
	
ok4:
	LEAQ	runtime·emptystring(SB), SI
	
copyresult4:
	LEAQ	err+72(FP), DI

	CLD
	MOVSQ
	MOVSQ

	CALL	runtime·exitsyscall(SB)
	RET

// func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr)
TEXT ·RawSyscall(SB),NOSPLIT,$0-56
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	$0, R10
	MOVQ	$0, R8
	MOVQ	$0, R9
	MOVQ	trap+0(FP), AX	// syscall entry
	SYSCALL
	MOVQ	AX, r1+40(SP)
	MOVQ	AX, r2+48(SP)
	MOVQ	AX, err+56(SP)
 	RET

// func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr)
TEXT ·RawSyscall6(SB),NOSPLIT,$0-80
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	a4+32(FP), R10
	MOVQ	a5+40(FP), R8
	MOVQ	a6+48(FP), R9
	MOVQ	trap+0(FP), AX	// syscall entry
	SYSCALL
	MOVQ	AX, r1+56(FP)
	MOVQ	DX, r2+64(FP)
	MOVQ	$0, err+72(FP)
	RET

// func gettimeofday(tv *Timeval) (err uintptr)
TEXT ·gettimeofday(SB),NOSPLIT,$0-16
	MOVQ	tv+0(FP), DI
	MOVQ	$0, SI
	MOVQ	runtime·__vdso_gettimeofday_sym(SB), AX
	CALL	AX

	CMPQ	AX, $0xfffffffffffff001
	JLS	ok7
	NEGQ	AX
	MOVQ	AX, err+8(FP)
	RET
ok7:
	MOVQ	$0, err+8(FP)
	RET


//func exit(code int)
// Import runtime·exit for cleanly exiting.
TEXT ·exit(SB),NOSPLIT,$8-8
	NO_LOCAL_POINTERS
	MOVQ	code+0(FP), AX
	MOVQ	AX, 0(SP)
	CALL	runtime·exit(SB)
	RET

//func seek(placeholder uintptr, fd int, offset int64, whence int) (newoffset int64, err string)
//TEXT ·seek(SB),NOSPLIT,$0-56
