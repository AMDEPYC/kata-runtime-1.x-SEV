// cgo -godefs -- -Wall -Werror -static -I/tmp/include linux/types.go | go run mkpost.go
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build mipsle,linux

package unix

const (
	sizeofPtr      = 0x4
	sizeofShort    = 0x2
	sizeofInt      = 0x4
	sizeofLong     = 0x4
	sizeofLongLong = 0x8
	PathMax        = 0x1000
)

type (
	_C_short     int16
	_C_int       int32
	_C_long      int32
	_C_long_long int64
)

type Timespec struct {
	Sec  int32
	Nsec int32
}

type Timeval struct {
	Sec  int32
	Usec int32
}

type Timex struct {
	Modes     uint32
	Offset    int32
	Freq      int32
	Maxerror  int32
	Esterror  int32
	Status    int32
	Constant  int32
	Precision int32
	Tolerance int32
	Time      Timeval
	Tick      int32
	Ppsfreq   int32
	Jitter    int32
	Shift     int32
	Stabil    int32
	Jitcnt    int32
	Calcnt    int32
	Errcnt    int32
	Stbcnt    int32
	Tai       int32
	Pad_cgo_0 [44]byte
}

type Time_t int32

type Tms struct {
	Utime  int32
	Stime  int32
	Cutime int32
	Cstime int32
}

type Utimbuf struct {
	Actime  int32
	Modtime int32
}

type Rusage struct {
	Utime    Timeval
	Stime    Timeval
	Maxrss   int32
	Ixrss    int32
	Idrss    int32
	Isrss    int32
	Minflt   int32
	Majflt   int32
	Nswap    int32
	Inblock  int32
	Oublock  int32
	Msgsnd   int32
	Msgrcv   int32
	Nsignals int32
	Nvcsw    int32
	Nivcsw   int32
}

type Rlimit struct {
	Cur uint64
	Max uint64
}

type _Gid_t uint32

type Stat_t struct {
	Dev     uint32
	Pad1    [3]int32
	Ino     uint64
	Mode    uint32
	Nlink   uint32
	Uid     uint32
	Gid     uint32
	Rdev    uint32
	Pad2    [3]int32
	Size    int64
	Atim    Timespec
	Mtim    Timespec
	Ctim    Timespec
	Blksize int32
	Pad4    int32
	Blocks  int64
	Pad5    [14]int32
}

type Statfs_t struct {
	Type      int32
	Bsize     int32
	Frsize    int32
	Pad_cgo_0 [4]byte
	Blocks    uint64
	Bfree     uint64
	Files     uint64
	Ffree     uint64
	Bavail    uint64
	Fsid      Fsid
	Namelen   int32
	Flags     int32
	Spare     [5]int32
	Pad_cgo_1 [4]byte
}

type Dirent struct {
	Ino       uint64
	Off       int64
	Reclen    uint16
	Type      uint8
	Name      [256]int8
	Pad_cgo_0 [5]byte
}

type Fsid struct {
	X__val [2]int32
}

type Flock_t struct {
	Type      int16
	Whence    int16
	Pad_cgo_0 [4]byte
	Start     int64
	Len       int64
	Pid       int32
	Pad_cgo_1 [4]byte
}

type FscryptPolicy struct {
	Version                   uint8
	Contents_encryption_mode  uint8
	Filenames_encryption_mode uint8
	Flags                     uint8
	Master_key_descriptor     [8]uint8
}

type FscryptKey struct {
	Mode uint32
	Raw  [64]uint8
	Size uint32
}

type KeyctlDHParams struct {
	Private int32
	Prime   int32
	Base    int32
}

const (
	FADV_NORMAL     = 0x0
	FADV_RANDOM     = 0x1
	FADV_SEQUENTIAL = 0x2
	FADV_WILLNEED   = 0x3
	FADV_DONTNEED   = 0x4
	FADV_NOREUSE    = 0x5
)

type RawSockaddrInet4 struct {
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]uint8
}

