// go run mksyscall.go -tags darwin,arm64,go1.13 syscall_darwin.1_13.go
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build darwin,arm64,go1.13

package unix

import (
	"syscall"
	"unsafe"
)

var _ syscall.Errno

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func closedir(dir uintptr) (err error) ***REMOVED***
	_, _, e1 := syscall_syscall(funcPC(libc_closedir_trampoline), uintptr(dir), 0, 0)
	if e1 != 0 ***REMOVED***
		err = errnoErr(e1)
	***REMOVED***
	return
***REMOVED***

func libc_closedir_trampoline()

//go:cgo_import_dynamic libc_closedir closedir "/usr/lib/libSystem.B.dylib"

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func readdir_r(dir uintptr, entry *Dirent, result **Dirent) (res Errno) ***REMOVED***
	r0, _, _ := syscall_syscall(funcPC(libc_readdir_r_trampoline), uintptr(dir), uintptr(unsafe.Pointer(entry)), uintptr(unsafe.Pointer(result)))
	res = Errno(r0)
	return
***REMOVED***

func libc_readdir_r_trampoline()

//go:cgo_import_dynamic libc_readdir_r readdir_r "/usr/lib/libSystem.B.dylib"
