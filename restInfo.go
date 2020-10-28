package main


// restInfo
type RestInfo struct {
	// 注释
	Comment string
	// 请求头
	Header map[string]string
	// 请求方法
	Method string
	// 路径
	Path string
	// body
	Body string
}

