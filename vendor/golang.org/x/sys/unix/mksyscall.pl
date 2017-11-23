#!/usr/bin/env perl
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This program reads a file containing function prototypes
# (like syscall_darwin.go) and generates system call bodies.
# The prototypes are marked by lines beginning with "//sys"
# and read like func declarations if //sys is replaced by func, but:
#	* The parameter lists must give a name for each argument.
#	  This includes return parameters.
#	* The parameter lists must give a type for each argument:
#	  the (x, y, z int) shorthand is not allowed.
#	* If the return parameter is an error number, it must be named errno.

# A line beginning with //sysnb is like //sys, except that the
# goroutine will not be suspended during the execution of the system
# call.  This must only be used for system calls which can never
# block, as otherwise the system call could cause all goroutines to
# hang.

use strict;

my $cmdline = "mksyscall.pl " . join(' ', @ARGV);
my $errors = 0;
my $_32bit = "";
my $plan9 = 0;
my $openbsd = 0;
my $netbsd = 0;
my $dragonfly = 0;
my $arm = 0; # 64-bit value should use (even, odd)-pair
my $tags = "";  # build tags

if($ARGV[0] eq "-b32") ***REMOVED***
	$_32bit = "big-endian";
	shift;
***REMOVED*** elsif($ARGV[0] eq "-l32") ***REMOVED***
	$_32bit = "little-endian";
	shift;
***REMOVED***
if($ARGV[0] eq "-plan9") ***REMOVED***
	$plan9 = 1;
	shift;
***REMOVED***
if($ARGV[0] eq "-openbsd") ***REMOVED***
	$openbsd = 1;
	shift;
***REMOVED***
if($ARGV[0] eq "-netbsd") ***REMOVED***
	$netbsd = 1;
	shift;
***REMOVED***
if($ARGV[0] eq "-dragonfly") ***REMOVED***
	$dragonfly = 1;
	shift;
***REMOVED***
if($ARGV[0] eq "-arm") ***REMOVED***
	$arm = 1;
	shift;
***REMOVED***
if($ARGV[0] eq "-tags") ***REMOVED***
	shift;
	$tags = $ARGV[0];
	shift;
***REMOVED***

if($ARGV[0] =~ /^-/) ***REMOVED***
	print STDERR "usage: mksyscall.pl [-b32 | -l32] [-tags x,y] [file ...]\n";
	exit 1;
***REMOVED***

# Check that we are using the new build system if we should
if($ENV***REMOVED***'GOOS'***REMOVED*** eq "linux" && $ENV***REMOVED***'GOARCH'***REMOVED*** ne "sparc64") ***REMOVED***
	if($ENV***REMOVED***'GOLANG_SYS_BUILD'***REMOVED*** ne "docker") ***REMOVED***
		print STDERR "In the new build system, mksyscall should not be called directly.\n";
		print STDERR "See README.md\n";
		exit 1;
	***REMOVED***
***REMOVED***


sub parseparamlist($) ***REMOVED***
	my ($list) = @_;
	$list =~ s/^\s*//;
	$list =~ s/\s*$//;
	if($list eq "") ***REMOVED***
		return ();
	***REMOVED***
	return split(/\s*,\s*/, $list);
***REMOVED***

sub parseparam($) ***REMOVED***
	my ($p) = @_;
	if($p !~ /^(\S*) (\S*)$/) ***REMOVED***
		print STDERR "$ARGV:$.: malformed parameter: $p\n";
		$errors = 1;
		return ("xx", "int");
	***REMOVED***
	return ($1, $2);
***REMOVED***

