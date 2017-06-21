// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

// setldt(int entry, int address, int limit)
TEXT runtime·setldt(SB),NOSPLIT,$0
	RET

TEXT runtime·open(SB),NOSPLIT,$0-20
	MOVQ	name+0(FP), DI
	MOVL	mode+8(FP), SI
	MOVL	perm+12(FP), DX
	MOVL	$14, AX			// syscall entry
	SYSCALL
	MOVQ	AX, ret+16(FP)
	RET

TEXT runtime·pread(SB),NOSPLIT,$0
	MOVL	a0+0(FP), DI
	MOVQ	a1+8(FP), SI
	MOVL	a2+16(FP), DX
	MOVQ	a3+24(FP), R10
	MOVQ	$50, AX
	SYSCALL
	MOVQ	AX, ret+32(FP)
	RET

TEXT runtime·pwrite(SB),NOSPLIT,$0
	MOVL	fd+0(FP), DI
	MOVQ	p+8(FP), SI
	MOVL	n+16(FP), DX
	MOVQ	a4+24(FP), R10
	MOVQ	$51, AX
	SYSCALL
	MOVQ	AX, ret+32(FP)
	RET

// int64 seek(int32, int64, int32)
TEXT runtime·seek(SB),NOSPLIT,$32
	MOVL	fd+0(FP), DI
	MOVQ	off+4(FP), SI
	MOVL	whence+12(FP), DX
	MOVQ	$39, AX
	SYSCALL
	MOVQ	AX, ret+16(FP)
	RET

TEXT runtime·closefd(SB),NOSPLIT,$0
	MOVL	a0+0(FP), DI
	MOVQ	$4, AX
	SYSCALL
	MOVQ AX, ret+8(FP)
	RET

TEXT runtime·exits(SB),NOSPLIT,$0
	MOVQ	a0+0(FP), DI
	MOVQ	$8, AX
	SYSCALL
	RET

TEXT runtime·brk_(SB),NOSPLIT,$0
	MOVQ	a0+0(FP), DI
	MOVQ	$24, AX
	SYSCALL
	MOVQ AX, ret+8(FP)
	RET

TEXT runtime·sleep(SB),NOSPLIT,$0
	MOVL	a0+0(FP), DI
	MOVQ	$17, AX
	SYSCALL
	MOVQ AX, ret+8(FP)
	RET

TEXT runtime·plan9_semacquire(SB),NOSPLIT,$0
	MOVQ	a0+0(FP), DI
	MOVQ	a1+8(FP), SI
	MOVQ	$37, AX
	SYSCALL
	MOVQ AX, ret+16(FP)
	RET

TEXT runtime·plan9_tsemacquire(SB),NOSPLIT,$0
	MOVQ	a0+0(FP), DI
	MOVQ	a1+8(FP), SI
	MOVQ	$52, AX
	SYSCALL
	MOVQ AX, ret+16(FP)
	RET

TEXT runtime·nsec(SB),NOSPLIT,$0
	MOVQ	$53, AX
	SYSCALL
	MOVQ	AX, ret+8(FP)
	RET

// func walltime() (sec int64, nsec int32)
TEXT runtime·walltime(SB),NOSPLIT,$8-12
	CALL	runtime·nanotime(SB)
	MOVQ	0(SP), AX

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

TEXT runtime·notify(SB),NOSPLIT,$0
	MOVQ	a0+0(FP), DI
	MOVQ	$28, AX
	SYSCALL
	MOVQ AX, ret+8(FP)
	RET

TEXT runtime·noted(SB),NOSPLIT,$0
	MOVL	a0+0(FP), DI
	MOVQ	$29, AX
	SYSCALL
	MOVQ AX, ret+8(FP)
	RET
	
TEXT runtime·plan9_semrelease(SB),NOSPLIT,$0
	MOVQ	a0+0(FP), DI
	MOVQ	a1+8(FP), SI
	MOVQ	$38, AX
	SYSCALL
	MOVQ AX, ret+16(FP)
	RET

TEXT runtime·rfork(SB),NOSPLIT,$0
	MOVL	a0+0(FP), DI
	MOVQ	$19, AX
	SYSCALL
	MOVQ AX, ret+8(FP)
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
	MOVL	24(AX), AX
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

	// check that g exists
	MOVQ	g(AX), BX
	CMPQ	BX, $0
	JNE	3(PC)
	CALL	runtime·badsignal2(SB) // will exit
	RET

	// save args
	MOVQ	DI, CX
	MOVQ	SI, DX

	// change stack
	MOVQ	g_m(BX), BX
	MOVQ	m_gsignal(BX), R10
	MOVQ	(g_stack+stack_hi)(R10), BP
	MOVQ	BP, SP

	// make room for args and g
	SUBQ	$128, SP

	// save g
	MOVQ	g(AX), BP
	MOVQ	BP, 32(SP)

	// g = m->gsignal
	MOVQ	R10, g(AX)

	// load args and call sighandler
	MOVQ	CX, 0(SP)
	MOVQ	DX, 8(SP)
	MOVQ	BP, 16(SP)

	CALL	runtime·sighandler(SB)
	MOVL	24(SP), AX

	// restore g
	get_tls(BX)
	MOVQ	32(SP), R10
	MOVQ	R10, g(BX)

	// call noted(AX)
	MOVQ	AX, 0(SP)
	CALL	runtime·noted(SB)
	RET

TEXT runtime·setfpmasks(SB),NOSPLIT,$8
	STMXCSR	0(SP)
	MOVL	0(SP), AX
	ANDL	$~0x3F, AX
	ORL	$(0x3F<<7), AX
	MOVL	AX, 0(SP)
	LDMXCSR	0(SP)
	RET

#define ERRMAX 128	/* from os_plan9.h */

// void errstr(int8 *buf, int32 len)
TEXT errstr<>(SB),NOSPLIT,$0
	MOVQ	buf+0(FP), DI
	MOVQ	len+8(FP), SI
	MOVQ    $41, AX
	SYSCALL
	RET

// func errstr() string
// Only used by package syscall.
// Grab error string due to a syscall made
// in entersyscall mode, without going
// through the allocator (issue 4994).
// See ../syscall/asm_plan9_amd64.s:/·Syscall/
TEXT runtime·errstr(SB),NOSPLIT,$16-16
	get_tls(AX)
	MOVQ	g(AX), BX
	MOVQ	g_m(BX), BX
	MOVQ	(m_mOS+mOS_errstr)(BX), CX
	MOVQ	CX, 0(SP)
	MOVQ	$ERRMAX, 8(SP)
	CALL	errstr<>(SB)
	CALL	runtime·findnull(SB)
	MOVQ	8(SP), AX
	MOVQ	AX, ret_len+8(FP)
	MOVQ	0(SP), AX
	MOVQ	AX, ret_base+0(FP)
	RET
