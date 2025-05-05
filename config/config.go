package config

// Config 存储代理服务器配置
type Config struct {
	Port            int    // 服务器监听端口
	MaxContentSize  int64  // 最大下载内容大小（字节）
	Timeout         int    // 下载超时时间（秒）
	UserAgent       string // 请求的 User-Agent
	AllowedProtocols []string // 允许的协议
}

// NewConfig 创建默认配置
func NewConfig(port int) *Config {
	return &Config{
		Port:            port,
		MaxContentSize:  1024 * 1024 * 1024, // 默认最大 1GB
		Timeout:         300,                 // 默认 5 分钟超时
		UserAgent:       "DownProxy/1.0",
		AllowedProtocols: []string{"http", "https"},
	}
}