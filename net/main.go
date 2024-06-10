package net

import (
	timeUtil "github.com/daida459031925/common/time"
	"net"
	"net/http"
	"time"
)

const (
	httpClientTimeout = 10
	httpDialerTimeout = 10
	httpKeepAlive     = 30
)

// HttpClient 基本设置选项
type HttpClient struct {
	client  *http.Client                                 // http 客户端
	timeout time.Duration                                // 超时时间
	proxy   string                                       // 代理地址
	dialer  func(network, addr string) (net.Conn, error) // tcp or udp 连接
}

// NewHttpClient 创建连接
func NewHttpClient() (*HttpClient, error) {
	httpClient := &HttpClient{
		client: &http.Client{},
	}
	//默认十秒钟接口返回
	httpClient.SetTimeOut(httpClientTimeout)
	return httpClient, nil
}

// SetTimeOut 设置http请求超时时间
func (c *HttpClient) SetTimeOut(second int) {
	duration := timeUtil.GetSecond(second)
	c.timeout = duration
	c.client.Timeout = duration
}

// SetTransport 传输基本设置
func (c *HttpClient) SetTransport() {

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		//指定一个空闲连接的最大空闲时间，超过此时间的连接将会被关闭。这可以帮助客户端保持对服务器的长连接数量的限制
		IdleConnTimeout: 90 * time.Second,
		//TLSClientConfig用于设置TLS配置，包括TLS版本、密码套件、CA证书等
		TLSHandshakeTimeout: 10 * time.Second,
		//预期继续超时 不知道是不是代理超时时间
		ExpectContinueTimeout: 1 * time.Second,
		//远程服务器代理，默认没有代理
		Proxy: nil,
	}
	c.client.Transport = transport
}

