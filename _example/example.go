package main

import (
	"fmt"
	"github.com/hipai/gohttp"
)

func main() {
	var client = gohttp.NewHttpClient()

	//request := gohttp.Request{Url: "http://t.cn/RgrVhQ8", Method: gohttp.GET}
	req := gohttp.NewParam("http://dwz.cn/u3lNTYmt", gohttp.GET)
	resp := gohttp.Do(req, client)
	fmt.Print(resp.Text)
	fmt.Print(resp.Cookies)
}
