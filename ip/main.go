package ip

import (
	"net"
)

// NewIp 创建IP地址
func NewIp(a, b, c, d byte, port int) net.TCPAddr {
	ip := net.IPv4(a, b, c, d)
	return net.TCPAddr{IP: ip, Port: port}
}
