package config

type ProxyConfig struct {
	Files string `yaml:"files"`
}

type EntryConfig struct {
	ProxyEntries []ProxyEntry `yaml:"proxy_entries"`
}

type ProxyEntry struct {
	Name    string  `yaml:"name"`
	Target  string  `yaml:"target"` // 当前路由对应的后端地址
	Matches []Match `yaml:"matches"`
}

type Match struct {
	Path   string  `yaml:"path"`   // 路由
	Method string  `yaml:"method"` // 路由的请求方式
	Params []Param `yaml:"params"` // 参数
}

type Param struct {
	Location ParamLocation `yaml:"location"`
	// SessionKey session_key
	SessionKey string `yaml:"session_key"`
	// Rename 参数重命名,如果 rename 为空，则使用 session_key 为参数名
	Rename string `yaml:"rename"`
	// CustomValue 自定义的参数值，如果 session 中没有存某个值，支持自定义
	CustomValue interface{} `yaml:"custom_value"`
}
