//
// @project IPRangeTree 2016 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2019
//

package iprangetree

import (
	"net"
)

// IPItemAccessor data accessor interface
type IPItemAccessor interface {
	GetData() interface{}
	GetStartIP() net.IP
	GetEndIP() net.IP
}
