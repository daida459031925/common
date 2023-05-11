package net

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	//"github.com/ProtonVPN/go-openvpn/openvpn3"
	"github.com/daida459031925/common/util/strutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestSslVpn(t *testing.T) {

	// SSL VPN的服务器地址和端口号
	endpoint := "https://219.240.210.103:1702"

	// 跳过SSL证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 创建一个HTTP客户端
	client := &http.Client{Transport: tr}

	// 创建一个HTTP请求
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 添加HTTP Basic Authentication的用户名和密码
	req.SetBasicAuth("vpn", "vpn")

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 读取HTTP响应
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 输出HTTP响应内容
	fmt.Println(string(body))

}

func TestNet(t *testing.T) {
	//219.240.210.103 tcp1702 udp1673
	//121.81.101.1:1938
	//proxyUrl, err := url.Parse("https://vpn566878232.opengw.net:1215")
	//proxyUrl, err := url.Parse("https://180.65.68.20:1215")
	proxyUrl, err := url.Parse("https://219.240.210.103:1702")
	//proxyUrl, err := url.Parse("https://14.38.164.215:1807")
	//proxyUrl, err := url.Parse("https://119.23.42.160:1795")
	//proxyUrl, err := url.Parse("https://vpn:vpn@118.41.53.106:1299")
	//proxyUrl, err := url.Parse("https://14.55.85.96:1939")
	//"73.151.247.26", "vpn", "vpn"
	conn, err := net.Dial("tcp", proxyUrl.Host)
	//udpAddr, err := net.ResolveUDPAddr("udp", proxyUrl.Host)
	if err != nil {
		fmt.Printf("Error creating HTTP client: %s\n", err.Error())
		return
	}

	//conn, err := net.DialUDP("udp", nil, udpAddr)

	if conn == nil {
		fmt.Printf("Failed to connect to proxy server\n")
		return
	}

	defer conn.Close()

	// 发送账号密码
	username := "vpn"
	password := "vpn"
	authPacket := make([]byte, 2+len(username)+len(password))
	binary.BigEndian.PutUint16(authPacket[:2], uint16(len(username)))
	copy(authPacket[2:], []byte(username))
	copy(authPacket[2+len(username):], []byte(password))
	if _, err = conn.Write(authPacket); err != nil {
		fmt.Println(err)
		return
	}

	// 建立 TLS 连接
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	tlsConn := tls.Client(conn, tlsConfig)
	if e := tlsConn.Handshake(); e != nil {
		fmt.Printf("Failed to handshake: %v", e)
		return
	}
	defer tlsConn.Close()

	// 设置读取超时时间
	tlsConn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// 向代理服务器发送ping数据包
	pingData := []byte("ping")
	_, err = tlsConn.Write(pingData)
	if err != nil {
		log.Fatal("向代理服务器发送ping数据包失败:", err)
	}

	// 读取代理服务器的响应数据包
	buffer := make([]byte, 1024)
	n, err := tlsConn.Read(buffer)
	if err != nil {
		log.Fatal("读取代理服务器响应数据包失败:", err)
	}

	// 验证连接是否成功
	if string(buffer[:n]) != "ping" {
		log.Fatal("连接代理服务器失败")
	}
	//
	//requestData := map[string]interface{}{
	//	"prompt":     "你好",
	//	"max_tokens": 5,
	//	"n":          1,
	//	"stop":       "\n",
	//}
	//requestBody, _ := json.Marshal(requestData)
	//
	//// 创建 API 请求
	//apiEndpoint := "https://api.openai.com/v1/completions"
	////apiEndpoint := "https://api.openai.com/v1/chat/completions"
	//// 发送 HTTP 请求
	//req, err := http.NewRequest("GET", apiEndpoint, bytes.NewBuffer(requestBody))
	//if err != nil {
	//	panic(err)
	//}
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "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")
	//
	//client := &http.Client{Transport: &http.Transport{
	//	Proxy: http.ProxyURL(proxyUrl),
	//	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
	//		return conn, nil
	//	},
	//}}
	//resp, err := client.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//// 处理响应
	//fmt.Printf("HTTP response status: %d\n", resp.StatusCode)
}

func TestNet1(t *testing.T) {
	// 配置TLS客户端
	//caCert, err := x509.SystemCertPool()
	//if err != nil {
	//	panic(err)
	//}
	//tlsConfig := &tls.Config{
	//	InsecureSkipVerify: false,
	//	RootCAs:            caCert,
	//}

	// 建立TCP连接
	tlsConn, err := net.Dial("udp", "118.41.53.106:1299")
	if err != nil {
		panic(err)
	}

	// 建立TLS加密连接
	//tlsConn := tls.Client(tcpConn, tlsConfig)
	//err = tlsConn.Handshake()
	//if err != nil {
	//	panic(err)
	//}

	// 读取OpenVPN服务器发送的握手信息
	reader := bufio.NewReader(tlsConn)
	//line, _, err := reader.ReadLine()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(line))

	// 发送OpenVPN客户端认证信息
	fmt.Fprint(tlsConn, "AUTHENTICATE \"vpn\" \"vpn\"\n")
	line, _, e := reader.ReadLine()
	if e != nil {
		panic(e)
	}
	fmt.Println(string(line))

	// 发送OpenVPN客户端配置信息
	fmt.Fprint(tlsConn, "SET_CONFIG echo \"Hello, world!\"\n")
	line, _, err = reader.ReadLine()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(line))

	// 关闭连接
	err = tlsConn.Close()
	if err != nil {
		panic(err)
	}
}

