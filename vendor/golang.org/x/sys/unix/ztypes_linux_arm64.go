// cgo -godefs -- -Wall -Werror -static -I/tmp/include -fsigned-char linux/types.go | go run mkpost.go
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build arm64,linux

package unix

const (
	SizeofPtr  = 0x8
	SizeofLong = 0x8
)

type (
	_C_long int64
)

type Timespec struct ***REMOVED***
	Sec  int64
	Nsec int64
***REMOVED***

type Timeval struct ***REMOVED***
	Sec  int64
	Usec int64
***REMOVED***

type Timex struct ***REMOVED***
	Modes     uint32
	Offset    int64
	Freq      int64
	Maxerror  int64
	Esterror  int64
	Status    int32
	Constant  int64
	Precision int64
	Tolerance int64
	Time      Timeval
	Tick      int64
	Ppsfreq   int64
	Jitter    int64
	Shift     int32
	Stabil    int64
	Jitcnt    int64
	Calcnt    int64
	Errcnt    int64
	Stbcnt    int64
	Tai       int32
	_         [44]byte
***REMOVED***

type Time_t int64

type Tms struct ***REMOVED***
	Utime  int64
	Stime  int64
	Cutime int64
	Cstime int64
***REMOVED***

type Utimbuf struct ***REMOVED***
	Actime  int64
	Modtime int64
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

type Stat_t struct ***REMOVED***
	Dev     uint64
	Ino     uint64
	Mode    uint32
	Nlink   uint32
	Uid     uint32
	Gid     uint32
	Rdev    uint64
	_       uint64
	Size    int64
	Blksize int32
	_       int32
	Blocks  int64
	Atim    Timespec
	Mtim    Timespec
	Ctim    Timespec
	_       [2]int32
***REMOVED***

type Dirent struct ***REMOVED***
	Ino    uint64
	Off    int64
	Reclen uint16
	Type   uint8
	Name   [256]int8
	_      [5]byte
***REMOVED***

type Flock_t struct ***REMOVED***
	Type   int16
	Whence int16
	Start  int64
	Len    int64
	Pid    int32
	_      [4]byte
***REMOVED***

type DmNameList struct ***REMOVED***
	Dev  uint64
	Next uint32
	Name [0]byte
	_    [4]byte
***REMOVED***

const (
	FADV_DONTNEED = 0x4
	FADV_NOREUSE  = 0x5
)

type RawSockaddr struct ***REMOVED***
	Family uint16
	Data   [14]int8
***REMOVED***

type RawSockaddrAny struct ***REMOVED***
	Addr RawSockaddr
	Pad  [96]int8
***REMOVED***

type Iovec struct ***REMOVED***
	Base *byte
	Len  uint64
***REMOVED***

type Msghdr struct ***REMOVED***
	Name       *byte
	Namelen    uint32
	Iov        *Iovec
	Iovlen     uint64
	Control    *byte
	Controllen uint64
	Flags      int32
	_          [4]byte
***REMOVED***

type Cmsghdr struct ***REMOVED***
	Len   uint64
	Level int32
	Type  int32
***REMOVED***

const (
	SizeofIovec   = 0x10
	SizeofMsghdr  = 0x38
	SizeofCmsghdr = 0x10
)

const (
	SizeofSockFprog = 0x10
)

type PtraceRegs struct ***REMOVED***
	Regs   [31]uint64
	Sp     uint64
	Pc     uint64
	Pstate uint64
***REMOVED***

type FdSet struct ***REMOVED***
	Bits [16]int64
***REMOVED***

type Sysinfo_t struct ***REMOVED***
	Uptime    int64
	Loads     [3]uint64
	Totalram  uint64
	Freeram   uint64
	Sharedram uint64
	Bufferram uint64
	Totalswap uint64
	Freeswap  uint64
	Procs     uint16
	Pad       uint16
	Totalhigh uint64
	Freehigh  uint64
	Unit      uint32
	_         [0]int8
	_         [4]byte
***REMOVED***

type Ustat_t struct ***REMOVED***
	Tfree  int32
	Tinode uint64
	Fname  [6]int8
	Fpack  [6]int8
	_      [4]byte
