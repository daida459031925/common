package main

func main() {
	//proxyUrl, _ := url.Parse("http://<username>:<password>@<proxy_host>:<proxy_port>")
	//httpClient := &http.Client{
	//	Transport: &http.Transport{
	//		Proxy: http.ProxyURL(proxyUrl),
	//	},
	//}
	//
	//client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	//client.HTTPClient = httpClient
	//
	//completionRequest := &openai.CompletionRequest{
	//	Model:     "davinci",
	//	Prompt:    "请描述您的产品",
	//	MaxTokens: 60,
	//}
	//
	//response, err := client.CreateCompletion(completionRequest)
	//if err != nil {
	//	fmt.Printf("Error: %s\n", err.Error())
	//	return
	//}
	//
	//fmt.Println(response.Choices[0].Text)
}

//import (
//	"fmt"
//	"net/http"
//	"net/url"
//	"os"
//)
//
//func main() {
//	proxyUrl, _ := url.Parse("http://<username>:<password>@<proxy_host>:<proxy_port>")
//	httpClient := &http.Client{
//		Transport: &http.Transport{
//			Proxy: http.ProxyURL(proxyUrl),
//		},
//	}
//
//	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
//	client.HTTPClient = httpClient
//
//	completionRequest := &openai.CompletionRequest{
//		Model:     "davinci",
//		Prompt:    "请描述您的产品",
//		MaxTokens: 60,
//	}
//
//	response, err := client.CreateCompletion(completionRequest)
//	if err != nil {
//		fmt.Printf("Error: %s\n", err.Error())
//		return
//	}
//
//	fmt.Println(response.Choices[0].Text)
//}
