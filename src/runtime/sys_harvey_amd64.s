// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"


// begin generated code

TEXT runtime·open(SB),NOSPLIT,$0
	MOVQ	arg0+0(FP), DI
	MOVL	arg1+8(FP), SI
	MOVQ	$2, AX
	SYSCALL
	MOVL	AX, ret+16(FP)
	RET

TEXT runtime·pread(SB),NOSPLIT,$0
	MOVL	arg0+0(FP), DI
	MOVQ	arg1+8(FP), SI
	MOVL	arg2+16(FP), DX
	MOVQ	arg3+24(FP), R10
	MOVQ	$17, AX
	SYSCALL
	MOVL	AX, ret+32(FP)
	RET

TEXT runtime·pwrite(SB),NOSPLIT,$0
	MOVL	arg0+0(FP), DI
	MOVQ	arg1+8(FP), SI
	MOVL	arg2+16(FP), DX
	MOVQ	arg3+24(FP), R10
	MOVQ	$18, AX
	SYSCALL
	MOVL	AX, ret+32(FP)
	RET

TEXT runtime·closefd(SB),NOSPLIT,$0
	MOVL	arg0+0(FP), DI
	MOVQ	$3, AX
	SYSCALL
	MOVL	AX, ret+8(FP)
	RET

TEXT runtime·exits(SB),NOSPLIT,$0
	MOVQ	arg0+0(FP), DI
	MOVQ	$4104, AX
	SYSCALL
	RET

TEXT runtime·brk_(SB),NOSPLIT,$0
	MOVQ	arg0+0(FP), DI
	MOVQ	$4120, AX
	SYSCALL
	MOVQ	AX, ret+8(FP)
	RET

TEXT runtime·sleep(SB),NOSPLIT,$0
	MOVL	arg0+0(FP), DI
	MOVQ	$4113, AX
	SYSCALL
	MOVL	AX, ret+8(FP)
	RET

TEXT runtime·plan9_semacquire(SB),NOSPLIT,$0
	MOVQ	arg0+0(FP), DI
	MOVL	arg1+8(FP), SI
	MOVQ	$4133, AX
	SYSCALL
	MOVL	AX, ret+16(FP)
	RET

TEXT runtime·plan9_tsemacquire(SB),NOSPLIT,$0
	MOVQ	arg0+0(FP), DI
	MOVL	arg1+8(FP), SI
	MOVQ	$4148, AX
	SYSCALL
	MOVL	AX, ret+16(FP)
	RET

TEXT runtime·notify(SB),NOSPLIT,$0
	MOVQ	arg0+0(FP), DI
	MOVQ	$4124, AX
	SYSCALL
	MOVL	AX, ret+8(FP)
	RET

TEXT runtime·noted(SB),NOSPLIT,$0
	MOVL	arg0+0(FP), DI
	MOVQ	$4125, AX
	SYSCALL
	MOVL	AX, ret+8(FP)
	RET

TEXT runtime·plan9_semrelease(SB),NOSPLIT,$0
	MOVQ	arg0+0(FP), DI
	MOVL	arg1+8(FP), SI
	MOVQ	$4134, AX
	SYSCALL
	MOVL	AX, ret+16(FP)
	RET

TEXT runtime·rfork(SB),NOSPLIT,$0
	MOVL	arg0+0(FP), DI
	MOVQ	$4115, AX
	SYSCALL
	MOVL	AX, ret+8(FP)
	RET

// void errstr(int8 *buf, int32 len)
TEXT runtime·errstr(SB),NOSPLIT,$0
	MOVQ	arg0+0(FP), DI
	MOVL	arg1+8(FP), SI
	MOVQ	$4137, AX
	SYSCALL
	MOVL	AX, ret+16(FP)
	RET

TEXT runtime·m_errstr(SB),NOSPLIT,$16-16
	get_tls(AX)
	MOVQ	g(AX), BX
	MOVQ	g_m(BX), BX
	MOVQ	m_errstr(BX), DI
	MOVQ	DI, 0(SP)
	MOVQ	$128, SI		// src/runtime/os2_harvey.go _ERRMAX
	MOVQ	$4137, AX
	SYSCALL
	CALL	runtime·findnull(SB)
	MOVQ	0(SP), AX
	MOVQ	8(SP), DX
	MOVQ	AX, ret_base+0(FP)
	MOVQ	DX, ret_len+8(FP)
	RET


// setldt(int entry, int address, int limit)
TEXT runtime·setldt(SB),NOSPLIT,$0
	RET

