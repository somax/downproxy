package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/somax/downproxy/config"
)

// ProxyHandler 处理下载代理请求
type ProxyHandler struct {
	config *config.Config
	client *http.Client
}

// NewProxyHandler 创建新的代理处理器
func NewProxyHandler(cfg *config.Config) *ProxyHandler {
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	return &ProxyHandler{
		config: cfg,
		client: client,
	}
}

// HandleDownload 处理下载请求
func (p *ProxyHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	// 只允许 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, "只支持 GET 请求", http.StatusMethodNotAllowed)
		return
	}

	// 获取要下载的 URL
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "缺少 url 参数", http.StatusBadRequest)
		return
	}

	// 验证 URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "无效的 URL", http.StatusBadRequest)
		return
	}

	// 检查协议
	if !p.isAllowedProtocol(parsedURL.Scheme) {
		http.Error(w, "不支持的协议", http.StatusBadRequest)
		return
	}

	// 获取文件名（用于 Content-Disposition 头）
	filename := p.extractFilename(parsedURL.Path)
	log.Printf("从URL路径提取的文件名: %s", filename)

	// 创建请求
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		http.Error(w, "创建请求失败", http.StatusInternalServerError)
		return
	}

	// 设置 User-Agent
	req.Header.Set("User-Agent", p.config.UserAgent)

	// 转发一些原始请求的头信息
	for _, header := range []string{"Range", "If-Modified-Since", "If-None-Match"} {
		if value := r.Header.Get(header); value != "" {
			req.Header.Set(header, value)
		}
	}

	// 执行请求
	log.Printf("代理下载: %s", targetURL)
	resp, err := p.client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("下载失败: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		http.Error(w, fmt.Sprintf("远程服务器返回错误: %s", resp.Status), resp.StatusCode)
		return
	}

	// 检查内容大小
	if resp.ContentLength > p.config.MaxContentSize {
		http.Error(w, "文件太大，超过允许的最大大小", http.StatusForbidden)
		return
	}

	// 尝试从响应头中获取文件名
	if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "" {
		log.Printf("原始Content-Disposition头: %s", contentDisposition)
		if newFilename := p.extractFilenameFromHeader(contentDisposition); newFilename != "" {
			log.Printf("从Content-Disposition头提取的文件名: %s", newFilename)
			filename = newFilename
		}
	}

	log.Printf("最终使用的文件名: %s", filename)

	// 设置响应头
	for key, values := range resp.Header {
		// 跳过一些不应该转发的头
		if key == "Server" || key == "Date" || key == "Connection" || key == "Content-Disposition" {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 设置内容处置头，使浏览器下载文件而不是显示
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// 设置状态码
	w.WriteHeader(resp.StatusCode)

	// 使用有限制的读取器，确保不超过最大大小
	limitedReader := io.LimitReader(resp.Body, p.config.MaxContentSize)

	// 将内容流式传输到客户端
	_, err = io.Copy(w, limitedReader)
	if err != nil {
		log.Printf("传输数据时出错: %v", err)
		// 此时已经开始发送响应，无法发送错误状态
	}

	log.Printf("下载完成: %s", targetURL)
}

// isAllowedProtocol 检查协议是否被允许
func (p *ProxyHandler) isAllowedProtocol(scheme string) bool {
	scheme = strings.ToLower(scheme)
	for _, allowed := range p.config.AllowedProtocols {
		if scheme == allowed {
			return true
		}
	}
	return false
}

// extractFilename 从路径中提取文件名
func (p *ProxyHandler) extractFilename(path string) string {
	// 从路径中提取文件名
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		filename := parts[len(parts)-1]
		if filename != "" {
			// 处理可能包含查询参数的情况
			queryIndex := strings.Index(filename, "?")
			if queryIndex > 0 {
				filename = filename[:queryIndex]
			}

			// 处理 URL 编码
			if decodedFilename, err := url.QueryUnescape(filename); err == nil {
				filename = decodedFilename
			} else {
				log.Printf("URL解码失败: %v", err) // 保留错误日志
			}

			return filename
		}
	}

	// 如果无法从路径中提取文件名，则返回默认文件名
	return "download"
}

// 新增函数：从 Content-Disposition 头中提取文件名
func (p *ProxyHandler) extractFilenameFromHeader(header string) string {
	if strings.Contains(header, "filename=") {
		parts := strings.Split(header, "filename=")
		if len(parts) > 1 {
			filename := parts[1]
			// 处理引号
			filename = strings.Trim(filename, `"'`)
			// 处理可能的额外参数
			if idx := strings.Index(filename, ";"); idx > 0 {
				filename = filename[:idx]
			}
			return filename
		}
	}
	return ""
}
