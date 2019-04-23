//
// @project IPRangeTree 2016 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2019
//

package iprangetree

import (
	"net"

	"github.com/google/btree"
)

// IPItemFix IP range
type IPItemFix struct {
	StartIP ipFix
	EndIP   ipFix
	Data    interface{}
}

// GetData object
func (i *IPItemFix) GetData() interface{} {
	return i.Data
}

// Less camparing for btree
func (i *IPItemFix) Less(then btree.Item) bool {
	switch ip := then.(type) {
	case *IPItemFix:
		return i.StartIP.Compare(&ip.StartIP) < 0
	case *ipFix:
		return i.EndIP.Compare(ip) < 0
	}
	return false
}

// Compare with the second item
func (i *IPItemFix) Compare(it interface{}) int {
	switch ip := it.(type) {
	case *IPItemFix:
		return i.StartIP.Compare(&ip.StartIP)
	case *ipFix:
		return i.EndIP.Compare(ip)
	}
	return 0
}

// Has IP in range
func (i *IPItemFix) Has(ip net.IP) bool {
	return compare(i.StartIP.IP(), ip) <= 0 && compare(i.EndIP.IP(), ip) >= 0
}

// GetStartIP net object
func (i *IPItemFix) GetStartIP() net.IP {
	return i.StartIP.IP()
}

// GetEndIP net object
func (i *IPItemFix) GetEndIP() net.IP {
	return i.EndIP.IP()
}

var _ IPItemAccessor = (*IPItemFix)(nil)