TEXT runtime·nanotime(SB),NOSPLIT,$0
	MOVQ	$4149, AX
	SYSCALL
	MOVQ	AX, ret+0(FP)
	RET

// func now() (sec int64, nsec int32)
TEXT time·now(SB),NOSPLIT,$8-12
	MOVQ	$4149, AX
	SYSCALL

	// generated code for
	//	func f(x uint64) (uint64, uint64) { return x/1000000000, x%100000000 }
	// adapted to reduce duplication
	MOVQ	AX, CX
	MOVQ	$1360296554856532783, AX
	MULQ	CX
	ADDQ	CX, DX
	RCRQ	$1, DX
	SHRQ	$29, DX
	MOVQ	DX, sec+0(FP)
	IMULQ	$1000000000, DX
	SUBQ	DX, CX
	MOVL	CX, nsec+8(FP)
	RET

// int64 seek(int32, int64, int32)
// Convenience wrapper around _seek, the actual system call.
TEXT runtime·seek(SB),NOSPLIT,$32
	LEAQ	ret+24(FP), DI
	MOVL	fd+0(FP), SI
	MOVQ	offset+8(FP), DX
	MOVL	whence+16(FP), R10
	MOVQ	$4135, AX
	SYSCALL
	CMPL	AX, $-1
	JNE	seekok
	MOVQ	$-1, ret+24(FP)
seekok:
	RET

TEXT runtime·tstart_plan9(SB),NOSPLIT,$0
	MOVQ	newm+0(FP), CX
	MOVQ	m_g0(CX), DX

	// Layout new m scheduler stack on os stack.
	MOVQ	SP, AX
	MOVQ	AX, (g_stack+stack_hi)(DX)
	SUBQ	$(64*1024), AX		// stack size
	MOVQ	AX, (g_stack+stack_lo)(DX)
	MOVQ	AX, g_stackguard0(DX)
	MOVQ	AX, g_stackguard1(DX)

	// Initialize procid from TOS struct.
	MOVQ	_tos(SB), AX
	MOVL	64(AX), AX
	MOVQ	AX, m_procid(CX)	// save pid as m->procid

	// Finally, initialize g.
	get_tls(BX)
	MOVQ	DX, g(BX)

	CALL	runtime·stackcheck(SB)	// smashes AX, CX
	CALL	runtime·mstart(SB)

	MOVQ	$0x1234, 0x1234		// not reached
	RET

// This is needed by asm_amd64.s
TEXT runtime·settls(SB),NOSPLIT,$0
	RET

// void sigtramp(void *ureg, int8 *note)
TEXT runtime·sigtramp(SB),NOSPLIT,$0
	get_tls(AX)
	CMPQ	AX, $0
	JEQ sigtramp_badtls1

	// check that g exists
	MOVQ	g(AX), BX
	CMPQ	BX, $0
	JEQ sigtramp_badsig1

	// change stack
	MOVQ	g_m(BX), BX
	MOVQ	m_gsignal(BX), R10
	MOVQ	(g_stack+stack_hi)(R10), BP
	CMPQ	BP, $0
	JEQ sigtramp_badstack1

	// check for 16-alignment of stack
	MOVQ	BP, DX
	ANDQ	$-0x10, DX
	CMPQ	DX, $0
	JEQ sigtramp_badstack1

	MOVQ	BP, SP

	PUSHQ	g(AX) // stash g
	MOVQ	R10, g(AX) // g = m->gsignal

	PUSHQ	AX // retval
	PUSHQ	BX // gp
	PUSHQ	SI // note
	PUSHQ	DI // ureg
	CALL	runtime·sighandler(SB)// func sighandler(_ureg *ureg, note *byte, gp *g) int
	POPQ	AX // ureg
	POPQ	AX // note
	POPQ	AX // gp
	POPQ	DI // retval -> DI (syscall first arg)

	get_tls(AX)
	POPQ	BX // stashed g
	MOVQ	BX, g(AX) // restore g
	MOVQ	$4125, AX	// syscall noted(int)
	SYSCALL
sigtramp_badtls1:
	CALL	runtime·badtls1(SB) // will exit
sigtramp_badsig1:
	CALL	runtime·badsignal1(SB) // will exit
sigtramp_badstack1:
	CALL	runtime·badstack1(SB) // will exit

TEXT runtime·setfpmasks(SB),NOSPLIT,$8
	STMXCSR	0(SP)
	MOVL	0(SP), AX
	ANDL	$~0x3F, AX
	ORL	$(0x3F<<7), AX
	MOVL	AX, 0(SP)
	LDMXCSR	0(SP)
	RET


