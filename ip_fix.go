//
// @project IPRangeTree 2016 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2019
//

package iprangetree

import (
	"net"

	"github.com/google/btree"
)

// Fixed size IP data type
type ipFix [net.IPv6len + 1]byte

func ip2fix(ip net.IP) (nIP ipFix) {
	var nIPCursor = net.IPv6len - 1
	var isIPv4 = true
	for i := len(ip) - 1; i >= 0; i-- {
		nIP[nIPCursor] = ip[i]
		if isIPv4 {
			switch {
			case nIPCursor == 11 || nIPCursor == 10:
				if nIP[nIPCursor] != 0xff {
					isIPv4 = false
				}
			case nIPCursor < 10:
				if nIP[nIPCursor] != 0 {
					isIPv4 = false
				}
			}
		}
		nIPCursor--
	}
	if isIPv4 {
		nIP[net.IPv6len] = 1
	} else {
		nIP[net.IPv6len] = 2
	}
	return
}

func (ip ipFix) IsIPv4() bool {
	return ip[net.IPv6len] == 1
}

func (ip ipFix) IsEmpty() bool {
	return ip[net.IPv6len] == 0
}

//go:notinheap
func (ip *ipFix) Compare(ip2 *ipFix) (result int) {
	i := 0
	if v1, v2 := ip.IsIPv4(), ip2.IsIPv4(); !v1 && v2 {
		return 1
	} else if v1 && !v2 {
		return -1
	} else if v1 {
		i = 12
	}
	for ; i < net.IPv6len; i++ {
		if ip[i] < ip2[i] {
			return -1
		} else if ip[i] > ip2[i] {
			return 1
		}
	}
	return 0
}

func (ip ipFix) IP() net.IP {
	if ip.IsIPv4() {
		return net.IP(ip[12:net.IPv6len])
	}
	return net.IP(ip[:net.IPv6len])
}

// IPv4 simple format
func (ip ipFix) IPv4() ipV4 {
	if ip.IsIPv4() {
		return toIPv4(ip[12:net.IPv6len])
	}
	return undefinedIPv4
}

// Less camparing for btree
//go:notinheap
func (ip *ipFix) Less(then btree.Item) bool {
	switch v := then.(type) {
	case *IPItemFix:
		return ip.Compare(&v.StartIP) < 0
	case *IPItemV4:
		return ip.IsIPv4() && ip.IPv4().Compare(v.StartIP) < 0
	case ipV4:
		return ip.IsIPv4() && ip.IPv4().Compare(v) < 0
	case *ipFix:
		return ip.Compare(v) < 0
	}
	return false
}
