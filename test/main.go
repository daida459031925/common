package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	HTTP  = "http://"
	HTTPS = "https://"
)

var (
	httpSync sync.Once
	client   *http.Client
)

// 使用sync的once可以只执行一次方法，或者使用init
func newNetClient() *http.Client {
	httpSync.Do(func() {
		client = &http.Client{
			//它涵盖整个交互过程，从发起连接到接收响应报文结束
			Timeout:   30 * time.Second,
			Transport: http.DefaultTransport,
			//Dialer.Timeout限制创建一个tcp连接使用的时间
			//Transport.TLSHandshakeTimeout 限制TLS握手使用的时间
			//Transport.ResponseHeaderTimeout 限制读取响应报文头使用的时间
			//Transport.ExpectContinueTimeout 限制客户端在发送一个包含：100-continue的http报文头后，等待收到一个go-ahead响应报文所用的时间。在1.6中，此设置对HTTP/2无效。（在1.6.2中提供了一个特定的封装DefaultTransport）
			//Transport.IdleConnTimeout 连接最大空闲时间，超过这个时间就会被关闭
			//Transport.ExpectContinueTimeout  等待服务器的第一个响应headers的时间，0表示没有超时，则body会立刻发送，无需等待服务器批准，这个时间不包括发送请求header的时间
			//DisableKeepAlives true为代表开启长连接
			//MaxIdleConns 是长连接在关闭之前，连接池对所有host的最大链接数量
			//MaxIdleConnsPerHost： 连接池对每个host的最大链接数量(MaxIdleConnsPerHost <=MaxIdleConns,如果客户端只需要访问一个host，那么最好将MaxIdleConnsPerHost与MaxIdleConns设置为相同，这样逻辑更加清晰)
		}
	})

	//http.d
	//client.Transport.RoundTrip()

	return client
}

// 自定义常量将http的请求模式进行封装不能传入字符串,使用定义的常量
const (
	GET     httpType = http.MethodGet
	HEAD    httpType = http.MethodHead
	POST    httpType = http.MethodPost
	PUT     httpType = http.MethodPut
	PATCH   httpType = http.MethodPatch
	DELETE  httpType = http.MethodDelete
	CONNECT httpType = http.MethodConnect
	OPTIONS httpType = http.MethodOptions
	TRACE   httpType = http.MethodTrace
)

// 工具类统一采用结构体（对象）形式，简单算法内容采用直接调用模型
type (
	httpType string

	//http.Client内部都存在协程
	httpx struct {
		url    string
		body   map[string]any
		header map[string]any
	}

	Resp struct {
		header http.Header
		result Result
		code   int
	}
	// Result 其他工具类中返还的数据内容
	Result struct {
		Status    int16        `json:"status"` //状态类型
		Msg       string       `json:"msg"`    //错误时候的返回信息
		Data      any          `json:"data"`   //返还的数据
		Date      string       `json:"date"`   //记录数据返回时间
		funcSlice []func() any //记录当前需要执行的所有任务
	}
)

// Url 使用自定义的http进行初始化请求，添加数据模式采用链式编程
func Url(url string) *httpx {
	if len(url) <= 0 {
		return nil
	}

	return &httpx{
		url: url,
	}
}

// Http 设置请求内容为http请求
func (h *httpx) Http() *httpx {
	h.url = join(HTTP, h.url)
	return h
}

// Https 设置请求为https请求
func (h *httpx) Https() *httpx {
	h.url = join(HTTPS, h.url)
	return h
}

// 判断传入的url中是否添加了http://或https://
func join(types, url string) string {
	b := bytes.Buffer{}
	if !(strings.Contains(url, HTTP) || strings.Contains(url, HTTPS)) {
		b.WriteString(types)
	}
	b.WriteString(url)
	return b.String()
}

// GetHeader 返回resp的header
func (r *Resp) GetHeader() http.Header {
	if r != nil {
		return r.header
	}
	return nil
}

// GetContent 返回response的[]byte格式内容
func (r *Resp) GetContent() Result {
	return r.result
}

