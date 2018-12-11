package gohttp

import (
	"fmt"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	GET     = "GET"
	HEAD    = "HEAD"
	POST    = "POST"
	PUT     = "PUT"
	PATCH   = "PATCH"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
)

//MimeType Map
var MimeMap = map[string]string{
	"text/html":                ".html",
	"application/xml":          ".xml",
	"text/xml":                 ".xml",
	"text/css":                 ".css",
	"application/json":         ".json",
	"application/javascript":   ".js",
	"text/plain":               ".txt",
	"application/octet-stream": ".file",
	"application/pdf":          ".pdf",
	"application/msword":       ".doc",
	"application/vnd.ms-word":  ".doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   ".docx",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         ".xlsx",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
	"application/msexcel":          ".xls",
	"application/vnd.ms-excel":     ".xls",
	"application/mspowerpoint":     ".ppt",
	"application/zip":              ".zip",
	"application/x-zip-compressed": ".zip",
	"application/x-rar-compressed": ".rar",
	"image/jpeg":                   ".jpg",
	"image/jpg":                    ".jpg",
	"image/gif":                    ".gif",
	"image/png":                    ".png",
	"application/x-download":       ".download",
	"application/x-javascript":     ".json",
	"application":                  ".application",
}

type Request struct {
	//请求地址
	Url string
	//请求方式
	Method string
	//ua
	UserAgent string
	//文件大小长度
	FileLength int
	//支持采集的文件类型
	MimeType map[string]string
	//头信息
	HeaderMap map[string]string
	//表单数据
	FormBody map[string]string
}

type Response struct {
	//最终URL
	Url string
	//响应字节流
	Body []byte
	//响应状态码
	Code int
	//源码
	Text string
	//内容类型
	MimeType string
	//响应cookie
	Cookies string
	//文件类型
	FileType string
	//网页编码
	Charset string
	//响应时间
	ResponseTime int64
	//下载时间
	DownloadTime int64
	//总耗时
	TotalTime int64
}

func NewRequest(url, method string) *Request {
	return &Request{
		Url:       url,
		Method:    method,
		HeaderMap: make(map[string]string),
		FormBody:  make(map[string]string),
		MimeType:  make(map[string]string),
	}
}

func Do(req *Request, client *http.Client) *Response {
	resp := &Response{}
	uri := req.Url
	//判断url是否合法
	if !strings.HasPrefix(uri, "http") {
		log.Println("===============URL不合法===============")
		return resp
	}
	//判断请求方式并新建
	request, err := http.NewRequest("", uri, nil)
	if err != nil {
		print("===============处理错误===============")
	}

	//设置请求头信息
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded") //POST必填
	if req.UserAgent == "" {
		//设置默认UA
		req.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"
	}
	request.Header.Add("User-Agent", req.UserAgent)
	//浏览器支持的 MIME 类型分别是 text/html、application/xhtml+xml、application/xml 和 */*，
	//优先顺序是它们从左到右的排列顺序。
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	for name, value := range req.HeaderMap {
		request.Header.Add(name, value)
	}

	//POST请求处理
	switch req.Method {
	case POST:
		request.Method = "POST"
		values := url.Values{}
		for name, value := range req.FormBody {
			values.Add(name, value)
		}
		request.PostForm = values
	}

	//发起请求
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return resp
	}
	//保存响应cookie
	for _, v := range response.Cookies() {
		resp.Cookies += v.String()
	}

	//获取响应状态(重定向会自动处理)
	httpCode := response.StatusCode
	resp.Code = httpCode
	defer response.Body.Close()
	//处理body

	//最终url
	resp.Url = response.Request.URL.String()
	//读取内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	resp.Body = body

	//buf := new(bytes.Buffer)
	//buf.ReadFrom(response.Body)
	//buf.Bytes()

	//out, err := os.Create("output.txt")
	//defer out.Close()
	//n, err := io.Copy(out, response.Body)
	//contentType:=response.Header.Get("Content-Type")
	//返回cookies
	cookies := ""
	for _, cookie := range response.Cookies() {
		cookies += cookie.String() + ";"
	}
	resp.Cookies = cookies
	//检测ContentType
	mimeType := http.DetectContentType(body)
	resp.MimeType = mimeType
	//判断是否需要的类型
	fileType, ok := MimeMap[strings.TrimSpace(strings.Split(mimeType, ";")[0])]
	if ok {
		resp.FileType = fileType
		switch fileType {
		case ".html", ".xml", ".css", ".json", ".js", ".txt":
			//检测encoding
			enc, charSet, _ := charset.DetermineEncoding(body, mimeType)
			fmt.Println("charset:", charSet)
			if charSet == "utf-8" {
				resp.Text = string(body)
			} else { //转码
				bts, _ := enc.NewDecoder().Bytes(body)
				resp.Text = string(bts)
			}
		}
	}
	log.Println(response.Status + " | " + response.Proto + " | " + resp.MimeType + " | " + uri)
	return resp
}
