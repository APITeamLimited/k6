// cgo -godefs types_netbsd.go | go run mkpost.go
// Code generated by the command above; see README.md. DO NOT EDIT.

//go:build amd64 && netbsd
// +build amd64,netbsd

package unix

const (
	SizeofPtr      = 0x8
	SizeofShort    = 0x2
	SizeofInt      = 0x4
	SizeofLong     = 0x8
	SizeofLongLong = 0x8
)

type (
	_C_short     int16
	_C_int       int32
	_C_long      int64
	_C_long_long int64
)

type Timespec struct ***REMOVED***
	Sec  int64
	Nsec int64
***REMOVED***

type Timeval struct ***REMOVED***
	Sec       int64
	Usec      int32
	Pad_cgo_0 [4]byte
***REMOVED***

type Rusage struct ***REMOVED***
	Utime    Timeval
	Stime    Timeval
	Maxrss   int64
	Ixrss    int64
	Idrss    int64
	Isrss    int64
	Minflt   int64
	Majflt   int64
	Nswap    int64
	Inblock  int64
	Oublock  int64
	Msgsnd   int64
	Msgrcv   int64
	Nsignals int64
	Nvcsw    int64
	Nivcsw   int64
***REMOVED***

type Rlimit struct ***REMOVED***
	Cur uint64
	Max uint64
***REMOVED***

type _Gid_t uint32

type Stat_t struct ***REMOVED***
	Dev     uint64
	Mode    uint32
	_       [4]byte
	Ino     uint64
	Nlink   uint32
	Uid     uint32
	Gid     uint32
	_       [4]byte
	Rdev    uint64
	Atim    Timespec
	Mtim    Timespec
	Ctim    Timespec
	Btim    Timespec
	Size    int64
	Blocks  int64
	Blksize uint32
	Flags   uint32
	Gen     uint32
	Spare   [2]uint32
	_       [4]byte
***REMOVED***

type Statfs_t [0]byte

type Statvfs_t struct ***REMOVED***
	Flag        uint64
	Bsize       uint64
	Frsize      uint64
	Iosize      uint64
	Blocks      uint64
	Bfree       uint64
	Bavail      uint64
	Bresvd      uint64
	Files       uint64
	Ffree       uint64
	Favail      uint64
	Fresvd      uint64
	Syncreads   uint64
	Syncwrites  uint64
	Asyncreads  uint64
	Asyncwrites uint64
	Fsidx       Fsid
	Fsid        uint64
	Namemax     uint64
	Owner       uint32
	Spare       [4]uint32
	Fstypename  [32]byte
	Mntonname   [1024]byte
	Mntfromname [1024]byte
	_           [4]byte
***REMOVED***

type Flock_t struct ***REMOVED***
	Start  int64
	Len    int64
	Pid    int32
	Type   int16
	Whence int16
***REMOVED***

type Dirent struct ***REMOVED***
	Fileno    uint64
	Reclen    uint16
	Namlen    uint16
	Type      uint8
	Name      [512]int8
	Pad_cgo_0 [3]byte
***REMOVED***

type Fsid struct ***REMOVED***
	X__fsid_val [2]int32
***REMOVED***

const (
	PathMax = 0x400
)

const (
	ST_WAIT   = 0x1
	ST_NOWAIT = 0x2
)

const (
	FADV_NORMAL     = 0x0
	FADV_RANDOM     = 0x1
	FADV_SEQUENTIAL = 0x2
	FADV_WILLNEED   = 0x3
	FADV_DONTNEED   = 0x4
	FADV_NOREUSE    = 0x5
)

type RawSockaddrInet4 struct ***REMOVED***
	Len    uint8
	Family uint8
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]int8
***REMOVED***

type RawSockaddrInet6 struct ***REMOVED***
	Len      uint8
	Family   uint8
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
***REMOVED***

type RawSockaddrUnix struct ***REMOVED***
	Len    uint8
	Family uint8
	Path   [104]int8
***REMOVED***

