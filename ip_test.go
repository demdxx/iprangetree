package iprangetree

import (
	"net"
	"sync/atomic"
	"testing"
)

func Test_IPCompareV4(t *testing.T) {
	var tests = []struct {
		ip1    ipV4
		ip2    ipV4
		result int
	}{
		{
			ip1:    toIPv4(net.ParseIP("89.18.12.13")),
			ip2:    toIPv4(net.ParseIP("89.20.12.13")),
			result: -1,
		},
		{
			ip1:    toIPv4(net.ParseIP("89.18.12.17")),
			ip2:    toIPv4(net.ParseIP("89.18.12.13")),
			result: 1,
		},
		{
			ip1:    toIPv4(net.ParseIP("89.18.12.17")),
			ip2:    toIPv4(net.ParseIP("89.18.12.17")),
			result: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.ip1.IP().String()+" ? "+test.ip2.IP().String(), func(t *testing.T) {
			if result := test.ip1.Compare(test.ip2); result != test.result {
				t.Errorf("Incorrect comparison of %s & %s got %d but have to be %d",
					test.ip1.IP().String(), test.ip2.IP().String(), result, test.result)
			}
		})
	}
}

func Test_IPCompareFix(t *testing.T) {
	var tests = []struct {
		ip1    ipFix
		ip2    ipFix
		result int
		ip1v4  bool
		ip2v4  bool
	}{
		{
			ip1:    ip2fix(net.ParseIP("89.18.12.13")),
			ip1v4:  true,
			ip2:    ip2fix(net.ParseIP("89.20.12.13")),
			ip2v4:  true,
			result: -1,
		},
		{
			ip1:    ip2fix(net.ParseIP("89.18.12.17")),
			ip1v4:  true,
			ip2:    ip2fix(net.ParseIP("89.18.12.13")),
			ip2v4:  true,
			result: 1,
		},
		{
			ip1:    ip2fix(net.ParseIP("89.18.12.17")),
			ip1v4:  true,
			ip2:    ip2fix(net.ParseIP("89.18.12.17")),
			ip2v4:  true,
			result: 0,
		},
		{
			ip1:    ip2fix(net.ParseIP("2001:0db8:0000:0000:0000:ff00:0042:8329")),
			ip1v4:  false,
			ip2:    ip2fix(net.ParseIP("2001:0db8:0000:0000:0000:ff00:0042:8329")),
			ip2v4:  false,
			result: 0,
		},
		{
			ip1:    ip2fix(net.ParseIP("2001:0db8:0000:0000:0001:ff00:0042:8329")),
			ip1v4:  false,
			ip2:    ip2fix(net.ParseIP("2001:0db8:0000:0000:0000:ff00:0042:8329")),
			ip2v4:  false,
			result: 1,
		},
		{
			ip1:    ip2fix(net.ParseIP("2001:0db8:0000:0000:0000:ff00:0042:8329")),
			ip1v4:  false,
			ip2:    ip2fix(net.ParseIP("2001:0db8:0000:0000:0001:ff00:0042:8329")),
			ip2v4:  false,
			result: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.ip1.IP().String()+" ? "+test.ip2.IP().String(), func(t *testing.T) {
			if test.ip1.IsIPv4() != test.ip1v4 {
				t.Errorf("IP1 must be V4: %v", test.ip1v4)
			}
			if test.ip2.IsIPv4() != test.ip2v4 {
				t.Errorf("IP2 must be V4: %v", test.ip2v4)
			}
			if result := test.ip1.Compare(&test.ip2); result != test.result {
				t.Errorf("Incorrect comparison of %s & %s got %d but have to be %d",
					test.ip1.IP().String(), test.ip2.IP().String(), result, test.result)
			}
		})
	}
}

func Benchmark_IPCompare(b *testing.B) {
	var ips = []ipFix{
		ip2fix(net.ParseIP("127.0.0.1")),
		ip2fix(net.ParseIP("189.1.3.4")),
		ip2fix(net.ParseIP("111.2.3.4")),
	}
	var (
		idx   = int32(0)
		iplen = int32(len(ips))
	)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ip1 := &ips[atomic.LoadInt32(&idx)%iplen]
			atomic.AddInt32(&idx, 1)
			ip2 := &ips[atomic.LoadInt32(&idx)%iplen]
			atomic.AddInt32(&idx, 1)

			_ = ip1.Compare(ip2)
		}
	})
}

func Benchmark_IPv4Compare(b *testing.B) {
	var ips = []ipV4{
		toIPv4(net.ParseIP("127.0.0.1")),
		toIPv4(net.ParseIP("189.1.3.4")),
		toIPv4(net.ParseIP("111.2.3.4")),
	}
	var (
		idx   = int32(0)
		iplen = int32(len(ips))
	)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ip1 := ips[atomic.LoadInt32(&idx)%iplen]
			atomic.AddInt32(&idx, 1)
			ip2 := ips[atomic.LoadInt32(&idx)%iplen]
			atomic.AddInt32(&idx, 1)

			_ = ip1.Compare(ip2)
		}
	})
}

func Benchmark_IPCompare2(b *testing.B) {
	var ips = []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("189.1.3.4"),
		net.ParseIP("111.2.3.4"),
	}
	var (
		idx   = int32(0)
		iplen = int32(len(ips))
	)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ip1 := ips[atomic.LoadInt32(&idx)%iplen]
			atomic.AddInt32(&idx, 1)
			ip2 := ips[atomic.LoadInt32(&idx)%iplen]
			atomic.AddInt32(&idx, 1)

			_ = compare(ip1, ip2)
		}
	})
}
