// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Plan 9 system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and
// wrap it in our own nicer implementation.

package syscall

import "unsafe"

const ImplementsGetwd = true

// ErrorString implements Error's String method by returning itself.
type ErrorString string

func (e ErrorString) Error() string { return string(e) }

// NewError converts s to an ErrorString, which satisfies the Error interface.
func NewError(s string) error { return ErrorString(s) }

func (e ErrorString) Temporary() bool {
	return e == EINTR || e == EMFILE || e.Timeout()
}

func (e ErrorString) Timeout() bool {
	return e == EBUSY || e == ETIMEDOUT
}

// A Note is a string describing a process note.
// It implements the os.Signal interface.
type Note string

func (n Note) Signal() {}

func (n Note) String() string {
	return string(n)
}

var (
	Stdin  = 0
	Stdout = 1
	Stderr = 2
)

// For testing: clients can set this flag to force
// creation of IPv6 sockets to return EAFNOSUPPORT.
var SocketDisableIPv6 bool

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err ErrorString)
func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err ErrorString)
func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr)
func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr)

func atoi(b []byte) (n uint) {
	n = 0
	for i := 0; i < len(b); i++ {
		n = n*10 + uint(b[i]-'0')
	}
	return
}

func cstring(s []byte) string {
	for i := range s {
		if s[i] == 0 {
			return string(s[0:i])
		}
	}
	return string(s)
}

func errstr() string {
	var buf [ERRMAX]byte

	RawSyscall(SYS_ERRSTR, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), 0)

	buf[len(buf)-1] = 0
	return cstring(buf[:])
}

// Implemented in assembly to import from runtime.
func exit(code int)

func Exit(code int) { exit(code) }

func readnum(path string) (uint, error) {
	var b [12]byte

	fd, e := Open(path, O_RDONLY)
	if e != nil {
		return 0, e
	}
	defer Close(fd)

	n, e := Pread(fd, b[:], 0)

	if e != nil {
		return 0, e
	}

	m := 0
	for ; m < n && b[m] == ' '; m++ {
	}

	return atoi(b[m : n-1]), nil
}

func Getpid() (pid int) {
	n, _ := readnum("#c/pid")
	return int(n)
}

func Getppid() (ppid int) {
	n, _ := readnum("#c/ppid")
	return int(n)
}

func Read(fd int, p []byte) (n int, err error) {
	return Pread(fd, p, -1)
}

func Write(fd int, p []byte) (n int, err error) {
	return Pwrite(fd, p, -1)
}

var ioSync int64

//sys	fd2path(fd int, buf []byte) (err error)
func Fd2path(fd int) (path string, err error) {
	var buf [512]byte

	e := fd2path(fd, buf[:])
	if e != nil {
		return "", e
	}
	return cstring(buf[:]), nil
}

//sys	pipe(p *[2]int32) (err error)
func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return NewError("bad arg in system call")
	}
	var pp [2]int32
	err = pipe(&pp)
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
}

// Underlying system call writes to newoffset via pointer.
// Assembly kept as pastey from Syscall() as possible.
func Seekcall(trap uintptr, fd int, offset int64, whence int) (r1 int64, r2 uintptr, err ErrorString)

func Seek(fd int, offset int64, whence int) (int64, error) {
	newoffset, _, err := Seekcall(SYS_SEEK, fd, offset, whence)
	if newoffset == -1 {
		return newoffset, err
	}
	return newoffset, nil
}

func Mkdir(path string, mode uint32) (err error) {
	fd, err := Create(path, O_RDONLY, DMDIR|mode)

	if fd != -1 {
		Close(fd)
	}

	return
}

type Waitmsg struct {
	Pid  int
	Time [3]uint32
	Msg  string
}

func (w Waitmsg) Exited() bool   { return true }
func (w Waitmsg) Signaled() bool { return false }

func (w Waitmsg) ExitStatus() int {
	if len(w.Msg) == 0 {
		// a normal exit returns no message
		return 0
	}
	return 1
}

//sys	await(s []byte) (n int, err error)
func Await(w *Waitmsg) (err error) {
	var buf [512]byte
	var f [5][]byte

	n, err := await(buf[:])

	if err != nil || w == nil {
		return
	}

	nf := 0
	p := 0
	for i := 0; i < n && nf < len(f)-1; i++ {
		if buf[i] == ' ' {
			f[nf] = buf[p:i]
			p = i + 1
			nf++
		}
	}
	f[nf] = buf[p:]
	nf++

	if nf != len(f) {
		return NewError("invalid wait message")
	}
	w.Pid = int(atoi(f[0]))
	w.Time[0] = uint32(atoi(f[1]))
	w.Time[1] = uint32(atoi(f[2]))
	w.Time[2] = uint32(atoi(f[3]))
	w.Msg = cstring(f[4])
	if w.Msg == "''" {
		// await() returns '' for no error
		w.Msg = ""
	}
	return
}

func Unmount(name, old string) (err error) {
	Fixwd()
	oldp, err := BytePtrFromString(old)
	if err != nil {
		return err
	}
	oldptr := uintptr(unsafe.Pointer(oldp))

	var r0 uintptr
	var e ErrorString

	// bind(2) man page: If name is zero, everything bound or mounted upon old is unbound or unmounted.
	if name == "" {
		r0, _, e = Syscall(SYS_UNMOUNT, _zero, oldptr, 0)
	} else {
		namep, err := BytePtrFromString(name)
		if err != nil {
			return err
		}
		r0, _, e = Syscall(SYS_UNMOUNT, uintptr(unsafe.Pointer(namep)), oldptr, 0)
		use(unsafe.Pointer(namep))
	}
	use(unsafe.Pointer(oldp))

	if int32(r0) == -1 {
		err = e
	}
	return
}