// GetStatusCode 返回response的http code
func (r *Resp) GetStatusCode() int {
	if r != nil {
		return r.code
	}
	return 500
}

// 获取返还值
func getRespBody(resp *http.Response) (response *Resp, e error) {
	if resp == nil {
		return nil, errors.New("response is nil")
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r Result
	err = json.Unmarshal(content, &r)
	if e != nil {
		return nil, err
	}

	return &Resp{
		header: resp.Header,
		result: r,
		code:   resp.StatusCode,
	}, nil
}

func main() {
	resp, e := http.Get(Url("www.baidu.com").Http().url)
	if e != nil {
		return
	}
	defer resp.Body.Close()

	_, e = io.ReadAll(resp.Body)
	if e != nil {
		return
	}

	//fmt.Println(string(body))
	//client.Do()
}

func (h *httpx) NewRequest(ht httpType) (*http.Request, error) {
	if len(strings.TrimSpace(h.url)) <= 0 {
		return nil, errors.New("url不能为空")
	}

	var reader io.Reader
	if h.body != nil && len(h.body) > 0 {
		body, e := json.Marshal(h.body)
		if e != nil {
			return nil, errors.New("解析body失败")
		}
		reader = strings.NewReader(string(body))
	}

	req, e := http.NewRequest(fmt.Sprintf("%s", ht), h.url, reader)
	if e != nil {
		return nil, errors.New("http访问连接创建失败")
	}

	if h.header != nil && len(h.header) > 0 {
		for s := range h.header {
			req.Header.Add(s, fmt.Sprintf("%s", h.header[s]))
		}
	}

	return req, e
}

// go test -bench=方法名 -benchmem
// ns/op平均每次多少时间 1s=1000ms 1ms=1000us 1us=1000ns
// allocs/op进行多少次内存分配
// B/op标识每次操作分配多少字节
func BenchmarkName(b *testing.B) {
	main1()
	//for i := 0; i < b.N; i++ {
	//	a()
	//}
}

var once sync.Once

func a() {

	once.Do(func() {
		fmt.Println(fmt.Sprintf("%s", GET))
		fmt.Println(1<<63 - 1)
		fmt.Println("1")
	})
	once.Do(func() {
		fmt.Print("2")
	})
}

func main1() {
	f1, err := os.Open("/home/sga/图片/2022-12-06 16-44-03屏幕截图1111111111111111111111111.png")
	if err != nil {
		panic(err)
	}
	f2, err := os.Open("/home/sga/图片/166908228357506.png")
	if err != nil {
		panic(err)
	}
	f3, err := os.Create("/home/sga/图片/appicon_2.jpg")
	if err != nil {
		panic(err)
	}
	m1, err := png.Decode(f1)
	if err != nil {
		panic(err)
	}
	//获取目标图片最小像素点和最大像素点
	bounds := m1.Bounds()
	fmt.Println(bounds.Max)
	fmt.Println(bounds.Min)
	m2, err := png.Decode(f2)
	if err != nil {
		panic(err)
	}
	p := image.Pt(0, 0)
	m := image.NewRGBA(bounds)
	white := color.RGBA{255, 255, 255, 255}
	//绘制图片
	draw.Draw(m, bounds, &image.Uniform{white}, p, draw.Src)
	draw.Draw(m, bounds, m1, p, draw.Src)
	draw.Draw(m, image.Rect(100, 200, 300, 600), m2, image.Pt(0, 0), draw.Src)
	err = jpeg.Encode(f3, m, &jpeg.Options{90})
	if err != nil {
		panic(err)
	}
	fmt.Printf("okn")
}

//"https://blog.csdn.net/ydl1128/article/details/126259943"
//"https://www.cnblogs.com/superhin/p/16332720.html"
//"https://gitee.com/dianjiu/gokit"
//"https://blog.csdn.net/asd1126163471/article/details/127020095"
//"https://www.php.cn/be/go/475737.html"
//"https://www.cnblogs.com/hsyw/p/16104591.html"
//"https://laravelacademy.org/post/21003"
