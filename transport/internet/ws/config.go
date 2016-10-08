package ws

import (
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"

	"github.com/golang/protobuf/proto"
)

func init() {
	internet.RegisterNetworkConfigCreator(v2net.Network_WebSocket, func() proto.Message {
		return new(Config)
	})
}
