package internet

import (
	"errors"
	"net"

	v2net "v2ray.com/core/common/net"
)

var (
	ErrUnsupportedStreamType = errors.New("Unsupported stream type.")
)

type DialerOptions struct {
	Stream *StreamConfig
}

type Dialer func(src v2net.Address, dest v2net.Destination, options DialerOptions) (Connection, error)

var (
	TCPDialer    Dialer
	KCPDialer    Dialer
	RawTCPDialer Dialer
	UDPDialer    Dialer
	WSDialer     Dialer
)

func Dial(src v2net.Address, dest v2net.Destination, settings *StreamConfig) (Connection, error) {

	var connection Connection
	var err error
	dialerOptions := DialerOptions{
		Stream: settings,
	}
	if dest.Network == v2net.Network_TCP {
		switch settings.Network {
		case v2net.Network_TCP:
			connection, err = TCPDialer(src, dest, dialerOptions)
		case v2net.Network_KCP:
			connection, err = KCPDialer(src, dest, dialerOptions)
		case v2net.Network_WebSocket:
			connection, err = WSDialer(src, dest, dialerOptions)

			// This check has to be the last one.
		case v2net.Network_RawTCP:
			connection, err = RawTCPDialer(src, dest, dialerOptions)
		default:
			return nil, ErrUnsupportedStreamType
		}
		if err != nil {
			return nil, err
		}

		return connection, nil
	}

	return UDPDialer(src, dest, dialerOptions)
}

func DialToDest(src v2net.Address, dest v2net.Destination) (net.Conn, error) {
	return effectiveSystemDialer.Dial(src, dest)
}