type RawSockaddrInet6 struct {
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

type RawSockaddrUnix struct {
	Family uint16
	Path   [108]int8
}

type RawSockaddrLinklayer struct {
	Family   uint16
	Protocol uint16
	Ifindex  int32
	Hatype   uint16
	Pkttype  uint8
	Halen    uint8
	Addr     [8]uint8
}

type RawSockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
}

type RawSockaddrHCI struct {
	Family  uint16
	Dev     uint16
	Channel uint16
}

type RawSockaddrCAN struct {
	Family    uint16
	Pad_cgo_0 [2]byte
	Ifindex   int32
	Addr      [8]byte
}

type RawSockaddrALG struct {
	Family uint16
	Type   [14]uint8
	Feat   uint32
	Mask   uint32
	Name   [64]uint8
}

type RawSockaddrVM struct {
	Family    uint16
	Reserved1 uint16
	Port      uint32
	Cid       uint32
	Zero      [4]uint8
}

type RawSockaddr struct {
	Family uint16
	Data   [14]int8
}

type RawSockaddrAny struct {
	Addr RawSockaddr
	Pad  [96]int8
}

type _Socklen uint32

type Linger struct {
	Onoff  int32
	Linger int32
}

type Iovec struct {
	Base *byte
	Len  uint32
}

type IPMreq struct {
	Multiaddr [4]byte /* in_addr */
	Interface [4]byte /* in_addr */
}

type IPMreqn struct {
	Multiaddr [4]byte /* in_addr */
	Address   [4]byte /* in_addr */
	Ifindex   int32
}

type IPv6Mreq struct {
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
}

type PacketMreq struct {
	Ifindex int32
	Type    uint16
	Alen    uint16
	Address [8]uint8
}

type Msghdr struct {
	Name       *byte
	Namelen    uint32
	Iov        *Iovec
	Iovlen     uint32
	Control    *byte
	Controllen uint32
	Flags      int32
}

type Cmsghdr struct {
	Len   uint32
	Level int32
	Type  int32
}

type Inet4Pktinfo struct {
	Ifindex  int32
	Spec_dst [4]byte /* in_addr */
	Addr     [4]byte /* in_addr */
}

type Inet6Pktinfo struct {
	Addr    [16]byte /* in6_addr */
	Ifindex uint32
}

type IPv6MTUInfo struct {
	Addr RawSockaddrInet6
	Mtu  uint32
}

type ICMPv6Filter struct {
	Data [8]uint32
}

type Ucred struct {
	Pid int32
	Uid uint32
	Gid uint32
}

type TCPInfo struct {
	State          uint8
	Ca_state       uint8
	Retransmits    uint8
	Probes         uint8
	Backoff        uint8
	Options        uint8
	Pad_cgo_0      [2]byte
	Rto            uint32
	Ato            uint32
	Snd_mss        uint32
	Rcv_mss        uint32
	Unacked        uint32
	Sacked         uint32
	Lost           uint32
	Retrans        uint32
	Fackets        uint32
	Last_data_sent uint32
	Last_ack_sent  uint32
	Last_data_recv uint32
	Last_ack_recv  uint32
	Pmtu           uint32
	Rcv_ssthresh   uint32
	Rtt            uint32
	Rttvar         uint32
	Snd_ssthresh   uint32
	Snd_cwnd       uint32
	Advmss         uint32
	Reordering     uint32
	Rcv_rtt        uint32
	Rcv_space      uint32
	Total_retrans  uint32
}