type RawSockaddrDatalink struct ***REMOVED***
	Len    uint8
	Family uint8
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [12]int8
***REMOVED***

type RawSockaddr struct ***REMOVED***
	Len    uint8
	Family uint8
	Data   [14]int8
***REMOVED***

type RawSockaddrAny struct ***REMOVED***
	Addr RawSockaddr
	Pad  [92]int8
***REMOVED***

type _Socklen uint32

type Linger struct ***REMOVED***
	Onoff  int32
	Linger int32
***REMOVED***

type Iovec struct ***REMOVED***
	Base *byte
	Len  uint64
***REMOVED***

type IPMreq struct ***REMOVED***
	Multiaddr [4]byte /* in_addr */
	Interface [4]byte /* in_addr */
***REMOVED***

type IPv6Mreq struct ***REMOVED***
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
***REMOVED***

type Msghdr struct ***REMOVED***
	Name       *byte
	Namelen    uint32
	Pad_cgo_0  [4]byte
	Iov        *Iovec
	Iovlen     int32
	Pad_cgo_1  [4]byte
	Control    *byte
	Controllen uint32
	Flags      int32
***REMOVED***

type Cmsghdr struct ***REMOVED***
	Len   uint32
	Level int32
	Type  int32
***REMOVED***

type Inet6Pktinfo struct ***REMOVED***
	Addr    [16]byte /* in6_addr */
	Ifindex uint32
***REMOVED***

type IPv6MTUInfo struct ***REMOVED***
	Addr RawSockaddrInet6
	Mtu  uint32
***REMOVED***

type ICMPv6Filter struct ***REMOVED***
	Filt [8]uint32
***REMOVED***

const (
	SizeofSockaddrInet4    = 0x10
	SizeofSockaddrInet6    = 0x1c
	SizeofSockaddrAny      = 0x6c
	SizeofSockaddrUnix     = 0x6a
	SizeofSockaddrDatalink = 0x14
	SizeofLinger           = 0x8
	SizeofIovec            = 0x10
	SizeofIPMreq           = 0x8
	SizeofIPv6Mreq         = 0x14
	SizeofMsghdr           = 0x30
	SizeofCmsghdr          = 0xc
	SizeofInet6Pktinfo     = 0x14
	SizeofIPv6MTUInfo      = 0x20
	SizeofICMPv6Filter     = 0x20
)

const (
	PTRACE_TRACEME = 0x0
	PTRACE_CONT    = 0x7
	PTRACE_KILL    = 0x8
)

type Kevent_t struct ***REMOVED***
	Ident     uint64
	Filter    uint32
	Flags     uint32
	Fflags    uint32
	Pad_cgo_0 [4]byte
	Data      int64
	Udata     int64
***REMOVED***

type FdSet struct ***REMOVED***
	Bits [8]uint32
***REMOVED***

const (
	SizeofIfMsghdr         = 0x98
	SizeofIfData           = 0x88
	SizeofIfaMsghdr        = 0x18
	SizeofIfAnnounceMsghdr = 0x18
	SizeofRtMsghdr         = 0x78
	SizeofRtMetrics        = 0x50
)

type IfMsghdr struct ***REMOVED***
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	Pad_cgo_0 [2]byte
	Data      IfData
***REMOVED***

type IfData struct ***REMOVED***
	Type       uint8
	Addrlen    uint8
	Hdrlen     uint8
	Pad_cgo_0  [1]byte
	Link_state int32
	Mtu        uint64
	Metric     uint64
	Baudrate   uint64
	Ipackets   uint64
	Ierrors    uint64
	Opackets   uint64
	Oerrors    uint64
	Collisions uint64
	Ibytes     uint64
	Obytes     uint64
	Imcasts    uint64
	Omcasts    uint64
	Iqdrops    uint64
	Noproto    uint64
	Lastchange Timespec
***REMOVED***

type IfaMsghdr struct ***REMOVED***
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Metric    int32
	Index     uint16
	Pad_cgo_0 [6]byte
