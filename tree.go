//
// @project IPRangeTree 2016 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2019
//

package iprangetree

import (
	"errors"
	"net"
	"strings"

	"github.com/google/btree"
)

// Errors list
var (
	ErrInvalidItem        = errors.New("[iprangetree] invalid IP range item")
	ErrItemRaplaced       = errors.New("[iprangetree] item was replaced")
	ErrInvalidRangeFormat = errors.New("[iprangetree] invalid range format")
	ErrInvalidRangeValues = errors.New("[iprangetree] start IP is less then end IP")
)

// IPTree base
type IPTree struct {
	mixed bool // If in IPv6 could conatains ranges with IPv4
	ipV4  *btree.BTree
	ipV6  *btree.BTree
}

// New tree object
func New(degree int) *IPTree {
	return &IPTree{
		mixed: false,
		ipV4:  btree.New(degree),
		ipV6:  btree.New(degree),
	}
}

// AddRangeByString data by IP range string or single IP
// Example:
// 	db.AddRangeByString("127.0.0.1", data)
// 	db.AddRangeByString("127.0.0.1-127.0.0.2", data)
func (t *IPTree) AddRangeByString(ipRange string, val interface{}) (err error) {
	ipRange = strings.Trim(ipRange, " \t\n-–+")
	switch {
	case strings.Contains(ipRange, "-"):
		if arr := strings.Split(ipRange, "-"); len(arr) == 2 {
			return t.AddRange(net.ParseIP(arr[0]), net.ParseIP(arr[1]), val)
		}
		err = ErrInvalidRangeFormat
	case strings.Contains(ipRange, "/"):
		var ip net.IP
		var inet *net.IPNet
		if ip, inet, err = net.ParseCIDR(ipRange); err == nil {
			return t.AddRange(ip, lastIP(inet.IP, inet.Mask), val)
		}
	default:
		ip := net.ParseIP(ipRange)
		return t.AddRange(ip, ip, val)
	}
	return
}

// AddRange IPs vith value
func (t *IPTree) AddRange(ip1, ip2 net.IP, val interface{}) (err error) {
	var (
		startIP = ip2fix(ip1)
		endIP   = ip2fix(ip2)
		old     interface{}
	)

	if startIP.IsEmpty() && endIP.IsEmpty() {
		return ErrInvalidItem
	}

	if startIP.IsEmpty() {
		startIP = endIP
	} else if endIP.IsEmpty() {
		endIP = startIP
	}

	if startIP.Compare(&endIP) > 0 {
		return ErrInvalidRangeValues
	}

	if startIP.IsIPv4() && endIP.IsIPv4() {
		old = t.ipV4.ReplaceOrInsert(&IPItemV4{
			StartIP: startIP.IPv4(),
			EndIP:   endIP.IPv4(),
			Data:    val,
		})
	} else {
		if startIP.IsIPv4() || endIP.IsIPv4() {
			t.mixed = true
		}
		old = t.ipV6.ReplaceOrInsert(&IPItemFix{
			StartIP: startIP,
			EndIP:   endIP,
			Data:    val,
		})
	}
	if old != nil {
		err = ErrItemRaplaced
	}
	return
}

// LookupByString parse the IP value and search data in the IP tree
func (t *IPTree) LookupByString(ip string) IPItemAccessor {
	return t.Lookup(net.ParseIP(ip))
}

// Lookup to search the IP value in the IP tree
//go:notinheap
func (t *IPTree) Lookup(ip net.IP) (response IPItemAccessor) {
	ipFix := ip2fix(ip)

	if ipFix.IsIPv4() {
		ipv4 := ipFix.IPv4()
		t.ipV4.AscendGreaterOrEqual(ipv4, func(item btree.Item) bool {
			if response != nil {
				return false
			}

			it := item.(*IPItemV4)
			switch ipv4.Compare(it.StartIP) {
			case 1:
				if ipv4.Compare(it.EndIP) <= 0 {
					response = it
					return false
				}
			case 0:
				response = it
				return false
			case -1:
				return false
			}
			return true
		})
	}

	if response == nil && (!ipFix.IsIPv4() || t.mixed) {
		t.ipV6.AscendGreaterOrEqual(&ipFix, func(item btree.Item) bool {
			if response != nil {
				return false
			}

			it := item.(*IPItemFix)
			switch ipFix.Compare(&it.StartIP) {
			case 1:
				if ipFix.Compare(&it.EndIP) <= 0 {
					response = it
					return false
				}
			case 0:
				response = it
				return false
			case -1:
				return false
			}
			return true
		})
	}
	return
}