***REMOVED***

type EpollEvent struct ***REMOVED***
	Events uint32
	PadFd  int32
	Fd     int32
	Pad    int32
***REMOVED***

const (
	POLLRDHUP = 0x2000
)

type Sigset_t struct ***REMOVED***
	Val [16]uint64
***REMOVED***

const _C__NSIG = 0x41

type Termios struct ***REMOVED***
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [19]uint8
	Ispeed uint32
	Ospeed uint32
***REMOVED***

type Taskstats struct ***REMOVED***
	Version                   uint16
	Ac_exitcode               uint32
	Ac_flag                   uint8
	Ac_nice                   uint8
	Cpu_count                 uint64
	Cpu_delay_total           uint64
	Blkio_count               uint64
	Blkio_delay_total         uint64
	Swapin_count              uint64
	Swapin_delay_total        uint64
	Cpu_run_real_total        uint64
	Cpu_run_virtual_total     uint64
	Ac_comm                   [32]int8
	Ac_sched                  uint8
	Ac_pad                    [3]uint8
	_                         [4]byte
	Ac_uid                    uint32
	Ac_gid                    uint32
	Ac_pid                    uint32
	Ac_ppid                   uint32
	Ac_btime                  uint32
	Ac_etime                  uint64
	Ac_utime                  uint64
	Ac_stime                  uint64
	Ac_minflt                 uint64
	Ac_majflt                 uint64
	Coremem                   uint64
	Virtmem                   uint64
	Hiwater_rss               uint64
	Hiwater_vm                uint64
	Read_char                 uint64
	Write_char                uint64
	Read_syscalls             uint64
	Write_syscalls            uint64
	Read_bytes                uint64
	Write_bytes               uint64
	Cancelled_write_bytes     uint64
	Nvcsw                     uint64
	Nivcsw                    uint64
	Ac_utimescaled            uint64
	Ac_stimescaled            uint64
	Cpu_scaled_run_real_total uint64
	Freepages_count           uint64
	Freepages_delay_total     uint64
	Thrashing_count           uint64
	Thrashing_delay_total     uint64
	Ac_btime64                uint64
***REMOVED***

type cpuMask uint64

const (
	_NCPUBITS = 0x40
)

const (
	CBitFieldMaskBit0  = 0x1
	CBitFieldMaskBit1  = 0x2
	CBitFieldMaskBit2  = 0x4
	CBitFieldMaskBit3  = 0x8
	CBitFieldMaskBit4  = 0x10
	CBitFieldMaskBit5  = 0x20
	CBitFieldMaskBit6  = 0x40
	CBitFieldMaskBit7  = 0x80
	CBitFieldMaskBit8  = 0x100
	CBitFieldMaskBit9  = 0x200
	CBitFieldMaskBit10 = 0x400
	CBitFieldMaskBit11 = 0x800
	CBitFieldMaskBit12 = 0x1000
	CBitFieldMaskBit13 = 0x2000
	CBitFieldMaskBit14 = 0x4000
	CBitFieldMaskBit15 = 0x8000
	CBitFieldMaskBit16 = 0x10000
	CBitFieldMaskBit17 = 0x20000
	CBitFieldMaskBit18 = 0x40000
	CBitFieldMaskBit19 = 0x80000
	CBitFieldMaskBit20 = 0x100000
	CBitFieldMaskBit21 = 0x200000
	CBitFieldMaskBit22 = 0x400000
	CBitFieldMaskBit23 = 0x800000
	CBitFieldMaskBit24 = 0x1000000
	CBitFieldMaskBit25 = 0x2000000
	CBitFieldMaskBit26 = 0x4000000
	CBitFieldMaskBit27 = 0x8000000
	CBitFieldMaskBit28 = 0x10000000
	CBitFieldMaskBit29 = 0x20000000
	CBitFieldMaskBit30 = 0x40000000
	CBitFieldMaskBit31 = 0x80000000
	CBitFieldMaskBit32 = 0x100000000
	CBitFieldMaskBit33 = 0x200000000
	CBitFieldMaskBit34 = 0x400000000
	CBitFieldMaskBit35 = 0x800000000
	CBitFieldMaskBit36 = 0x1000000000
	CBitFieldMaskBit37 = 0x2000000000
	CBitFieldMaskBit38 = 0x4000000000
	CBitFieldMaskBit39 = 0x8000000000
	CBitFieldMaskBit40 = 0x10000000000
	CBitFieldMaskBit41 = 0x20000000000
	CBitFieldMaskBit42 = 0x40000000000
	CBitFieldMaskBit43 = 0x80000000000
	CBitFieldMaskBit44 = 0x100000000000
	CBitFieldMaskBit45 = 0x200000000000
	CBitFieldMaskBit46 = 0x400000000000
	CBitFieldMaskBit47 = 0x800000000000
	CBitFieldMaskBit48 = 0x1000000000000
	CBitFieldMaskBit49 = 0x2000000000000
	CBitFieldMaskBit50 = 0x4000000000000
	CBitFieldMaskBit51 = 0x8000000000000
	CBitFieldMaskBit52 = 0x10000000000000
	CBitFieldMaskBit53 = 0x20000000000000
	CBitFieldMaskBit54 = 0x40000000000000
	CBitFieldMaskBit55 = 0x80000000000000
	CBitFieldMaskBit56 = 0x100000000000000
	CBitFieldMaskBit57 = 0x200000000000000
	CBitFieldMaskBit58 = 0x400000000000000
	CBitFieldMaskBit59 = 0x800000000000000
	CBitFieldMaskBit60 = 0x1000000000000000
	CBitFieldMaskBit61 = 0x2000000000000000
	CBitFieldMaskBit62 = 0x4000000000000000
	CBitFieldMaskBit63 = 0x8000000000000000
)