//func (c *HttpClient) SetProxy(proxy string) {
//	if strUtil.TrimSpaceAndEmpty(proxy) {
//		return
//	}
//
//	u, b1 := strUtil.IsURL(proxy)
//	b2 := strUtil.IsIPv4(proxy)
//	b3 := strUtil.IsIPv6(proxy)
//
//	if !strUtil.IsIPOrURLToBool(b1, b2, b3) {
//		return
//	}
//
//	conn, e := net.Dial("tcp", u.Host)
//
//	if e != nil {
//		err.RuntimeExceptionTF(true, e)
//	}
//
//	dialer := &net.Dialer{
//		Timeout:   timeUtil.GetSecond(httpDialerTimeout),
//		KeepAlive: timeUtil.GetSecond(httpKeepAlive),
//	}
//
//	if b2 || b3 {
//
//	}
//
//	func(request *http.Request) (*url.URL, error) {
//		return url.Parse(proxy)
//	}
//
//	transport := &http.Transport{
//		DialContext: (&net.Dialer{
//			Timeout:   10 * time.Second,
//			KeepAlive: 30 * time.Second,
//		}).DialContext,
//		//指定一个空闲连接的最大空闲时间，超过此时间的连接将会被关闭。这可以帮助客户端保持对服务器的长连接数量的限制
//		IdleConnTimeout: 90 * time.Second,
//		//TLSClientConfig用于设置TLS配置，包括TLS版本、密码套件、CA证书等
//		TLSHandshakeTimeout: 10 * time.Second,
//		//预期继续超时 不知道是不是代理超时时间
//		ExpectContinueTimeout: 1 * time.Second,
//		//远程服务器代理，默认没有代理
//		Proxy: nil,
//	}
//	c.client.Transport = transport
//}
//
//// SetHeader 代理的话Proxy-Authorization可能使用的也是Authorization
//func SetHeader(key, name, password string, header http.Header) {
//	header.Set(key, strUtil.GetBasicAuth(name, password))
//}
//
//func NewHttpProxyClient(proxyURL, proxyUser, proxyPassword string) (*HttpClient, error) {
//	var proxy *url.URL = nil
//	if proxyURL != "" {
//		var e error = nil
//		proxy, e = url.Parse(proxyURL)
//		if e != nil {
//			return nil, e
//		}
//	}
//
//	client, e := NewHttpClient()
//	if e != nil {
//		err.RuntimeExceptionTF(true, e)
//	}
//
//	return &HttpClient{
//		client:        client.client,
//		proxyURL:      proxy,
//		proxyUser:     proxyUser,
//		proxyPassword: proxyPassword,
//	}, nil
//}
//
//func (c *HttpClient) SendRequest(method, urlStr string, headers map[string]string, body []byte) (*http.Response, error) {
//	request, err := http.NewRequest(method, urlStr, bytes.NewBuffer(body))
//	if err != nil {
//		return nil, err
//	}
//	for key, value := range headers {
//		request.Header.Set(key, value)
//	}
//	if c.proxyURL != nil {
//		if c.proxyUser != "" {
//			user := url.UserPassword(c.proxyUser, c.proxyPassword)
//			c.proxyURL.User = user
//		}
//		request.URL.Opaque = urlStr
//		request.URL.Scheme = c.proxyURL.Scheme
//		request.URL.Host = c.proxyURL.Host
//	}
//	return c.client.Do(request)
//}
//
//import (
//	"bytes"
//	err "github.com/daida459031925/common/error"
//	timeUtil "github.com/daida459031925/common/time"
//	"github.com/daida459031925/common/util/strutil"
//	"net"
//	"net/http"
//	"net/url"
//	"time"
//)
//
//const (
//	httpClientTimeout = 10
//	httpDialerTimeout = 10
//	httpKeepAlive     = 30
//)
//
//// HttpClient 基本设置选项
//type HttpClient struct {
//	client  *http.Client                                 // http 客户端
//	timeout time.Duration                                // 超时时间
//	proxy   string                                       // 代理地址
//	dialer  func(network, addr string) (net.Conn, error) // tcp or udp 连接
//}
//
//// NewHttpClient 创建连接
//func NewHttpClient() (*HttpClient, error) {
//	httpClient := &HttpClient{
//		client: &http.Client{},
//	}
//	//默认十秒钟接口返回
//	httpClient.SetTimeOut(httpClientTimeout)
//	return httpClient, nil
//}
//
//// SetTimeOut 设置http请求超时时间
//func (c *HttpClient) SetTimeOut(second int) {
//	duration := timeUtil.GetSecond(second)
//	c.timeout = duration
//	c.client.Timeout = duration
//}
//
//// SetTransport 传输基本设置
//func (c *HttpClient) SetTransport() {
//
//	transport := &http.Transport{
//		DialContext: (&net.Dialer{
//			Timeout:   10 * time.Second,
//			KeepAlive: 30 * time.Second,
//		}).DialContext,
//		//指定一个空闲连接的最大空闲时间，超过此时间的连接将会被关闭。这可以帮助客户端保持对服务器的长连接数量的限制
//		IdleConnTimeout: 90 * time.Second,
//		//TLSClientConfig用于设置TLS配置，包括TLS版本、密码套件、CA证书等
//		TLSHandshakeTimeout: 10 * time.Second,
//		//预期继续超时 不知道是不是代理超时时间
//		ExpectContinueTimeout: 1 * time.Second,
//		//远程服务器代理，默认没有代理
//		Proxy: nil,
//	}
//	c.client.Transport = transport
//}
//
//func (c *HttpClient) SetProxy(proxy string) {
//	if strUtil.TrimSpaceAndEmpty(proxy) {
//		return
//	}
//
//	u, b1 := strUtil.IsURL(proxy)
//	b2 := strUtil.IsIPv4(proxy)
//	b3 := strUtil.IsIPv6(proxy)
//
//	if !strUtil.IsIPOrURLToBool(b1, b2, b3) {
//		return
//	}
//
//	conn, e := net.Dial("tcp", u.Host)
//
//	if e != nil {
//		err.RuntimeExceptionTF(true, e)
//	}
//
//	dialer := &net.Dialer{
//		Timeout:   timeUtil.GetSecond(httpDialerTimeout),
//		KeepAlive: timeUtil.GetSecond(httpKeepAlive),
//	}
//
//	if b2 || b3 {
//
//	}
//
//	func(request *http.Request) (*url.URL, error) {
//		return url.Parse(proxy)
//	}
//
//	transport := &http.Transport{
//		DialContext: (&net.Dialer{
//			Timeout:   10 * time.Second,
//			KeepAlive: 30 * time.Second,
//		}).DialContext,
//		//指定一个空闲连接的最大空闲时间，超过此时间的连接将会被关闭。这可以帮助客户端保持对服务器的长连接数量的限制
//		IdleConnTimeout: 90 * time.Second,
//		//TLSClientConfig用于设置TLS配置，包括TLS版本、密码套件、CA证书等
//		TLSHandshakeTimeout: 10 * time.Second,
//		//预期继续超时 不知道是不是代理超时时间
//		ExpectContinueTimeout: 1 * time.Second,
//		//远程服务器代理，默认没有代理
//		Proxy: nil,
//	}
//	c.client.Transport = transport
//}
//
//// SetHeader 代理的话Proxy-Authorization可能使用的也是Authorization
//func SetHeader(key, name, password string, header http.Header) {
//	header.Set(key, strUtil.GetBasicAuth(name, password))
//}
//
//func NewHttpProxyClient(proxyURL, proxyUser, proxyPassword string) (*HttpClient, error) {
//	var proxy *url.URL = nil
//	if proxyURL != "" {
//		var e error = nil
//		proxy, e = url.Parse(proxyURL)
//		if e != nil {
//			return nil, e
//		}
//	}
//
//	client, e := NewHttpClient()
//	if e != nil {
//		err.RuntimeExceptionTF(true, e)
//	}
//
//	return &HttpClient{
//		client:        client.client,
//		proxyURL:      proxy,
//		proxyUser:     proxyUser,
//		proxyPassword: proxyPassword,
//	}, nil
//}
//
//func (c *HttpClient) SendRequest(method, urlStr string, headers map[string]string, body []byte) (*http.Response, error) {
//	request, err := http.NewRequest(method, urlStr, bytes.NewBuffer(body))
//	if err != nil {
//		return nil, err
//	}
//	for key, value := range headers {
//		request.Header.Set(key, value)
//	}
//	if c.proxyURL != nil {
//		if c.proxyUser != "" {
//			user := url.UserPassword(c.proxyUser, c.proxyPassword)
//			c.proxyURL.User = user
//		}
//		request.URL.Opaque = urlStr
//		request.URL.Scheme = c.proxyURL.Scheme
//		request.URL.Host = c.proxyURL.Host
//	}
//	return c.client.Do(request)
//}
