package gohttp

import (
	"crypto/tls"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

//初始化httpClient
func NewHttpClient() *http.Client {

	//构建承载者 Transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			//略过证书校验(不安全)
			InsecureSkipVerify: true,
			//允许远程服务器重复请求
			Renegotiation: tls.RenegotiateFreelyAsClient},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	//HttpClient连接
	httpClient := &http.Client{
		//使用默认Transport
		//Transport: http.DefaultTransport,
		Transport: tr,
		//手动重定向
		//CheckRedirect: func(req *http.Param, via []*http.Param) error {
		//	return http.ErrUseLastResponse
		//},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	httpClient.Jar = jar

	return httpClient
}

//初始化IP代理httpClient
func NewHttpClientProxy(proxyUrl string) *http.Client {
	//解析代理地址
	urlProxy, _ := url.Parse(proxyUrl)
	//HttpClient连接
	httpClient := &http.Client{
		Transport: &http.Transport{
			//获取代理函数
			Proxy: http.ProxyURL(urlProxy),
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		//CheckRedirect: func(req *http.Param, via []*http.Param) error {
		//	return http.ErrUseLastResponse
		//},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	httpClient.Jar = jar
	return httpClient
}

//初始化SOCKS5代理httpClient
func NewHttpClientSocks5(addr string, auth *proxy.Auth) (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", addr,
		auth,
		&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		},
	)
	//HttpClient连接
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial:                  dialer.Dial,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		//CheckRedirect: func(req *http.Param, via []*http.Param) error {
		//	return http.ErrUseLastResponse
		//},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	httpClient.Jar = jar
	return httpClient, err
}