const (
	SizeofSockaddrInet4     = 0x10
	SizeofSockaddrInet6     = 0x1c
	SizeofSockaddrAny       = 0x70
	SizeofSockaddrUnix      = 0x6e
	SizeofSockaddrLinklayer = 0x14
	SizeofSockaddrNetlink   = 0xc
	SizeofSockaddrHCI       = 0x6
	SizeofSockaddrCAN       = 0x10
	SizeofSockaddrALG       = 0x58
	SizeofSockaddrVM        = 0x10
	SizeofLinger            = 0x8
	SizeofIovec             = 0x8
	SizeofIPMreq            = 0x8
	SizeofIPMreqn           = 0xc
	SizeofIPv6Mreq          = 0x14
	SizeofPacketMreq        = 0x10
	SizeofMsghdr            = 0x1c
	SizeofCmsghdr           = 0xc
	SizeofInet4Pktinfo      = 0xc
	SizeofInet6Pktinfo      = 0x14
	SizeofIPv6MTUInfo       = 0x20
	SizeofICMPv6Filter      = 0x20
	SizeofUcred             = 0xc
	SizeofTCPInfo           = 0x68
)

const (
	IFA_UNSPEC          = 0x0
	IFA_ADDRESS         = 0x1
	IFA_LOCAL           = 0x2
	IFA_LABEL           = 0x3
	IFA_BROADCAST       = 0x4
	IFA_ANYCAST         = 0x5
	IFA_CACHEINFO       = 0x6
	IFA_MULTICAST       = 0x7
	IFLA_UNSPEC         = 0x0
	IFLA_ADDRESS        = 0x1
	IFLA_BROADCAST      = 0x2
	IFLA_IFNAME         = 0x3
	IFLA_MTU            = 0x4
	IFLA_LINK           = 0x5
	IFLA_QDISC          = 0x6
	IFLA_STATS          = 0x7
	IFLA_COST           = 0x8
	IFLA_PRIORITY       = 0x9
	IFLA_MASTER         = 0xa
	IFLA_WIRELESS       = 0xb
	IFLA_PROTINFO       = 0xc
	IFLA_TXQLEN         = 0xd
	IFLA_MAP            = 0xe
	IFLA_WEIGHT         = 0xf
	IFLA_OPERSTATE      = 0x10
	IFLA_LINKMODE       = 0x11
	IFLA_LINKINFO       = 0x12
	IFLA_NET_NS_PID     = 0x13
	IFLA_IFALIAS        = 0x14
	IFLA_MAX            = 0x2c
	RT_SCOPE_UNIVERSE   = 0x0
	RT_SCOPE_SITE       = 0xc8
	RT_SCOPE_LINK       = 0xfd
	RT_SCOPE_HOST       = 0xfe
	RT_SCOPE_NOWHERE    = 0xff
	RT_TABLE_UNSPEC     = 0x0
	RT_TABLE_COMPAT     = 0xfc
	RT_TABLE_DEFAULT    = 0xfd
	RT_TABLE_MAIN       = 0xfe
	RT_TABLE_LOCAL      = 0xff
	RT_TABLE_MAX        = 0xffffffff
	RTA_UNSPEC          = 0x0
	RTA_DST             = 0x1
	RTA_SRC             = 0x2
	RTA_IIF             = 0x3
	RTA_OIF             = 0x4
	RTA_GATEWAY         = 0x5
	RTA_PRIORITY        = 0x6
	RTA_PREFSRC         = 0x7
	RTA_METRICS         = 0x8
	RTA_MULTIPATH       = 0x9
	RTA_FLOW            = 0xb
	RTA_CACHEINFO       = 0xc
	RTA_TABLE           = 0xf
	RTN_UNSPEC          = 0x0
	RTN_UNICAST         = 0x1
	RTN_LOCAL           = 0x2
	RTN_BROADCAST       = 0x3
	RTN_ANYCAST         = 0x4
	RTN_MULTICAST       = 0x5
	RTN_BLACKHOLE       = 0x6
	RTN_UNREACHABLE     = 0x7
	RTN_PROHIBIT        = 0x8
	RTN_THROW           = 0x9
	RTN_NAT             = 0xa
	RTN_XRESOLVE        = 0xb
	RTNLGRP_NONE        = 0x0
	RTNLGRP_LINK        = 0x1
	RTNLGRP_NOTIFY      = 0x2
	RTNLGRP_NEIGH       = 0x3
	RTNLGRP_TC          = 0x4
	RTNLGRP_IPV4_IFADDR = 0x5
	RTNLGRP_IPV4_MROUTE = 0x6
	RTNLGRP_IPV4_ROUTE  = 0x7
	RTNLGRP_IPV4_RULE   = 0x8
	RTNLGRP_IPV6_IFADDR = 0x9
	RTNLGRP_IPV6_MROUTE = 0xa
	RTNLGRP_IPV6_ROUTE  = 0xb
	RTNLGRP_IPV6_IFINFO = 0xc
	RTNLGRP_IPV6_PREFIX = 0x12
	RTNLGRP_IPV6_RULE   = 0x13
	RTNLGRP_ND_USEROPT  = 0x14
	SizeofNlMsghdr      = 0x10
	SizeofNlMsgerr      = 0x14
	SizeofRtGenmsg      = 0x1
	SizeofNlAttr        = 0x4
	SizeofRtAttr        = 0x4
	SizeofIfInfomsg     = 0x10
	SizeofIfAddrmsg     = 0x8
	SizeofRtMsg         = 0xc
	SizeofRtNexthop     = 0x8
)

