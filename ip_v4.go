package iprangetree

import (
	"encoding/binary"
	"math"
	"net"

	"github.com/google/btree"
)

const undefinedIPv4 ipV4 = math.MaxUint32

var byteOrder binary.ByteOrder = binary.BigEndian

type ipV4 uint32

func toIPv4(ip net.IP) ipV4 {
	if ip = ip.To4(); ip == nil {
		return undefinedIPv4
	}
	return ipV4(byteOrder.Uint32([]byte(ip)))
}

func (ip ipV4) IsIPv4() bool {
	return true
}

func (ip ipV4) Compare(ip2 ipV4) int {
	if ip < ip2 {
		return -1
	} else if ip > ip2 {
		return 1
	}
	return 0
}

func (ip ipV4) IP() net.IP {
	data := make([]byte, 4)
	byteOrder.PutUint32(data, uint32(ip))
	return net.IP(data)
}

// Less camparing for btree
func (ip ipV4) Less(then btree.Item) bool {
	switch v := then.(type) {
	case *IPItemFix:
		return !v.StartIP.IsIPv4() || ip.Compare(v.StartIP.IPv4()) < 0
	case *IPItemV4:
		return ip.Compare(v.StartIP) < 0
	case ipV4:
		return ip.Compare(v) < 0
	case *ipFix:
		return !v.IsIPv4() || ip.Compare(v.IPv4()) < 0
	}
	return false
}
