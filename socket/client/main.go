package main

import (
	"github.com/daida459031925/common/fmt"
	"github.com/gogf/gf/v2/net/gtcp"
)

func NewConn(ip string, port int) (*gtcp.Conn, error) {
	conn, err := gtcp.NewConn(fmt.Sprintf("%s:%d", ip, port))
	return conn, err
}
