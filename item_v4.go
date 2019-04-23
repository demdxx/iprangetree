//
// @project IPRangeTree 2016 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2019
//

package iprangetree

import (
	"net"

	"github.com/google/btree"
)

// IPItemV4 IP range
type IPItemV4 struct {
	StartIP ipV4
	EndIP   ipV4
	Data    interface{}
}

// GetData object
func (i *IPItemV4) GetData() interface{} {
	return i.Data
}

// Less camparing for btree
func (i *IPItemV4) Less(then btree.Item) bool {
	switch ip := then.(type) {
	case *IPItemV4:
		return i.StartIP.Compare(ip.StartIP) < 0
	case *IPItemFix:
		return !ip.StartIP.IsIPv4() || i.StartIP.Compare(ip.StartIP.IPv4()) < 0
	case ipV4:
		return i.EndIP.Compare(ip) < 0
	case *ipFix:
		return !ip.IsIPv4() || i.EndIP.Compare(ip.IPv4()) < 0
	}
	return false
}

// Compare with the second item
func (i *IPItemV4) Compare(it interface{}) int {
	switch ip := it.(type) {
	case *IPItemV4:
		return i.StartIP.Compare(ip.StartIP)
	case *IPItemFix:
		if !ip.StartIP.IsIPv4() {
			return -1
		}
		return i.StartIP.Compare(ip.StartIP.IPv4())
	case ipV4:
		return i.EndIP.Compare(ip)
	case *ipFix:
		if !ip.IsIPv4() {
			return -1
		}
		return i.EndIP.Compare(ip.IPv4())
	}
	return 0
}

// Has IP in range
func (i *IPItemV4) Has(ip net.IP) bool {
	return compare(i.StartIP.IP(), ip) <= 0 && compare(i.EndIP.IP(), ip) >= 0
}

// GetStartIP net object
func (i *IPItemV4) GetStartIP() net.IP {
	return i.StartIP.IP()
}

// GetEndIP net object
func (i *IPItemV4) GetEndIP() net.IP {
	return i.EndIP.IP()
}

var _ IPItemAccessor = (*IPItemV4)(nil)
