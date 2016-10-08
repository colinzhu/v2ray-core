package udp

import (
	"net"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

type Connection struct {
	net.UDPConn
}

func (this *Connection) Reusable() bool {
	return false
}

func (this *Connection) SetReusable(b bool) {}

func init() {
	internet.UDPDialer = func(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
		conn, err := internet.DialToDest(src, dest)
		if err != nil {
			return nil, err
		}
		// TODO: handle dialer options
		return &Connection{
			UDPConn: *(conn.(*net.UDPConn)),
		}, nil
	}
}