func Fchdir(fd int) (err error) {
	path, err := Fd2path(fd)

	if err != nil {
		return
	}

	return Chdir(path)
}

type Timespec struct {
	Sec  int32
	Nsec int32
}

type Timeval struct {
	Sec  int32
	Usec int32
}

func NsecToTimeval(nsec int64) (tv Timeval) {
	nsec += 999 // round up to microsecond
	tv.Usec = int32(nsec % 1e9 / 1e3)
	tv.Sec = int32(nsec / 1e9)
	return
}

func nsec() int64 {
	var scratch int64

	r0, _, _ := Syscall(SYS_NSEC, uintptr(unsafe.Pointer(&scratch)), 0, 0)
	// TODO(aram): remove hack after I fix _nsec in the pc64 kernel.
	if r0 == 0 {
		return scratch
	}
	return int64(r0)
}

func Gettimeofday(tv *Timeval) error {
	nsec := nsec()
	*tv = NsecToTimeval(nsec)
	return nil
}

func Getpagesize() int { return 0x1000 }

func Getegid() (egid int) { return -1 }
func Geteuid() (euid int) { return -1 }
func Getgid() (gid int)   { return -1 }
func Getuid() (uid int)   { return -1 }

func Getgroups() (gids []int, err error) {
	return make([]int, 0), nil
}

func Open(path string, mode int) (fd int, err error) {
	Fixwd()
	return open(path, mode)
}

func Create(path string, mode int, perm uint32) (fd int, err error) {
	Fixwd()
	return create(path, mode, perm)
}

func Remove(path string) error {
	Fixwd()
	return remove(path)
}

func Stat(path string, edir []byte) (n int, err error) {
	Fixwd()
	return stat(path, edir)
}

func Bind(name string, old string, flag int) (err error) {
	Fixwd()
	return bind(name, old, flag)
}

func Mount(fd int, afd int, old string, flag int, aname string) (err error) {
	Fixwd()
	return mount(fd, afd, old, flag, aname)
}

func Wstat(path string, edir []byte) (err error) {
	Fixwd()
	return wstat(path, edir)
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func fd2path(fd int, buf []byte) (err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_FD2PATH, uintptr(fd), uintptr(_p0), uintptr(len(buf)))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func pipe(p *[2]int32) (err error) {
	r0, _, e1 := Syscall(SYS_PIPE, uintptr(unsafe.Pointer(p)), 0, 0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func await(s []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(s) > 0 {
		_p0 = unsafe.Pointer(&s[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_AWAIT, uintptr(_p0), uintptr(len(s)), 0)
	n = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func open(path string, mode int) (fd int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_OPEN, uintptr(unsafe.Pointer(_p0)), uintptr(mode), 0)
	use(unsafe.Pointer(_p0))
	fd = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func create(path string, mode int, perm uint32) (fd int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_CREATE, uintptr(unsafe.Pointer(_p0)), uintptr(mode), uintptr(perm))
	use(unsafe.Pointer(_p0))
	fd = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func remove(path string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_REMOVE, uintptr(unsafe.Pointer(_p0)), 0, 0)
	use(unsafe.Pointer(_p0))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func stat(path string, edir []byte) (n int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(edir) > 0 {
		_p1 = unsafe.Pointer(&edir[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_STAT, uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(edir)))
	use(unsafe.Pointer(_p0))
	n = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func bind(name string, old string, flag int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(name)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(old)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_BIND, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(flag))
	use(unsafe.Pointer(_p0))
	use(unsafe.Pointer(_p1))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func mount(fd int, afd int, old string, flag int, aname string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(old)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(aname)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall6(SYS_MOUNT, uintptr(fd), uintptr(afd), uintptr(unsafe.Pointer(_p0)), uintptr(flag), uintptr(unsafe.Pointer(_p1)), 0)
	use(unsafe.Pointer(_p0))
	use(unsafe.Pointer(_p1))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func wstat(path string, edir []byte) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(edir) > 0 {
		_p1 = unsafe.Pointer(&edir[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_WSTAT, uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(edir)))
	use(unsafe.Pointer(_p0))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func chdir(path string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_CHDIR, uintptr(unsafe.Pointer(_p0)), 0, 0)
	use(unsafe.Pointer(_p0))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Dup(oldfd int, newfd int) (fd int, err error) {
	r0, _, e1 := Syscall(SYS_DUP, uintptr(oldfd), uintptr(newfd), 0)
	fd = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Pread(fd int, p []byte, offset int64) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_PREAD, uintptr(fd), uintptr(_p0), uintptr(len(p)), uintptr(offset), uintptr(offset>>32), 0)
	n = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Pwrite(fd int, p []byte, offset int64) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_PWRITE, uintptr(fd), uintptr(_p0), uintptr(len(p)), uintptr(offset), uintptr(offset>>32), 0)
	n = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Close(fd int) (err error) {
	r0, _, e1 := Syscall(SYS_CLOSE, uintptr(fd), 0, 0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fstat(fd int, edir []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(edir) > 0 {
		_p0 = unsafe.Pointer(&edir[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_FSTAT, uintptr(fd), uintptr(_p0), uintptr(len(edir)))
	n = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fwstat(fd int, edir []byte) (err error) {
	var _p0 unsafe.Pointer
	if len(edir) > 0 {
		_p0 = unsafe.Pointer(&edir[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_FWSTAT, uintptr(fd), uintptr(_p0), uintptr(len(edir)))
	if int32(r0) == -1 {
		err = e1
	}
	return
}
