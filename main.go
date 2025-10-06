package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

// getClientIP 获取客户端真实IP地址
func getClientIP(r *http.Request) string {
	// 1. 优先检查 X-Forwarded-For 头（反向代理场景）
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 2. 检查 X-Real-IP 头（Nginx等反向代理）
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// 3. 直接从 RemoteAddr 获取
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// handler 处理HTTP请求
func handler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)

	// 设置响应头
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// 返回客户端IP
	fmt.Fprintf(w, "%s\n", clientIP)
}

func main() {
	// 命令行参数：监听端口
	port := flag.Int("port", 60080, "HTTP服务监听端口")
	flag.Parse()

	// 设置路由
	http.HandleFunc("/", handler)

	// 服务器监听地址
	addr := fmt.Sprintf(":%d", *port)

	log.Printf("服务器启动，监听端口 %d", *port)
	log.Printf("访问 http://localhost:%d 获取您的IP地址", *port)

	// 启动HTTP服务器
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
