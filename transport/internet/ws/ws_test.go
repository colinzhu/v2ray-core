package ws_test

import (
	"io/ioutil"
	"testing"
	"time"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
	. "v2ray.com/core/transport/internet/ws"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

func MarshalSettings(config *Config, assert *assert.Assert) *any.Any {
	anySettings, err := ptypes.MarshalAny(config)
	assert.Error(err).IsNil()

	return anySettings
}

func Test_Connect_ws(t *testing.T) {
	assert := assert.On(t)

	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 80), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path: "",
					}, assert),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss(t *testing.T) {
	assert := assert.On(t)
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path: "",
					}, assert),
				},
			},
			SecurityType: internet.SecurityType_TLS,
		},
	})
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_1_nil(t *testing.T) {
	assert := assert.On(t)
	conn, err := Dial(nil, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path: "",
					}, assert),
				},
			},
			SecurityType: internet.SecurityType_TLS,
		},
	})
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_ws_guess(t *testing.T) {
	assert := assert.On(t)
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 80), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path: "",
					}, assert),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_guess(t *testing.T) {
	assert := assert.On(t)
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path: "",
					}, assert),
				},
			},
			SecurityType: internet.SecurityType_TLS,
		},
	})
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_guess_fail(t *testing.T) {
	assert := assert.On(t)
	_, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("static.kkdev.org"), 443), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path: "",
					}, assert),
				},
			},
			SecurityType: internet.SecurityType_TLS,
		},
	})
	assert.Error(err).IsNotNil()
}

func Test_Connect_wss_guess_reuse(t *testing.T) {
	assert := assert.On(t)
	i := 3
	for i != 0 {
		conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443), internet.DialerOptions{
			Stream: &internet.StreamConfig{
				Network: v2net.Network_WebSocket,
				NetworkSettings: []*internet.NetworkSettings{
					{
						Network: v2net.Network_WebSocket,
						Settings: MarshalSettings(&Config{
							Path:            "",
							ConnectionReuse: true,
						}, assert),
					},
				},
				SecurityType: internet.SecurityType_TLS,
			},
		})
		assert.Error(err).IsNil()
		conn.Write([]byte("echo"))
		s := make(chan int)
		go func() {
			buf := make([]byte, 4)
			conn.Read(buf)
			str := string(buf)
			if str != "echo" {
				assert.Fail("Data mismatch")
			}
			s <- 0
		}()
		<-s
		if i == 0 {
			conn.SetDeadline(time.Now())
			conn.SetReadDeadline(time.Now())
			conn.SetWriteDeadline(time.Now())
			conn.SetReusable(false)
		}
		conn.Close()
		i--
	}
}

func Test_listenWSAndDial(t *testing.T) {
	assert := assert.On(t)
	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13142, internet.ListenOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path: "ws",
					}, assert),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	go func() {
		conn, err := listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		conn, err = listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		conn, err = listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		listen.Close()
	}()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13142), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path: "ws",
					}, assert),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	conn.Close()
	<-time.After(time.Second * 5)
	conn, err = Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13142), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path: "ws",
					}, assert),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	conn.Close()
	<-time.After(time.Second * 15)
	conn, err = Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13142), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path: "ws",
					}, assert),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	conn.Close()
}

func Test_listenWSAndDial_TLS(t *testing.T) {
	assert := assert.On(t)
	go func() {
		<-time.After(time.Second * 5)
		assert.Fail("Too slow")
	}()
	tlsSettings := &v2tls.Config{
		AllowInsecure: true,
		Certificate: []*v2tls.Certificate{
			{
				Certificate: ReadFile("./../../../testing/tls/cert.pem", assert),
				Key:         ReadFile("./../../../testing/tls/key.pem", assert),
			},
		},
	}
	anyTlsSettings, err := ptypes.MarshalAny(tlsSettings)
	assert.Error(err).IsNil()

	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13143, internet.ListenOptions{
		Stream: &internet.StreamConfig{
			SecurityType: internet.SecurityType_TLS,
			SecuritySettings: []*internet.SecuritySettings{{
				Type:     internet.SecurityType_TLS,
				Settings: anyTlsSettings,
			}},
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path:            "wss",
						ConnectionReuse: true,
					}, assert),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	go func() {
		conn, err := listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		listen.Close()
	}()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13143), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			SecurityType: internet.SecurityType_TLS,
			SecuritySettings: []*internet.SecuritySettings{{
				Type:     internet.SecurityType_TLS,
				Settings: anyTlsSettings,
			}},
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
					Settings: MarshalSettings(&Config{
						Path:            "wss",
						ConnectionReuse: true,
					}, assert),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	conn.Close()
}

func ReadFile(file string, assert *assert.Assert) []byte {
	b, err := ioutil.ReadFile(file)
	assert.Error(err).IsNil()
	return b
}
