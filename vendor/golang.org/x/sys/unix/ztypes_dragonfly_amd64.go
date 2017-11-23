// cgo -godefs types_dragonfly.go | go run mkpost.go
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build amd64,dragonfly

package unix

const (
	sizeofPtr      = 0x8
	sizeofShort    = 0x2
	sizeofInt      = 0x4
	sizeofLong     = 0x8
	sizeofLongLong = 0x8
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
	Sec  int64
	Usec int64
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
	Cur int64
	Max int64
***REMOVED***

type _Gid_t uint32

const (
	S_IFMT   = 0xf000
	S_IFIFO  = 0x1000
	S_IFCHR  = 0x2000
	S_IFDIR  = 0x4000
	S_IFBLK  = 0x6000
	S_IFREG  = 0x8000
	S_IFLNK  = 0xa000
	S_IFSOCK = 0xc000
	S_ISUID  = 0x800
	S_ISGID  = 0x400
	S_ISVTX  = 0x200
	S_IRUSR  = 0x100
	S_IWUSR  = 0x80
	S_IXUSR  = 0x40
)

type Stat_t struct ***REMOVED***
	Ino      uint64
	Nlink    uint32
	Dev      uint32
	Mode     uint16
	Padding1 uint16
	Uid      uint32
	Gid      uint32
	Rdev     uint32
	Atim     Timespec
	Mtim     Timespec
	Ctim     Timespec
	Size     int64
	Blocks   int64
	Blksize  uint32
	Flags    uint32
	Gen      uint32
	Lspare   int32
	Qspare1  int64
	Qspare2  int64
***REMOVED***

type Statfs_t struct ***REMOVED***
	Spare2      int64
	Bsize       int64
	Iosize      int64
	Blocks      int64
	Bfree       int64
	Bavail      int64
	Files       int64
	Ffree       int64
	Fsid        Fsid
	Owner       uint32
	Type        int32
	Flags       int32
	Pad_cgo_0   [4]byte
	Syncwrites  int64
	Asyncwrites int64
	Fstypename  [16]int8
	Mntonname   [80]int8
	Syncreads   int64
	Asyncreads  int64
	Spares1     int16
	Mntfromname [80]int8
	Spares2     int16
	Pad_cgo_1   [4]byte
	Spare       [2]int64
***REMOVED***

type Flock_t struct ***REMOVED***
	Start  int64
	Len    int64
	Pid    int32
	Type   int16
	Whence int16
***REMOVED***

type Dirent struct ***REMOVED***
	Fileno  uint64
	Namlen  uint16
	Type    uint8
	Unused1 uint8
	Unused2 uint32
	Name    [256]int8
***REMOVED***

type Fsid struct ***REMOVED***
	Val [2]int32
***REMOVED***

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
	Rcf    uint16
	Route  [16]uint16
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
	SizeofSockaddrDatalink = 0x36
	SizeofLinger           = 0x8
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
	Ident  uint64
	Filter int16
	Flags  uint16
	Fflags uint32
	Data   int64
	Udata  *byte
***REMOVED***

type FdSet struct ***REMOVED***
	Bits [16]uint64
***REMOVED***

const (
	SizeofIfMsghdr         = 0xb0
	SizeofIfData           = 0xa0
	SizeofIfaMsghdr        = 0x14
	SizeofIfmaMsghdr       = 0x10
	SizeofIfAnnounceMsghdr = 0x18
	SizeofRtMsghdr         = 0x98
	SizeofRtMetrics        = 0x70
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
	Physical   uint8
	Addrlen    uint8
	Hdrlen     uint8
	Recvquota  uint8
	Xmitquota  uint8
	Pad_cgo_0  [2]byte
	Mtu        uint64
	Metric     uint64
	Link_state uint64
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
	Hwassist   uint64
	Oqdrops    uint64
	Lastchange Timeval
***REMOVED***

type IfaMsghdr struct ***REMOVED***
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	Pad_cgo_0 [2]byte
	Metric    int32
***REMOVED***

type IfmaMsghdr struct ***REMOVED***
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	Pad_cgo_0 [2]byte
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
	Inits     uint64
	Rmx       RtMetrics
***REMOVED***

type RtMetrics struct ***REMOVED***
	Locks     uint64
	Mtu       uint64
	Pksent    uint64
	Expire    uint64
	Sendpipe  uint64
	Ssthresh  uint64
	Rtt       uint64
	Rttvar    uint64
	Recvpipe  uint64
	Hopcount  uint64
	Mssopt    uint16
	Pad       uint16
	Pad_cgo_0 [4]byte
	Msl       uint64
	Iwmaxsegs uint64
	Iwcapsegs uint64
***REMOVED***

const (
	SizeofBpfVersion = 0x4
	SizeofBpfStat    = 0x8
	SizeofBpfProgram = 0x10
	SizeofBpfInsn    = 0x8
	SizeofBpfHdr     = 0x20
)

type BpfVersion struct ***REMOVED***
	Major uint16
	Minor uint16
***REMOVED***

type BpfStat struct ***REMOVED***
	Recv uint32
	Drop uint32
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
	Tstamp    Timeval
	Caplen    uint32
	Datalen   uint32
	Hdrlen    uint16
	Pad_cgo_0 [6]byte
***REMOVED***

type Termios struct ***REMOVED***
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]uint8
	Ispeed uint32
	Ospeed uint32
***REMOVED***

type Winsize struct ***REMOVED***
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
***REMOVED***

const (
	AT_FDCWD            = 0xfffafdcd
	AT_SYMLINK_NOFOLLOW = 0x1
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
