package client

import (
	"github.com/daida459031925/common/fmt"
	"github.com/gogf/gf/v2/net/gtcp"
)

// NewConn 创建链接
func NewConn(ip string, port int) (*gtcp.Conn, error) {
	conn, err := gtcp.NewConn(fmt.Sprintf("%s:%d", ip, port))
	return conn, err
}

// SendString 向指定链接发送消息
func SendString(conn *gtcp.Conn, str string, back ...func(conn *gtcp.Conn, e error)) {
	SendByte(conn, []byte(str), back...)
}

// SendByte 向指定链接发送消息
func SendByte(conn *gtcp.Conn, b []byte, back ...func(conn *gtcp.Conn, e error)) {
	if err := conn.Send(b); err != nil {
		if len(back) > 0 {
			for i := range back {
				back[i](conn, err)
			}
		}
	}
}
