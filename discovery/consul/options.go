package consul

const (
	defaultAddr                  = "127.0.0.1:8500"
	defaultHealthCheck           = true
	defaultHealthCheckInterval   = 10
	defaultHealthCheckTimeout    = 5
	defaultCheckFailedDeregister = 300
)

type Option func(o *Options)

type Options struct {
	// 客户端连接地址
	// 内建客户端配置，默认为127.0.0.1:8500
	Addr string

	// 是否启用健康检查
	// 默认为true
	EnableHealthCheck bool

	// 健康检查时间间隔（秒），仅在启用健康检查后生效
	// 默认10秒
	HealthCheckInterval int

	// 健康检查超时时间（秒），仅在启用健康检查后生效
	// 默认5秒
	HealthCheckTimeout int

	// 健康检测失败后自动注销服务时间（秒）
	// 默认30秒
	CheckFailedDeregister int
}

func defaultOptions() *Options {
	return &Options{
		Addr:                  defaultAddr,
		EnableHealthCheck:     defaultHealthCheck,
		HealthCheckInterval:   defaultHealthCheckInterval,
		HealthCheckTimeout:    defaultHealthCheckTimeout,
		CheckFailedDeregister: defaultCheckFailedDeregister,
	}
}

// WithAddr 设置客户端连接地址
func WithAddr(addr string) Option {
	return func(o *Options) { o.Addr = addr }
}

// WithEnableHealthCheck 设置是否启用健康检查
func WithEnableHealthCheck(enable bool) Option {
	return func(o *Options) { o.EnableHealthCheck = enable }
}

// WithHealthCheckInterval 设置健康检查时间间隔
func WithHealthCheckInterval(interval int) Option {
	return func(o *Options) { o.HealthCheckInterval = interval }
}

// WithHealthCheckTimeout 设置健康检查超时时间
func WithHealthCheckTimeout(timeout int) Option {
	return func(o *Options) { o.HealthCheckTimeout = timeout }
}

// WithDeregisterCriticalServiceAfter 设置健康检测失败后自动注销服务时间
func WithDeregisterCriticalServiceAfter(after int) Option {
	return func(o *Options) { o.CheckFailedDeregister = after }
}
