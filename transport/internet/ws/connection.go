package ws

import (
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrInvalidConn = errors.New("Invalid Connection.")
)

type ConnectionManager interface {
	Recycle(string, *wsconn)
}

type Connection struct {
	dest     string
	conn     *wsconn
	listener ConnectionManager
	reusable bool
	config   *Config
}

func NewConnection(dest string, conn *wsconn, manager ConnectionManager, config *Config) *Connection {
	return &Connection{
		dest:     dest,
		conn:     conn,
		listener: manager,
		reusable: config.ConnectionReuse,
		config:   config,
	}
}

func (this *Connection) Read(b []byte) (int, error) {
	if this == nil || this.conn == nil {
		return 0, io.EOF
	}

	return this.conn.Read(b)
}

func (this *Connection) Write(b []byte) (int, error) {
	if this == nil || this.conn == nil {
		return 0, io.ErrClosedPipe
	}
	return this.conn.Write(b)
}

func (this *Connection) Close() error {
	if this == nil || this.conn == nil {
		return io.ErrClosedPipe
	}
	if this.Reusable() {
		this.listener.Recycle(this.dest, this.conn)
		return nil
	}
	err := this.conn.Close()
	this.conn = nil
	return err
}

func (this *Connection) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *Connection) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *Connection) SetDeadline(t time.Time) error {
	return this.conn.SetDeadline(t)
}

func (this *Connection) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}

func (this *Connection) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}

func (this *Connection) SetReusable(reusable bool) {
	if !this.config.ConnectionReuse {
		return
	}
	this.reusable = reusable
}

func (this *Connection) Reusable() bool {
	return this.reusable
}

func (this *Connection) SysFd() (int, error) {
	return getSysFd(this.conn)
}

func getSysFd(conn net.Conn) (int, error) {
	return 0, ErrInvalidConn
}