func TestNet2(t *testing.T) {
	// 创建一个Dialer对象
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// 设置代理地址和端口
	proxyUrl, _ := url.Parse("https://118.41.53.106:1299")

	// 使用Dialer对象建立TCP连接
	conn, err := dialer.Dial("udp", "www.google.com:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// 创建一个HTTP请求
	req, err := http.NewRequest("GET", "https://www.google.com", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 设置代理服务器
	req.Header.Set("Proxy-Authorization", strUtil.GetBasicAuth("vpn", "vpn")) // 如果需要认证
	transport := &http.Transport{
		Dial:  dialer.Dial,
		Proxy: http.ProxyURL(proxyUrl),
	}

	// 发送HTTP请求并打印结果
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}

func TestNet3(t *testing.T) {
	// 设置代理地址
	//proxyAddr, err := net.ResolveUDPAddr("udp", "118.41.53.106:1299")
	proxyAddr, err := net.ResolveUDPAddr("udp", "219.240.210.103:1673")
	if err != nil {
		fmt.Println("Resolve proxy address error:", err)
		return
	}

	// 设置本地地址
	localAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	if err != nil {
		fmt.Println("Resolve local address error:", err)
		return
	}

	// 创建UDP连接
	conn, err := net.DialUDP("udp", localAddr, proxyAddr)
	if err != nil {
		fmt.Println("Dial UDP error:", err)
		return
	}

	// 向连接发送数据
	_, err = conn.Write([]byte("Hello, UDP!"))
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}

	// 接收连接返回的数据
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println("Received from server:", string(buf[:n]))
}

func TestUdpVpnService(t *testing.T) {
	serverAddr, err := net.ResolveUDPAddr("udp", ":5000")
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			continue
		}

		go func() {
			fmt.Printf("Received %d bytes from %s\n", n, clientAddr.String())

			// TODO: Modify the data as needed, and forward it to the destination server

			response := []byte("Hello from the VPN server!")

			_, err = conn.WriteToUDP(response, clientAddr)
			if err != nil {
				fmt.Println("Error writing to client:", err)
				return
			}

			fmt.Printf("Sent %d bytes to %s\n", len(response), clientAddr.String())
		}()
	}
}

func TestUdpVpnClient(t *testing.T) {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5000")
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		os.Exit(1)
	}

	clientAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	if err != nil {
		fmt.Println("Error resolving client address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", clientAddr, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// TODO: Send requests to the server and receive responses

	request := []byte("Hello from the VPN client!")
	_, err = conn.Write(request)
	if err != nil {
		fmt.Println("Error writing to server:", err)
		os.Exit(1)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		os.Exit(1)
	}

	fmt.Printf("Received %d bytes from server: %s\n", n, buffer[:n])
}

func TestUdpVpnClient1(t *testing.T) {
	// 创建一个 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// 使用 VPN 连接创建一个 TCP 连接
				return net.Dial("udp", "118.41.53.106:1299")
			},
		},
	}

	// 发送 HTTP 请求
	//resp, err := client.Get("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	resp, err := client.Get("https://www.zhihu.com/question/326354272")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func TestUdpVpnClient2(t *testing.T) {

	// 设置代理服务器
	proxyURL, err := url.Parse("https://219.240.210.103:1702")
	if err != nil {
		panic(err)
	}

	// 创建HTTP客户端并设置代理
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// 设置Chat GPT接口URL
	gptURL := "https://api.openai.com/v1/engines/davinci-codex/completions"

	// 创建POST请求体
	body := url.Values{}
	body.Set("prompt", "Hello, ChatGPT!")
	body.Set("max_tokens", "10")

	// 创建HTTP请求
	req, err := http.NewRequest("POST", gptURL, bytes.NewBufferString(body.Encode()))
	if err != nil {
		panic(err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer <your-openai-api-key>")

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应数据
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(respBody))

}

//func TestOpenVpnClient(t *testing.T) {
//
//	// Create a new OpenVPN session
//	session, err := openvpn3.NewSession(context.Background(), &openvpn3.SessionConfig{
//		UserCredentials: &openvpn3.UserCredentials{
//			Username: "your-username",
//			Password: "your-password",
//		},
//		ServerList: []*openvpn3.Server{
//			{
//				Protocol:      "udp",
//				Hostname:      "your-server-hostname",
//				Port:          1194,
//				Cipher:        "AES-256-CBC",
//				Auth:          "SHA256",
//				ServerCA:      "your-server-ca-certificate",
//				ClientCert:    "your-client-certificate",
//				ClientKey:     "your-client-private-key",
//				TLSAuthKey:    "your-tls-auth-key",
//				TLSAuthKeyDir: "your-tls-auth-key-dir",
//			},
//		},
//	})
//
//	if err != nil {
//		fmt.Println("Failed to create OpenVPN session:", err)
//		return
//	}
//
//	// Start the OpenVPN session
//	err = session.Start()
//
//	if err != nil {
//		fmt.Println("Failed to start OpenVPN session:", err)
//		return
//	}
//
//	// Wait for the session to establish a connection
//	for {
//		status, err := session.GetStatus()
//
//		if err != nil {
//			fmt.Println("Failed to get OpenVPN session status:", err)
//			break
//		}
//
//		if status.ConnectionState == openvpn3.Connected {
//			fmt.Println("OpenVPN session connected!")
//			break
//		}
//
//		time.Sleep(1 * time.Second)
//	}
//
//	// Send some data through the OpenVPN connection
//	_, err = session.Send([]byte("Hello, OpenVPN!"))
//
//	if err != nil {
//		fmt.Println("Failed to send data through OpenVPN connection:", err)
//	}
//
//	// Close the OpenVPN session
//	err = session.Close()
//
//	if err != nil {
//		fmt.Println("Failed to close OpenVPN session:", err)
//	}
//
//}