***REMOVED***

type IfAnnounceMsghdr struct ***REMOVED***
	Msglen  uint16
	Version uint8
	Type    uint8
	Index   uint16
	Name    [16]int8
	What    uint16
***REMOVED***

type RtMsghdr struct ***REMOVED***
	Msglen    uint16
	Version   uint8
	Type      uint8
	Index     uint16
	Pad_cgo_0 [2]byte
	Flags     int32
	Addrs     int32
	Pid       int32
	Seq       int32
	Errno     int32
	Use       int32
	Inits     int32
	Pad_cgo_1 [4]byte
	Rmx       RtMetrics
***REMOVED***

type RtMetrics struct ***REMOVED***
	Locks    uint64
	Mtu      uint64
	Hopcount uint64
	Recvpipe uint64
	Sendpipe uint64
	Ssthresh uint64
	Rtt      uint64
	Rttvar   uint64
	Expire   int64
	Pksent   int64
***REMOVED***

type Mclpool [0]byte

const (
	SizeofBpfVersion = 0x4
	SizeofBpfStat    = 0x80
	SizeofBpfProgram = 0x10
	SizeofBpfInsn    = 0x8
	SizeofBpfHdr     = 0x20
)

type BpfVersion struct ***REMOVED***
	Major uint16
	Minor uint16
***REMOVED***

type BpfStat struct ***REMOVED***
	Recv    uint64
	Drop    uint64
	Capt    uint64
	Padding [13]uint64
***REMOVED***

type BpfProgram struct ***REMOVED***
	Len       uint32
	Pad_cgo_0 [4]byte
	Insns     *BpfInsn
***REMOVED***

type BpfInsn struct ***REMOVED***
	Code uint16
	Jt   uint8
	Jf   uint8
	K    uint32
***REMOVED***

type BpfHdr struct ***REMOVED***
	Tstamp    BpfTimeval
	Caplen    uint32
	Datalen   uint32
	Hdrlen    uint16
	Pad_cgo_0 [6]byte
***REMOVED***

type BpfTimeval struct ***REMOVED***
	Sec  int64
	Usec int64
***REMOVED***

type Termios struct ***REMOVED***
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]uint8
	Ispeed int32
	Ospeed int32
***REMOVED***

type Winsize struct ***REMOVED***
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
***REMOVED***

type Ptmget struct ***REMOVED***
	Cfd int32
	Sfd int32
	Cn  [1024]byte
	Sn  [1024]byte
***REMOVED***

const (
	AT_FDCWD            = -0x64
	AT_SYMLINK_FOLLOW   = 0x400
	AT_SYMLINK_NOFOLLOW = 0x200
)

type PollFd struct ***REMOVED***
	Fd      int32
	Events  int16
	Revents int16
***REMOVED***

const (
	POLLERR    = 0x8
	POLLHUP    = 0x10
	POLLIN     = 0x1
	POLLNVAL   = 0x20
	POLLOUT    = 0x4
	POLLPRI    = 0x2
	POLLRDBAND = 0x80
	POLLRDNORM = 0x40
	POLLWRBAND = 0x100
	POLLWRNORM = 0x4
)

type Sysctlnode struct ***REMOVED***
	Flags           uint32
	Num             int32
	Name            [32]int8
	Ver             uint32
	X__rsvd         uint32
	Un              [16]byte
	X_sysctl_size   [8]byte
	X_sysctl_func   [8]byte
	X_sysctl_parent [8]byte
	X_sysctl_desc   [8]byte
***REMOVED***

type Utsname struct ***REMOVED***
	Sysname  [256]byte
	Nodename [256]byte
	Release  [256]byte
	Version  [256]byte
	Machine  [256]byte
***REMOVED***

const SizeofClockinfo = 0x14

type Clockinfo struct ***REMOVED***
	Hz      int32
	Tick    int32
	Tickadj int32
	Stathz  int32
	Profhz  int32
***REMOVED***
