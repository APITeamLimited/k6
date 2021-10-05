// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix darwin freebsd linux netbsd openbsd solaris zos

package unix

import (
	"runtime"
)

// Round the length of a raw sockaddr up to align it properly.
func cmsgAlignOf(salen int) int ***REMOVED***
	salign := SizeofPtr

	// dragonfly needs to check ABI version at runtime, see cmsgAlignOf in
	// sockcmsg_dragonfly.go
	switch runtime.GOOS ***REMOVED***
	case "aix":
		// There is no alignment on AIX.
		salign = 1
	case "darwin", "ios", "illumos", "solaris":
		// NOTE: It seems like 64-bit Darwin, Illumos and Solaris
		// kernels still require 32-bit aligned access to network
		// subsystem.
		if SizeofPtr == 8 ***REMOVED***
			salign = 4
		***REMOVED***
	case "netbsd", "openbsd":
		// NetBSD and OpenBSD armv7 require 64-bit alignment.
		if runtime.GOARCH == "arm" ***REMOVED***
			salign = 8
		***REMOVED***
		// NetBSD aarch64 requires 128-bit alignment.
		if runtime.GOOS == "netbsd" && runtime.GOARCH == "arm64" ***REMOVED***
			salign = 16
		***REMOVED***
	case "zos":
		// z/OS socket macros use [32-bit] sizeof(int) alignment,
		// not pointer width.
		salign = SizeofInt
	***REMOVED***

	return (salen + salign - 1) & ^(salign - 1)
***REMOVED***