my $text = "";
while(<>) ***REMOVED***
	chomp;
	s/\s+/ /g;
	s/^\s+//;
	s/\s+$//;
	my $nonblock = /^\/\/sysnb /;
	next if !/^\/\/sys / && !$nonblock;

	# Line must be of the form
	#	func Open(path string, mode int, perm int) (fd int, errno error)
	# Split into name, in params, out params.
	if(!/^\/\/sys(nb)? (\w+)\(([^()]*)\)\s*(?:\(([^()]+)\))?\s*(?:=\s*((?i)SYS_[A-Z0-9_]+))?$/) ***REMOVED***
		print STDERR "$ARGV:$.: malformed //sys declaration\n";
		$errors = 1;
		next;
	***REMOVED***
	my ($func, $in, $out, $sysname) = ($2, $3, $4, $5);

	# Split argument lists on comma.
	my @in = parseparamlist($in);
	my @out = parseparamlist($out);

	# Try in vain to keep people from editing this file.
	# The theory is that they jump into the middle of the file
	# without reading the header.
	$text .= "// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT\n\n";

	# Go function header.
	my $out_decl = @out ? sprintf(" (%s)", join(', ', @out)) : "";
	$text .= sprintf "func %s(%s)%s ***REMOVED***\n", $func, join(', ', @in), $out_decl;

	# Check if err return available
	my $errvar = "";
	foreach my $p (@out) ***REMOVED***
		my ($name, $type) = parseparam($p);
		if($type eq "error") ***REMOVED***
			$errvar = $name;
			last;
		***REMOVED***
	***REMOVED***

	# Prepare arguments to Syscall.
	my @args = ();
	my $n = 0;
	foreach my $p (@in) ***REMOVED***
		my ($name, $type) = parseparam($p);
		if($type =~ /^\*/) ***REMOVED***
			push @args, "uintptr(unsafe.Pointer($name))";
		***REMOVED*** elsif($type eq "string" && $errvar ne "") ***REMOVED***
			$text .= "\tvar _p$n *byte\n";
			$text .= "\t_p$n, $errvar = BytePtrFromString($name)\n";
			$text .= "\tif $errvar != nil ***REMOVED***\n\t\treturn\n\t***REMOVED***\n";
			push @args, "uintptr(unsafe.Pointer(_p$n))";
			$n++;
		***REMOVED*** elsif($type eq "string") ***REMOVED***
			print STDERR "$ARGV:$.: $func uses string arguments, but has no error return\n";
			$text .= "\tvar _p$n *byte\n";
			$text .= "\t_p$n, _ = BytePtrFromString($name)\n";
			push @args, "uintptr(unsafe.Pointer(_p$n))";
			$n++;
		***REMOVED*** elsif($type =~ /^\[\](.*)/) ***REMOVED***
			# Convert slice into pointer, length.
			# Have to be careful not to take address of &a[0] if len == 0:
			# pass dummy pointer in that case.
			# Used to pass nil, but some OSes or simulators reject write(fd, nil, 0).
			$text .= "\tvar _p$n unsafe.Pointer\n";
			$text .= "\tif len($name) > 0 ***REMOVED***\n\t\t_p$n = unsafe.Pointer(\&$***REMOVED***name***REMOVED***[0])\n\t***REMOVED***";
			$text .= " else ***REMOVED***\n\t\t_p$n = unsafe.Pointer(&_zero)\n\t***REMOVED***";
			$text .= "\n";
			push @args, "uintptr(_p$n)", "uintptr(len($name))";
			$n++;
		***REMOVED*** elsif($type eq "int64" && ($openbsd || $netbsd)) ***REMOVED***
			push @args, "0";
			if($_32bit eq "big-endian") ***REMOVED***
				push @args, "uintptr($name>>32)", "uintptr($name)";
			***REMOVED*** elsif($_32bit eq "little-endian") ***REMOVED***
				push @args, "uintptr($name)", "uintptr($name>>32)";
			***REMOVED*** else ***REMOVED***
				push @args, "uintptr($name)";
			***REMOVED***
		***REMOVED*** elsif($type eq "int64" && $dragonfly) ***REMOVED***
			if ($func !~ /^extp(read|write)/i) ***REMOVED***
				push @args, "0";
			***REMOVED***
			if($_32bit eq "big-endian") ***REMOVED***
				push @args, "uintptr($name>>32)", "uintptr($name)";
			***REMOVED*** elsif($_32bit eq "little-endian") ***REMOVED***
				push @args, "uintptr($name)", "uintptr($name>>32)";
			***REMOVED*** else ***REMOVED***
				push @args, "uintptr($name)";
			***REMOVED***
		***REMOVED*** elsif($type eq "int64" && $_32bit ne "") ***REMOVED***
			if(@args % 2 && $arm) ***REMOVED***
				# arm abi specifies 64-bit argument uses
				# (even, odd) pair
				push @args, "0"
			***REMOVED***
			if($_32bit eq "big-endian") ***REMOVED***
				push @args, "uintptr($name>>32)", "uintptr($name)";
			***REMOVED*** else ***REMOVED***
				push @args, "uintptr($name)", "uintptr($name>>32)";
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			push @args, "uintptr($name)";
		***REMOVED***
	***REMOVED***

	# Determine which form to use; pad args with zeros.
	my $asm = "Syscall";
	if ($nonblock) ***REMOVED***
		$asm = "RawSyscall";
	***REMOVED***
	if(@args <= 3) ***REMOVED***
		while(@args < 3) ***REMOVED***
			push @args, "0";
		***REMOVED***
	***REMOVED*** elsif(@args <= 6) ***REMOVED***
		$asm .= "6";
		while(@args < 6) ***REMOVED***
			push @args, "0";
		***REMOVED***
	***REMOVED*** elsif(@args <= 9) ***REMOVED***
		$asm .= "9";
		while(@args < 9) ***REMOVED***
			push @args, "0";
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		print STDERR "$ARGV:$.: too many arguments to system call\n";
	***REMOVED***

	# System call number.
	if($sysname eq "") ***REMOVED***
		$sysname = "SYS_$func";
		$sysname =~ s/([a-z])([A-Z])/$***REMOVED***1***REMOVED***_$2/g;	# turn FooBar into Foo_Bar
		$sysname =~ y/a-z/A-Z/;
	***REMOVED***

	# Actual call.
	my $args = join(', ', @args);
	my $call = "$asm($sysname, $args)";

	# Assign return values.
	my $body = "";
	my @ret = ("_", "_", "_");
	my $do_errno = 0;
	for(my $i=0; $i<@out; $i++) ***REMOVED***
		my $p = $out[$i];
		my ($name, $type) = parseparam($p);
		my $reg = "";
		if($name eq "err" && !$plan9) ***REMOVED***
			$reg = "e1";
			$ret[2] = $reg;
			$do_errno = 1;
		***REMOVED*** elsif($name eq "err" && $plan9) ***REMOVED***
			$ret[0] = "r0";
			$ret[2] = "e1";
			next;
		***REMOVED*** else ***REMOVED***
			$reg = sprintf("r%d", $i);
			$ret[$i] = $reg;
		***REMOVED***
		if($type eq "bool") ***REMOVED***
			$reg = "$reg != 0";
		***REMOVED***
		if($type eq "int64" && $_32bit ne "") ***REMOVED***
			# 64-bit number in r1:r0 or r0:r1.
			if($i+2 > @out) ***REMOVED***
				print STDERR "$ARGV:$.: not enough registers for int64 return\n";
			***REMOVED***
			if($_32bit eq "big-endian") ***REMOVED***
				$reg = sprintf("int64(r%d)<<32 | int64(r%d)", $i, $i+1);
			***REMOVED*** else ***REMOVED***
				$reg = sprintf("int64(r%d)<<32 | int64(r%d)", $i+1, $i);
			***REMOVED***
			$ret[$i] = sprintf("r%d", $i);
			$ret[$i+1] = sprintf("r%d", $i+1);
		***REMOVED***
		if($reg ne "e1" || $plan9) ***REMOVED***
			$body .= "\t$name = $type($reg)\n";
		***REMOVED***
	***REMOVED***
	if ($ret[0] eq "_" && $ret[1] eq "_" && $ret[2] eq "_") ***REMOVED***
		$text .= "\t$call\n";
	***REMOVED*** else ***REMOVED***
		$text .= "\t$ret[0], $ret[1], $ret[2] := $call\n";
	***REMOVED***
	$text .= $body;

	if ($plan9 && $ret[2] eq "e1") ***REMOVED***
		$text .= "\tif int32(r0) == -1 ***REMOVED***\n";
		$text .= "\t\terr = e1\n";
		$text .= "\t***REMOVED***\n";
	***REMOVED*** elsif ($do_errno) ***REMOVED***
		$text .= "\tif e1 != 0 ***REMOVED***\n";
		$text .= "\t\terr = errnoErr(e1)\n";
		$text .= "\t***REMOVED***\n";
	***REMOVED***
	$text .= "\treturn\n";
	$text .= "***REMOVED***\n\n";
***REMOVED***

chomp $text;
chomp $text;

if($errors) ***REMOVED***
	exit 1;
***REMOVED***

print <<EOF;
// $cmdline
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build $tags

package unix

import (
	"syscall"
	"unsafe"
)

var _ syscall.Errno

$text
EOF
exit 0;
