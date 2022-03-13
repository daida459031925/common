package translate

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Translate struct {
	Url    string `json:"url"`    //https://translate.googleapis.com/translate_a/single
	Client string `json:"client"` //默认值(不要修改) gtx
	Sl     string `json:"sl"`     //来源语言 en zh-cn 语言
	Tl     string `json:"tl"`     //目标语言 en zh-cn 语言
	Dt     string `json:"dt"`     //默认值(不要修改) t
	Q      string `json:"q"`      //翻译的文本 建议先url-encode
}

func (translate *Translate) toString() string {
	client := "?client=" + translate.Client
	sl := "&sl=" + translate.Sl
	tl := "&tl=" + translate.Tl
	dt := "&dt=" + translate.Dt
	q := "&q=" + url.QueryEscape(translate.Q)
	var build strings.Builder
	build.WriteString(translate.Url)
	build.WriteString(client)
	build.WriteString(sl)
	build.WriteString(tl)
	build.WriteString(dt)
	build.WriteString(q)

	return build.String()
}

func TranslateConversion(text string) (string, error) {
	translate := Translate{
		Url:    "https://translate.googleapis.com/translate_a/single",
		Client: "gtx",
		Sl:     "zh-cn",
		Tl:     "en",
		Dt:     "t",
		Q:      text,
	}
	toString := translate.toString()
	resp, err := http.Get(toString)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//返回的json反序列化比较麻烦, 直接字符串拆解
	ss := string(bs)
	ss = strings.ReplaceAll(ss, "[", "")
	ss = strings.ReplaceAll(ss, "]", "")
	ss = strings.ReplaceAll(ss, "null,", "")
	ss = strings.Trim(ss, `"`)
	ps := strings.Split(ss, `","`)
	return ps[0], nil
}