type SockaddrStorage struct ***REMOVED***
	Family uint16
	_      [118]int8
	_      uint64
***REMOVED***

type HDGeometry struct ***REMOVED***
	Heads     uint8
	Sectors   uint8
	Cylinders uint16
	Start     uint64
***REMOVED***

type Statfs_t struct ***REMOVED***
	Type    int64
	Bsize   int64
	Blocks  uint64
	Bfree   uint64
	Bavail  uint64
	Files   uint64
	Ffree   uint64
	Fsid    Fsid
	Namelen int64
	Frsize  int64
	Flags   int64
	Spare   [4]int64
***REMOVED***

type TpacketHdr struct ***REMOVED***
	Status  uint64
	Len     uint32
	Snaplen uint32
	Mac     uint16
	Net     uint16
	Sec     uint32
	Usec    uint32
	_       [4]byte
***REMOVED***

const (
	SizeofTpacketHdr = 0x20
)

type RTCPLLInfo struct ***REMOVED***
	Ctrl    int32
	Value   int32
	Max     int32
	Min     int32
	Posmult int32
	Negmult int32
	Clock   int64
***REMOVED***

type BlkpgPartition struct ***REMOVED***
	Start   int64
	Length  int64
	Pno     int32
	Devname [64]uint8
	Volname [64]uint8
	_       [4]byte
***REMOVED***

const (
	BLKPG = 0x1269
)

type XDPUmemReg struct ***REMOVED***
	Addr     uint64
	Len      uint64
	Size     uint32
	Headroom uint32
	Flags    uint32
	_        [4]byte
***REMOVED***

type CryptoUserAlg struct ***REMOVED***
	Name        [64]int8
	Driver_name [64]int8
	Module_name [64]int8
	Type        uint32
	Mask        uint32
	Refcnt      uint32
	Flags       uint32
***REMOVED***

type CryptoStatAEAD struct ***REMOVED***
	Type         [64]int8
	Encrypt_cnt  uint64
	Encrypt_tlen uint64
	Decrypt_cnt  uint64
	Decrypt_tlen uint64
	Err_cnt      uint64
***REMOVED***

type CryptoStatAKCipher struct ***REMOVED***
	Type         [64]int8
	Encrypt_cnt  uint64
	Encrypt_tlen uint64
	Decrypt_cnt  uint64
	Decrypt_tlen uint64
	Verify_cnt   uint64
	Sign_cnt     uint64
	Err_cnt      uint64
