package service

import (
	"fmt"
	"net"
	"time"
)

func PortCheck(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", port), 500*time.Millisecond)
	if err != nil {
		return true
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
	return false
}
