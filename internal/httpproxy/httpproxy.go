package httpproxy

import (
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// handleHTTP 处理 http 请求
func handleHTTP(w http.ResponseWriter, req *http.Request, dialer proxy.Dialer) {
	tp := http.Transport{
		Dial: dialer.Dial,
	}
	// 发起 socks 代理请求
	resp, err := tp.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	// StatusCode 写 Header
	w.WriteHeader(resp.StatusCode)
	// 将 body 直接复制写入到响应中
	io.Copy(w, resp.Body)
	log.Println(req.RemoteAddr, req.RequestURI, resp.StatusCode)
}

// 将 Header 写入到 Response
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// 通过劫持来实现隧道
func handleTunnel(w http.ResponseWriter, req *http.Request, dialer proxy.Dialer) {
	// 劫持请求
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	srcConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	// proxy 方式拨号
	dstConn, err := dialer.Dial("tcp", req.Host)
	if err != nil {
		_ = srcConn.Close()
		return
	}

	srcConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// 使用 goroutine 来实现 互相复制响应
	go transfer(dstConn, srcConn)
	go transfer(srcConn, dstConn)
	log.Println(req.RemoteAddr, req.RequestURI, "200")
}

func transfer(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()

	io.Copy(dst, src)
}

// 定义 serverHTTP 处理程序
func serveHTTP(w http.ResponseWriter, req *http.Request) {
	d := &net.Dialer{
		Timeout: 10 * time.Second,
	}
	// 使用 socks5 proxy
	dialer, _ := proxy.SOCKS5("tcp", viper.GetString("SocksProxy"), nil, d)

	// https Method is CONNECT
	if req.Method == "CONNECT" {
		handleTunnel(w, req, dialer)
	} else {
		handleHTTP(w, req, dialer)
	}
}

// HttpProxy 开启你的监听服务
func HttpProxy() {
	err := http.ListenAndServe(viper.GetString("ListenAddress"), http.HandlerFunc(serveHTTP))
	if err != nil {
		fmt.Println("http error: ", err.Error())
	}
}