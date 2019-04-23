//
// @project IPRangeTree 2016 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2019
//

package iprangetree

import (
	"net"
)

// We awaiting IPv4 as size 4 and IPv6 as size 16
//go:notinheap
func compare(ip1, ip2 net.IP) int {
	if len(ip1) > len(ip2) {
		return 1
	} else if len(ip1) < len(ip2) {
		return -1
	}

	for i := 0; i < len(ip1); i++ {
		if ip1[i] < ip2[i] {
			return -1
		} else if ip1[i] > ip2[i] {
			return 1
		}
	}
	return 0
}

//go:notinheap
func lastIP(ip net.IP, mask net.IPMask) net.IP {
	var (
		n   = len(mask)
		j   = len(ip) - n
		out = make(net.IP, n)
	)

	for i := 0; i < n; i++ {
		out[i] = ip[j] | ^mask[i]
		j++
	}
	return out
}
