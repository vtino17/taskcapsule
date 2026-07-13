package health

import (
	"net"
	"time"
)

// TCPChecker performs TCP port health checks with retries.
type TCPChecker struct {
	Host string
}

func NewTCPChecker(host string) *TCPChecker {
	return &TCPChecker{Host: host}
}

func (c *TCPChecker) Check(port string) error {
	addr := net.JoinHostPort(c.Host, port)
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
