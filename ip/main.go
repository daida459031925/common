package ip

import (
	"net"
)

// 创建IP地址
func CreateIp(a, b, c, d byte, port int) net.TCPAddr {
	ip := net.IPv4(a, b, c, d)
	return net.TCPAddr{IP: ip, Port: port}
}