***REMOVED***

type CryptoStatCipher struct ***REMOVED***
	Type         [64]int8
	Encrypt_cnt  uint64
	Encrypt_tlen uint64
	Decrypt_cnt  uint64
	Decrypt_tlen uint64
	Err_cnt      uint64
***REMOVED***

type CryptoStatCompress struct ***REMOVED***
	Type            [64]int8
	Compress_cnt    uint64
	Compress_tlen   uint64
	Decompress_cnt  uint64
	Decompress_tlen uint64
	Err_cnt         uint64
***REMOVED***

type CryptoStatHash struct ***REMOVED***
	Type      [64]int8
	Hash_cnt  uint64
	Hash_tlen uint64
	Err_cnt   uint64
***REMOVED***

type CryptoStatKPP struct ***REMOVED***
	Type                      [64]int8
	Setsecret_cnt             uint64
	Generate_public_key_cnt   uint64
	Compute_shared_secret_cnt uint64
	Err_cnt                   uint64
***REMOVED***

type CryptoStatRNG struct ***REMOVED***
	Type          [64]int8
	Generate_cnt  uint64
	Generate_tlen uint64
	Seed_cnt      uint64
	Err_cnt       uint64
***REMOVED***

type CryptoStatLarval struct ***REMOVED***
	Type [64]int8
***REMOVED***

type CryptoReportLarval struct ***REMOVED***
	Type [64]int8
***REMOVED***

type CryptoReportHash struct ***REMOVED***
	Type       [64]int8
	Blocksize  uint32
	Digestsize uint32
***REMOVED***

type CryptoReportCipher struct ***REMOVED***
	Type        [64]int8
	Blocksize   uint32
	Min_keysize uint32
	Max_keysize uint32
***REMOVED***

type CryptoReportBlkCipher struct ***REMOVED***
	Type        [64]int8
	Geniv       [64]int8
	Blocksize   uint32
	Min_keysize uint32
	Max_keysize uint32
	Ivsize      uint32
***REMOVED***

type CryptoReportAEAD struct ***REMOVED***
	Type        [64]int8
	Geniv       [64]int8
	Blocksize   uint32
	Maxauthsize uint32
	Ivsize      uint32
***REMOVED***

type CryptoReportComp struct ***REMOVED***
	Type [64]int8
***REMOVED***

type CryptoReportRNG struct ***REMOVED***
	Type     [64]int8
	Seedsize uint32
***REMOVED***

type CryptoReportAKCipher struct ***REMOVED***
	Type [64]int8
***REMOVED***

type CryptoReportKPP struct ***REMOVED***
	Type [64]int8
***REMOVED***

type CryptoReportAcomp struct ***REMOVED***
	Type [64]int8
***REMOVED***

type LoopInfo struct ***REMOVED***
	Number           int32
	Device           uint32
	Inode            uint64
	Rdevice          uint32
	Offset           int32
	Encrypt_type     int32
	Encrypt_key_size int32
	Flags            int32
	Name             [64]int8
	Encrypt_key      [32]uint8
	Init             [2]uint64
	Reserved         [4]int8
	_                [4]byte
***REMOVED***

type TIPCSubscr struct ***REMOVED***
	Seq     TIPCServiceRange
	Timeout uint32
	Filter  uint32
	Handle  [8]int8
***REMOVED***

type TIPCSIOCLNReq struct ***REMOVED***
	Peer     uint32
	Id       uint32
	Linkname [68]int8
***REMOVED***

type TIPCSIOCNodeIDReq struct ***REMOVED***
	Peer uint32
	Id   [16]int8
***REMOVED***

type PPSKInfo struct ***REMOVED***
	Assert_sequence uint32
	Clear_sequence  uint32
	Assert_tu       PPSKTime
	Clear_tu        PPSKTime
	Current_mode    int32
	_               [4]byte
***REMOVED***

const (
	PPS_GETPARAMS = 0x800870a1
	PPS_SETPARAMS = 0x400870a2
	PPS_GETCAP    = 0x800870a3
	PPS_FETCH     = 0xc00870a4
)