type NlMsghdr struct {
	Len   uint32
	Type  uint16
	Flags uint16
	Seq   uint32
	Pid   uint32
}

type NlMsgerr struct {
	Error int32
	Msg   NlMsghdr
}

type RtGenmsg struct {
	Family uint8
}

type NlAttr struct {
	Len  uint16
	Type uint16
}

type RtAttr struct {
	Len  uint16
	Type uint16
}

type IfInfomsg struct {
	Family     uint8
	X__ifi_pad uint8
	Type       uint16
	Index      int32
	Flags      uint32
	Change     uint32
}

type IfAddrmsg struct {
	Family    uint8
	Prefixlen uint8
	Flags     uint8
	Scope     uint8
	Index     uint32
}

type RtMsg struct {
	Family   uint8
	Dst_len  uint8
	Src_len  uint8
	Tos      uint8
	Table    uint8
	Protocol uint8
	Scope    uint8
	Type     uint8
	Flags    uint32
}

type RtNexthop struct {
	Len     uint16
	Flags   uint8
	Hops    uint8
	Ifindex int32
}

const (
	SizeofSockFilter = 0x8
	SizeofSockFprog  = 0x8
)

type SockFilter struct {
	Code uint16
	Jt   uint8
	Jf   uint8
	K    uint32
}

type SockFprog struct {
	Len       uint16
	Pad_cgo_0 [2]byte
	Filter    *SockFilter
}

type InotifyEvent struct {
	Wd     int32
	Mask   uint32
	Cookie uint32
	Len    uint32
}

const SizeofInotifyEvent = 0x10

type PtraceRegs struct {
	Regs     [32]uint64
	Lo       uint64
	Hi       uint64
	Epc      uint64
	Badvaddr uint64
	Status   uint64
	Cause    uint64
}

type FdSet struct {
	Bits [32]int32
}

type Sysinfo_t struct {
	Uptime    int32
	Loads     [3]uint32
	Totalram  uint32
	Freeram   uint32
	Sharedram uint32
	Bufferram uint32
	Totalswap uint32
	Freeswap  uint32
	Procs     uint16
	Pad       uint16
	Totalhigh uint32
	Freehigh  uint32
	Unit      uint32
	X_f       [8]int8
}

type Utsname struct {
	Sysname    [65]byte
	Nodename   [65]byte
	Release    [65]byte
	Version    [65]byte
	Machine    [65]byte
	Domainname [65]byte
}

type Ustat_t struct {
	Tfree  int32
	Tinode uint32
	Fname  [6]int8
	Fpack  [6]int8
}

type EpollEvent struct {
	Events uint32
	PadFd  int32
	Fd     int32
	Pad    int32
}

const (
	AT_FDCWD            = -0x64
	AT_NO_AUTOMOUNT     = 0x800
	AT_REMOVEDIR        = 0x200
	AT_SYMLINK_FOLLOW   = 0x400
	AT_SYMLINK_NOFOLLOW = 0x100
)

