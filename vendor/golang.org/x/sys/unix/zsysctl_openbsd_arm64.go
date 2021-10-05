// go run mksysctl_openbsd.go
// Code generated by the command above; DO NOT EDIT.

//go:build arm64 && openbsd
// +build arm64,openbsd

package unix

type mibentry struct ***REMOVED***
	ctlname string
	ctloid  []_C_int
***REMOVED***

var sysctlMib = []mibentry***REMOVED***
	***REMOVED***"ddb.console", []_C_int***REMOVED***9, 6***REMOVED******REMOVED***,
	***REMOVED***"ddb.log", []_C_int***REMOVED***9, 7***REMOVED******REMOVED***,
	***REMOVED***"ddb.max_line", []_C_int***REMOVED***9, 3***REMOVED******REMOVED***,
	***REMOVED***"ddb.max_width", []_C_int***REMOVED***9, 2***REMOVED******REMOVED***,
	***REMOVED***"ddb.panic", []_C_int***REMOVED***9, 5***REMOVED******REMOVED***,
	***REMOVED***"ddb.profile", []_C_int***REMOVED***9, 9***REMOVED******REMOVED***,
	***REMOVED***"ddb.radix", []_C_int***REMOVED***9, 1***REMOVED******REMOVED***,
	***REMOVED***"ddb.tab_stop_width", []_C_int***REMOVED***9, 4***REMOVED******REMOVED***,
	***REMOVED***"ddb.trigger", []_C_int***REMOVED***9, 8***REMOVED******REMOVED***,
	***REMOVED***"fs.posix.setuid", []_C_int***REMOVED***3, 1, 1***REMOVED******REMOVED***,
	***REMOVED***"hw.allowpowerdown", []_C_int***REMOVED***6, 22***REMOVED******REMOVED***,
	***REMOVED***"hw.byteorder", []_C_int***REMOVED***6, 4***REMOVED******REMOVED***,
	***REMOVED***"hw.cpuspeed", []_C_int***REMOVED***6, 12***REMOVED******REMOVED***,
	***REMOVED***"hw.diskcount", []_C_int***REMOVED***6, 10***REMOVED******REMOVED***,
	***REMOVED***"hw.disknames", []_C_int***REMOVED***6, 8***REMOVED******REMOVED***,
	***REMOVED***"hw.diskstats", []_C_int***REMOVED***6, 9***REMOVED******REMOVED***,
	***REMOVED***"hw.machine", []_C_int***REMOVED***6, 1***REMOVED******REMOVED***,
	***REMOVED***"hw.model", []_C_int***REMOVED***6, 2***REMOVED******REMOVED***,
	***REMOVED***"hw.ncpu", []_C_int***REMOVED***6, 3***REMOVED******REMOVED***,
	***REMOVED***"hw.ncpufound", []_C_int***REMOVED***6, 21***REMOVED******REMOVED***,
	***REMOVED***"hw.ncpuonline", []_C_int***REMOVED***6, 25***REMOVED******REMOVED***,
	***REMOVED***"hw.pagesize", []_C_int***REMOVED***6, 7***REMOVED******REMOVED***,
	***REMOVED***"hw.perfpolicy", []_C_int***REMOVED***6, 23***REMOVED******REMOVED***,
	***REMOVED***"hw.physmem", []_C_int***REMOVED***6, 19***REMOVED******REMOVED***,
	***REMOVED***"hw.product", []_C_int***REMOVED***6, 15***REMOVED******REMOVED***,
	***REMOVED***"hw.serialno", []_C_int***REMOVED***6, 17***REMOVED******REMOVED***,
	***REMOVED***"hw.setperf", []_C_int***REMOVED***6, 13***REMOVED******REMOVED***,
	***REMOVED***"hw.smt", []_C_int***REMOVED***6, 24***REMOVED******REMOVED***,
	***REMOVED***"hw.usermem", []_C_int***REMOVED***6, 20***REMOVED******REMOVED***,
	***REMOVED***"hw.uuid", []_C_int***REMOVED***6, 18***REMOVED******REMOVED***,
	***REMOVED***"hw.vendor", []_C_int***REMOVED***6, 14***REMOVED******REMOVED***,
	***REMOVED***"hw.version", []_C_int***REMOVED***6, 16***REMOVED******REMOVED***,
	***REMOVED***"kern.allowkmem", []_C_int***REMOVED***1, 52***REMOVED******REMOVED***,
	***REMOVED***"kern.argmax", []_C_int***REMOVED***1, 8***REMOVED******REMOVED***,
	***REMOVED***"kern.audio", []_C_int***REMOVED***1, 84***REMOVED******REMOVED***,
	***REMOVED***"kern.boottime", []_C_int***REMOVED***1, 21***REMOVED******REMOVED***,
	***REMOVED***"kern.bufcachepercent", []_C_int***REMOVED***1, 72***REMOVED******REMOVED***,
	***REMOVED***"kern.ccpu", []_C_int***REMOVED***1, 45***REMOVED******REMOVED***,
	***REMOVED***"kern.clockrate", []_C_int***REMOVED***1, 12***REMOVED******REMOVED***,
	***REMOVED***"kern.consdev", []_C_int***REMOVED***1, 75***REMOVED******REMOVED***,
	***REMOVED***"kern.cp_time", []_C_int***REMOVED***1, 40***REMOVED******REMOVED***,
	***REMOVED***"kern.cp_time2", []_C_int***REMOVED***1, 71***REMOVED******REMOVED***,
	***REMOVED***"kern.cpustats", []_C_int***REMOVED***1, 85***REMOVED******REMOVED***,
	***REMOVED***"kern.domainname", []_C_int***REMOVED***1, 22***REMOVED******REMOVED***,
	***REMOVED***"kern.file", []_C_int***REMOVED***1, 73***REMOVED******REMOVED***,
	***REMOVED***"kern.forkstat", []_C_int***REMOVED***1, 42***REMOVED******REMOVED***,
	***REMOVED***"kern.fscale", []_C_int***REMOVED***1, 46***REMOVED******REMOVED***,
	***REMOVED***"kern.fsync", []_C_int***REMOVED***1, 33***REMOVED******REMOVED***,
	***REMOVED***"kern.global_ptrace", []_C_int***REMOVED***1, 81***REMOVED******REMOVED***,
	***REMOVED***"kern.hostid", []_C_int***REMOVED***1, 11***REMOVED******REMOVED***,
	***REMOVED***"kern.hostname", []_C_int***REMOVED***1, 10***REMOVED******REMOVED***,
	***REMOVED***"kern.intrcnt.nintrcnt", []_C_int***REMOVED***1, 63, 1***REMOVED******REMOVED***,
	***REMOVED***"kern.job_control", []_C_int***REMOVED***1, 19***REMOVED******REMOVED***,
	***REMOVED***"kern.malloc.buckets", []_C_int***REMOVED***1, 39, 1***REMOVED******REMOVED***,
	***REMOVED***"kern.malloc.kmemnames", []_C_int***REMOVED***1, 39, 3***REMOVED******REMOVED***,
	***REMOVED***"kern.maxclusters", []_C_int***REMOVED***1, 67***REMOVED******REMOVED***,
	***REMOVED***"kern.maxfiles", []_C_int***REMOVED***1, 7***REMOVED******REMOVED***,
	***REMOVED***"kern.maxlocksperuid", []_C_int***REMOVED***1, 70***REMOVED******REMOVED***,
	***REMOVED***"kern.maxpartitions", []_C_int***REMOVED***1, 23***REMOVED******REMOVED***,
	***REMOVED***"kern.maxproc", []_C_int***REMOVED***1, 6***REMOVED******REMOVED***,
	***REMOVED***"kern.maxthread", []_C_int***REMOVED***1, 25***REMOVED******REMOVED***,
	***REMOVED***"kern.maxvnodes", []_C_int***REMOVED***1, 5***REMOVED******REMOVED***,
	***REMOVED***"kern.mbstat", []_C_int***REMOVED***1, 59***REMOVED******REMOVED***,
	***REMOVED***"kern.msgbuf", []_C_int***REMOVED***1, 48***REMOVED******REMOVED***,
	***REMOVED***"kern.msgbufsize", []_C_int***REMOVED***1, 38***REMOVED******REMOVED***,
	***REMOVED***"kern.nchstats", []_C_int***REMOVED***1, 41***REMOVED******REMOVED***,
	***REMOVED***"kern.netlivelocks", []_C_int***REMOVED***1, 76***REMOVED******REMOVED***,
	***REMOVED***"kern.nfiles", []_C_int***REMOVED***1, 56***REMOVED******REMOVED***,
	***REMOVED***"kern.ngroups", []_C_int***REMOVED***1, 18***REMOVED******REMOVED***,
	***REMOVED***"kern.nosuidcoredump", []_C_int***REMOVED***1, 32***REMOVED******REMOVED***,
	***REMOVED***"kern.nprocs", []_C_int***REMOVED***1, 47***REMOVED******REMOVED***,
	***REMOVED***"kern.nselcoll", []_C_int***REMOVED***1, 43***REMOVED******REMOVED***,
	***REMOVED***"kern.nthreads", []_C_int***REMOVED***1, 26***REMOVED******REMOVED***,
	***REMOVED***"kern.numvnodes", []_C_int***REMOVED***1, 58***REMOVED******REMOVED***,
	***REMOVED***"kern.osrelease", []_C_int***REMOVED***1, 2***REMOVED******REMOVED***,
	***REMOVED***"kern.osrevision", []_C_int***REMOVED***1, 3***REMOVED******REMOVED***,
	***REMOVED***"kern.ostype", []_C_int***REMOVED***1, 1***REMOVED******REMOVED***,
	***REMOVED***"kern.osversion", []_C_int***REMOVED***1, 27***REMOVED******REMOVED***,
	***REMOVED***"kern.pool_debug", []_C_int***REMOVED***1, 77***REMOVED******REMOVED***,
	***REMOVED***"kern.posix1version", []_C_int***REMOVED***1, 17***REMOVED******REMOVED***,
	***REMOVED***"kern.proc", []_C_int***REMOVED***1, 66***REMOVED******REMOVED***,
	***REMOVED***"kern.rawpartition", []_C_int***REMOVED***1, 24***REMOVED******REMOVED***,
	***REMOVED***"kern.saved_ids", []_C_int***REMOVED***1, 20***REMOVED******REMOVED***,
	***REMOVED***"kern.securelevel", []_C_int***REMOVED***1, 9***REMOVED******REMOVED***,
	***REMOVED***"kern.seminfo", []_C_int***REMOVED***1, 61***REMOVED******REMOVED***,
	***REMOVED***"kern.shminfo", []_C_int***REMOVED***1, 62***REMOVED******REMOVED***,
	***REMOVED***"kern.somaxconn", []_C_int***REMOVED***1, 28***REMOVED******REMOVED***,
	***REMOVED***"kern.sominconn", []_C_int***REMOVED***1, 29***REMOVED******REMOVED***,
	***REMOVED***"kern.splassert", []_C_int***REMOVED***1, 54***REMOVED******REMOVED***,
	***REMOVED***"kern.stackgap_random", []_C_int***REMOVED***1, 50***REMOVED******REMOVED***,
	***REMOVED***"kern.sysvipc_info", []_C_int***REMOVED***1, 51***REMOVED******REMOVED***,
	***REMOVED***"kern.sysvmsg", []_C_int***REMOVED***1, 34***REMOVED******REMOVED***,
	***REMOVED***"kern.sysvsem", []_C_int***REMOVED***1, 35***REMOVED******REMOVED***,
	***REMOVED***"kern.sysvshm", []_C_int***REMOVED***1, 36***REMOVED******REMOVED***,
	***REMOVED***"kern.timecounter.choice", []_C_int***REMOVED***1, 69, 4***REMOVED******REMOVED***,
	***REMOVED***"kern.timecounter.hardware", []_C_int***REMOVED***1, 69, 3***REMOVED******REMOVED***,
	***REMOVED***"kern.timecounter.tick", []_C_int***REMOVED***1, 69, 1***REMOVED******REMOVED***,
	***REMOVED***"kern.timecounter.timestepwarnings", []_C_int***REMOVED***1, 69, 2***REMOVED******REMOVED***,
	***REMOVED***"kern.tty.tk_cancc", []_C_int***REMOVED***1, 44, 4***REMOVED******REMOVED***,
	***REMOVED***"kern.tty.tk_nin", []_C_int***REMOVED***1, 44, 1***REMOVED******REMOVED***,
	***REMOVED***"kern.tty.tk_nout", []_C_int***REMOVED***1, 44, 2***REMOVED******REMOVED***,
	***REMOVED***"kern.tty.tk_rawcc", []_C_int***REMOVED***1, 44, 3***REMOVED******REMOVED***,
	***REMOVED***"kern.tty.ttyinfo", []_C_int***REMOVED***1, 44, 5***REMOVED******REMOVED***,
	***REMOVED***"kern.ttycount", []_C_int***REMOVED***1, 57***REMOVED******REMOVED***,
	***REMOVED***"kern.version", []_C_int***REMOVED***1, 4***REMOVED******REMOVED***,
	***REMOVED***"kern.watchdog.auto", []_C_int***REMOVED***1, 64, 2***REMOVED******REMOVED***,
	***REMOVED***"kern.watchdog.period", []_C_int***REMOVED***1, 64, 1***REMOVED******REMOVED***,
	***REMOVED***"kern.witnesswatch", []_C_int***REMOVED***1, 53***REMOVED******REMOVED***,
	***REMOVED***"kern.wxabort", []_C_int***REMOVED***1, 74***REMOVED******REMOVED***,
	***REMOVED***"net.bpf.bufsize", []_C_int***REMOVED***4, 31, 1***REMOVED******REMOVED***,
	***REMOVED***"net.bpf.maxbufsize", []_C_int***REMOVED***4, 31, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ah.enable", []_C_int***REMOVED***4, 2, 51, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ah.stats", []_C_int***REMOVED***4, 2, 51, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.carp.allow", []_C_int***REMOVED***4, 2, 112, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.carp.log", []_C_int***REMOVED***4, 2, 112, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet.carp.preempt", []_C_int***REMOVED***4, 2, 112, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.carp.stats", []_C_int***REMOVED***4, 2, 112, 4***REMOVED******REMOVED***,
	***REMOVED***"net.inet.divert.recvspace", []_C_int***REMOVED***4, 2, 258, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.divert.sendspace", []_C_int***REMOVED***4, 2, 258, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.divert.stats", []_C_int***REMOVED***4, 2, 258, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet.esp.enable", []_C_int***REMOVED***4, 2, 50, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.esp.stats", []_C_int***REMOVED***4, 2, 50, 4***REMOVED******REMOVED***,
	***REMOVED***"net.inet.esp.udpencap", []_C_int***REMOVED***4, 2, 50, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.esp.udpencap_port", []_C_int***REMOVED***4, 2, 50, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet.etherip.allow", []_C_int***REMOVED***4, 2, 97, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.etherip.stats", []_C_int***REMOVED***4, 2, 97, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.gre.allow", []_C_int***REMOVED***4, 2, 47, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.gre.wccp", []_C_int***REMOVED***4, 2, 47, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.icmp.bmcastecho", []_C_int***REMOVED***4, 2, 1, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.icmp.errppslimit", []_C_int***REMOVED***4, 2, 1, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet.icmp.maskrepl", []_C_int***REMOVED***4, 2, 1, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.icmp.rediraccept", []_C_int***REMOVED***4, 2, 1, 4***REMOVED******REMOVED***,
	***REMOVED***"net.inet.icmp.redirtimeout", []_C_int***REMOVED***4, 2, 1, 5***REMOVED******REMOVED***,
	***REMOVED***"net.inet.icmp.stats", []_C_int***REMOVED***4, 2, 1, 7***REMOVED******REMOVED***,
	***REMOVED***"net.inet.icmp.tstamprepl", []_C_int***REMOVED***4, 2, 1, 6***REMOVED******REMOVED***,
	***REMOVED***"net.inet.igmp.stats", []_C_int***REMOVED***4, 2, 2, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.arpdown", []_C_int***REMOVED***4, 2, 0, 40***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.arpqueued", []_C_int***REMOVED***4, 2, 0, 36***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.arptimeout", []_C_int***REMOVED***4, 2, 0, 39***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.encdebug", []_C_int***REMOVED***4, 2, 0, 12***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.forwarding", []_C_int***REMOVED***4, 2, 0, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.ifq.congestion", []_C_int***REMOVED***4, 2, 0, 30, 4***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.ifq.drops", []_C_int***REMOVED***4, 2, 0, 30, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.ifq.len", []_C_int***REMOVED***4, 2, 0, 30, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.ifq.maxlen", []_C_int***REMOVED***4, 2, 0, 30, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.maxqueue", []_C_int***REMOVED***4, 2, 0, 11***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.mforwarding", []_C_int***REMOVED***4, 2, 0, 31***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.mrtmfc", []_C_int***REMOVED***4, 2, 0, 37***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.mrtproto", []_C_int***REMOVED***4, 2, 0, 34***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.mrtstats", []_C_int***REMOVED***4, 2, 0, 35***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.mrtvif", []_C_int***REMOVED***4, 2, 0, 38***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.mtu", []_C_int***REMOVED***4, 2, 0, 4***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.mtudisc", []_C_int***REMOVED***4, 2, 0, 27***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.mtudisctimeout", []_C_int***REMOVED***4, 2, 0, 28***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.multipath", []_C_int***REMOVED***4, 2, 0, 32***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.portfirst", []_C_int***REMOVED***4, 2, 0, 7***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.porthifirst", []_C_int***REMOVED***4, 2, 0, 9***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.porthilast", []_C_int***REMOVED***4, 2, 0, 10***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.portlast", []_C_int***REMOVED***4, 2, 0, 8***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.redirect", []_C_int***REMOVED***4, 2, 0, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.sourceroute", []_C_int***REMOVED***4, 2, 0, 5***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.stats", []_C_int***REMOVED***4, 2, 0, 33***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ip.ttl", []_C_int***REMOVED***4, 2, 0, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ipcomp.enable", []_C_int***REMOVED***4, 2, 108, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ipcomp.stats", []_C_int***REMOVED***4, 2, 108, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ipip.allow", []_C_int***REMOVED***4, 2, 4, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.ipip.stats", []_C_int***REMOVED***4, 2, 4, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.mobileip.allow", []_C_int***REMOVED***4, 2, 55, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.pfsync.stats", []_C_int***REMOVED***4, 2, 240, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.ackonpush", []_C_int***REMOVED***4, 2, 6, 13***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.always_keepalive", []_C_int***REMOVED***4, 2, 6, 22***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.baddynamic", []_C_int***REMOVED***4, 2, 6, 6***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.drop", []_C_int***REMOVED***4, 2, 6, 19***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.ecn", []_C_int***REMOVED***4, 2, 6, 14***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.ident", []_C_int***REMOVED***4, 2, 6, 9***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.keepidle", []_C_int***REMOVED***4, 2, 6, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.keepinittime", []_C_int***REMOVED***4, 2, 6, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.keepintvl", []_C_int***REMOVED***4, 2, 6, 4***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.mssdflt", []_C_int***REMOVED***4, 2, 6, 11***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.reasslimit", []_C_int***REMOVED***4, 2, 6, 18***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.rfc1323", []_C_int***REMOVED***4, 2, 6, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.rfc3390", []_C_int***REMOVED***4, 2, 6, 17***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.rootonly", []_C_int***REMOVED***4, 2, 6, 24***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.rstppslimit", []_C_int***REMOVED***4, 2, 6, 12***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.sack", []_C_int***REMOVED***4, 2, 6, 10***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.sackholelimit", []_C_int***REMOVED***4, 2, 6, 20***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.slowhz", []_C_int***REMOVED***4, 2, 6, 5***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.stats", []_C_int***REMOVED***4, 2, 6, 21***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.synbucketlimit", []_C_int***REMOVED***4, 2, 6, 16***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.syncachelimit", []_C_int***REMOVED***4, 2, 6, 15***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.synhashsize", []_C_int***REMOVED***4, 2, 6, 25***REMOVED******REMOVED***,
	***REMOVED***"net.inet.tcp.synuselimit", []_C_int***REMOVED***4, 2, 6, 23***REMOVED******REMOVED***,
	***REMOVED***"net.inet.udp.baddynamic", []_C_int***REMOVED***4, 2, 17, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet.udp.checksum", []_C_int***REMOVED***4, 2, 17, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet.udp.recvspace", []_C_int***REMOVED***4, 2, 17, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet.udp.rootonly", []_C_int***REMOVED***4, 2, 17, 6***REMOVED******REMOVED***,
	***REMOVED***"net.inet.udp.sendspace", []_C_int***REMOVED***4, 2, 17, 4***REMOVED******REMOVED***,
	***REMOVED***"net.inet.udp.stats", []_C_int***REMOVED***4, 2, 17, 5***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.divert.recvspace", []_C_int***REMOVED***4, 24, 86, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.divert.sendspace", []_C_int***REMOVED***4, 24, 86, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.divert.stats", []_C_int***REMOVED***4, 24, 86, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.icmp6.errppslimit", []_C_int***REMOVED***4, 24, 30, 14***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.icmp6.mtudisc_hiwat", []_C_int***REMOVED***4, 24, 30, 16***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.icmp6.mtudisc_lowat", []_C_int***REMOVED***4, 24, 30, 17***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.icmp6.nd6_debug", []_C_int***REMOVED***4, 24, 30, 18***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.icmp6.nd6_delay", []_C_int***REMOVED***4, 24, 30, 8***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.icmp6.nd6_maxnudhint", []_C_int***REMOVED***4, 24, 30, 15***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.icmp6.nd6_mmaxtries", []_C_int***REMOVED***4, 24, 30, 10***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.icmp6.nd6_umaxtries", []_C_int***REMOVED***4, 24, 30, 9***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.icmp6.redirtimeout", []_C_int***REMOVED***4, 24, 30, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.auto_flowlabel", []_C_int***REMOVED***4, 24, 17, 17***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.dad_count", []_C_int***REMOVED***4, 24, 17, 16***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.dad_pending", []_C_int***REMOVED***4, 24, 17, 49***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.defmcasthlim", []_C_int***REMOVED***4, 24, 17, 18***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.forwarding", []_C_int***REMOVED***4, 24, 17, 1***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.forwsrcrt", []_C_int***REMOVED***4, 24, 17, 5***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.hdrnestlimit", []_C_int***REMOVED***4, 24, 17, 15***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.hlim", []_C_int***REMOVED***4, 24, 17, 3***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.log_interval", []_C_int***REMOVED***4, 24, 17, 14***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.maxdynroutes", []_C_int***REMOVED***4, 24, 17, 48***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.maxfragpackets", []_C_int***REMOVED***4, 24, 17, 9***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.maxfrags", []_C_int***REMOVED***4, 24, 17, 41***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.mforwarding", []_C_int***REMOVED***4, 24, 17, 42***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.mrtmfc", []_C_int***REMOVED***4, 24, 17, 53***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.mrtmif", []_C_int***REMOVED***4, 24, 17, 52***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.mrtproto", []_C_int***REMOVED***4, 24, 17, 8***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.mtudisctimeout", []_C_int***REMOVED***4, 24, 17, 50***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.multicast_mtudisc", []_C_int***REMOVED***4, 24, 17, 44***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.multipath", []_C_int***REMOVED***4, 24, 17, 43***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.neighborgcthresh", []_C_int***REMOVED***4, 24, 17, 45***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.redirect", []_C_int***REMOVED***4, 24, 17, 2***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.soiikey", []_C_int***REMOVED***4, 24, 17, 54***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.sourcecheck", []_C_int***REMOVED***4, 24, 17, 10***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.sourcecheck_logint", []_C_int***REMOVED***4, 24, 17, 11***REMOVED******REMOVED***,
	***REMOVED***"net.inet6.ip6.use_deprecated", []_C_int***REMOVED***4, 24, 17, 21***REMOVED******REMOVED***,
	***REMOVED***"net.key.sadb_dump", []_C_int***REMOVED***4, 30, 1***REMOVED******REMOVED***,
	***REMOVED***"net.key.spd_dump", []_C_int***REMOVED***4, 30, 2***REMOVED******REMOVED***,
	***REMOVED***"net.mpls.ifq.congestion", []_C_int***REMOVED***4, 33, 3, 4***REMOVED******REMOVED***,
	***REMOVED***"net.mpls.ifq.drops", []_C_int***REMOVED***4, 33, 3, 3***REMOVED******REMOVED***,
	***REMOVED***"net.mpls.ifq.len", []_C_int***REMOVED***4, 33, 3, 1***REMOVED******REMOVED***,
	***REMOVED***"net.mpls.ifq.maxlen", []_C_int***REMOVED***4, 33, 3, 2***REMOVED******REMOVED***,
	***REMOVED***"net.mpls.mapttl_ip", []_C_int***REMOVED***4, 33, 5***REMOVED******REMOVED***,
	***REMOVED***"net.mpls.mapttl_ip6", []_C_int***REMOVED***4, 33, 6***REMOVED******REMOVED***,
	***REMOVED***"net.mpls.maxloop_inkernel", []_C_int***REMOVED***4, 33, 4***REMOVED******REMOVED***,
	***REMOVED***"net.mpls.ttl", []_C_int***REMOVED***4, 33, 2***REMOVED******REMOVED***,
	***REMOVED***"net.pflow.stats", []_C_int***REMOVED***4, 34, 1***REMOVED******REMOVED***,
	***REMOVED***"net.pipex.enable", []_C_int***REMOVED***4, 35, 1***REMOVED******REMOVED***,
	***REMOVED***"vm.anonmin", []_C_int***REMOVED***2, 7***REMOVED******REMOVED***,
	***REMOVED***"vm.loadavg", []_C_int***REMOVED***2, 2***REMOVED******REMOVED***,
	***REMOVED***"vm.malloc_conf", []_C_int***REMOVED***2, 12***REMOVED******REMOVED***,
	***REMOVED***"vm.maxslp", []_C_int***REMOVED***2, 10***REMOVED******REMOVED***,
	***REMOVED***"vm.nkmempages", []_C_int***REMOVED***2, 6***REMOVED******REMOVED***,
	***REMOVED***"vm.psstrings", []_C_int***REMOVED***2, 3***REMOVED******REMOVED***,
	***REMOVED***"vm.swapencrypt.enable", []_C_int***REMOVED***2, 5, 0***REMOVED******REMOVED***,
	***REMOVED***"vm.swapencrypt.keyscreated", []_C_int***REMOVED***2, 5, 1***REMOVED******REMOVED***,
	***REMOVED***"vm.swapencrypt.keysdeleted", []_C_int***REMOVED***2, 5, 2***REMOVED******REMOVED***,
	***REMOVED***"vm.uspace", []_C_int***REMOVED***2, 11***REMOVED******REMOVED***,
	***REMOVED***"vm.uvmexp", []_C_int***REMOVED***2, 4***REMOVED******REMOVED***,
	***REMOVED***"vm.vmmeter", []_C_int***REMOVED***2, 1***REMOVED******REMOVED***,
	***REMOVED***"vm.vnodemin", []_C_int***REMOVED***2, 9***REMOVED******REMOVED***,
	***REMOVED***"vm.vtextmin", []_C_int***REMOVED***2, 8***REMOVED******REMOVED***,
***REMOVED***
