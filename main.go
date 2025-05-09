package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/somax/downproxy/config"
	"github.com/somax/downproxy/proxy"
	"github.com/somax/downproxy/static"
)

func main() {
	// 解析命令行参数
	port := flag.Int("port", 9527, "代理服务器监听端口")
	flag.Parse()

	// 初始化配置
	cfg := config.NewConfig(*port)

	// 创建代理服务处理器
	proxyHandler := proxy.NewProxyHandler(cfg)

	// 注册HTTP路由
	// 添加主页处理
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, static.IndexHTML)
	})

	// 添加下载处理路由
	http.HandleFunc("/download", proxyHandler.HandleDownload)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("下载代理服务已启动，监听端口: %d", cfg.Port)
	log.Printf("访问主页: http://localhost:%d/", cfg.Port)
	log.Printf("使用示例: http://localhost:%d/download?url=https://example.com/file.zip", cfg.Port)

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