type PollFd struct {
	Fd      int32
	Events  int16
	Revents int16
}

const (
	POLLIN    = 0x1
	POLLPRI   = 0x2
	POLLOUT   = 0x4
	POLLRDHUP = 0x2000
	POLLERR   = 0x8
	POLLHUP   = 0x10
	POLLNVAL  = 0x20
)

type Sigset_t struct {
	X__val [32]uint32
}

const RNDGETENTCNT = 0x40045200

const PERF_IOC_FLAG_GROUP = 0x1

type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [23]uint8
	Ispeed uint32
	Ospeed uint32
}

type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

type Taskstats struct {
	Version                   uint16
	Pad_cgo_0                 [2]byte
	Ac_exitcode               uint32
	Ac_flag                   uint8
	Ac_nice                   uint8
	Pad_cgo_1                 [6]byte
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
	Pad_cgo_2                 [4]byte
	Ac_uid                    uint32
	Ac_gid                    uint32
	Ac_pid                    uint32
	Ac_ppid                   uint32
	Ac_btime                  uint32
	Pad_cgo_3                 [4]byte
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
}

const (
	TASKSTATS_CMD_UNSPEC                  = 0x0
	TASKSTATS_CMD_GET                     = 0x1
	TASKSTATS_CMD_NEW                     = 0x2
	TASKSTATS_TYPE_UNSPEC                 = 0x0
	TASKSTATS_TYPE_PID                    = 0x1
	TASKSTATS_TYPE_TGID                   = 0x2
	TASKSTATS_TYPE_STATS                  = 0x3
	TASKSTATS_TYPE_AGGR_PID               = 0x4
	TASKSTATS_TYPE_AGGR_TGID              = 0x5
	TASKSTATS_TYPE_NULL                   = 0x6
	TASKSTATS_CMD_ATTR_UNSPEC             = 0x0
	TASKSTATS_CMD_ATTR_PID                = 0x1
	TASKSTATS_CMD_ATTR_TGID               = 0x2
	TASKSTATS_CMD_ATTR_REGISTER_CPUMASK   = 0x3
	TASKSTATS_CMD_ATTR_DEREGISTER_CPUMASK = 0x4
)

type Genlmsghdr struct {
	Cmd      uint8
	Version  uint8
	Reserved uint16
}

const (
	CTRL_CMD_UNSPEC            = 0x0
	CTRL_CMD_NEWFAMILY         = 0x1
	CTRL_CMD_DELFAMILY         = 0x2
	CTRL_CMD_GETFAMILY         = 0x3
	CTRL_CMD_NEWOPS            = 0x4
	CTRL_CMD_DELOPS            = 0x5
	CTRL_CMD_GETOPS            = 0x6
	CTRL_CMD_NEWMCAST_GRP      = 0x7
	CTRL_CMD_DELMCAST_GRP      = 0x8
	CTRL_CMD_GETMCAST_GRP      = 0x9
	CTRL_ATTR_UNSPEC           = 0x0
	CTRL_ATTR_FAMILY_ID        = 0x1
	CTRL_ATTR_FAMILY_NAME      = 0x2
	CTRL_ATTR_VERSION          = 0x3
	CTRL_ATTR_HDRSIZE          = 0x4
	CTRL_ATTR_MAXATTR          = 0x5
	CTRL_ATTR_OPS              = 0x6
	CTRL_ATTR_MCAST_GROUPS     = 0x7
	CTRL_ATTR_OP_UNSPEC        = 0x0
	CTRL_ATTR_OP_ID            = 0x1
	CTRL_ATTR_OP_FLAGS         = 0x2
	CTRL_ATTR_MCAST_GRP_UNSPEC = 0x0
	CTRL_ATTR_MCAST_GRP_NAME   = 0x1
	CTRL_ATTR_MCAST_GRP_ID     = 0x2
)
