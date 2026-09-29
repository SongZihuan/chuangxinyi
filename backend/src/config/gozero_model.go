package config

import (
	"github.com/spf13/viper"
	"github.com/wuntsong-org/go-zero-plus/rest"
)

type MiddlewaresConf struct {
	// Trace bool `json:"trace" yaml:"trace"` // 链路追踪，不使用
	// Log        bool `json:"log" yaml:"log"`  // 日志中间件，不使用
	Prometheus bool `json:"prometheus" yaml:"prometheus" mapstructure:"prometheus"` // 普罗米修斯
	//MaxConns   bool `json:"maxConns" yaml:"maxConns" mapstructure:"maxConns"`       // 限流（默认关闭）
	Breaker  bool `json:"breaker" yaml:"breaker" mapstructure:"breaker"`    // 熔断（默认关闭）
	Shedding bool `json:"shedding" yaml:"shedding" mapstructure:"shedding"` // 负载监控中心（默认关闭）
	//Timeout    bool `json:"timeout" yaml:"timeout" mapstructure:"timeout"`          // 超时
	// Recover    bool `json:"recover" yaml:"recover"`  // 回复 必须启用
	Metrics bool `json:"metrics" yaml:"metrics" mapstructure:"metrics"` // 指标
	//MaxBytes bool `json:"maxBytes" yaml:"maxBytes" mapstructure:"maxBytes"` // 最大content
	Gunzip bool `json:"gunzip" yaml:"gunzip" mapstructure:"gunzip"` // 压缩管理
}

type DevConfig struct {
	Host       string `json:"host" yaml:"host" mapstructure:"host"`
	Port       int64  `json:"port" yaml:"port" mapstructure:"port"` // 默认6060
	MetricsUrl string `json:"MetricsUrl" yaml:"metricsUrl" mapstructure:"metricsUrl"`
}

type GoZeroHttpConfig struct {
	Host           string          `json:"host" yaml:"host" mapstructure:"host"`
	Port           int64           `json:"port" yaml:"port" mapstructure:"port"`
	MaxConns       int64           `json:"maxConns" yaml:"maxConns" mapstructure:"maxConns"` // 默认10000
	MaxBytes       int64           `json:"maxBytes" yaml:"maxBytes" mapstructure:"maxBytes"` // 默认1048576
	Timeout        int64           `json:"timeout" yaml:"timeout" mapstructure:"timeout"`    // 默认3000
	LogServiceName string          `json:"logServiceName" yaml:"logServiceName" mapstructure:"logServiceName"`
	LogLevel       string          `json:"logLevel" yaml:"logLevel" mapstructure:"logLevel"` // debug,info,error,severe
	Dev            DevConfig       `json:"dev" yaml:"dev" mapstructure:"dev"`
	Middlewares    MiddlewaresConf `json:"middlewares" yaml:"middlewares" mapstructure:"middlewares"`
}

func (c *GoZeroHttpConfig) GetRestConfig() rest.RestConf {
	mode := BackendConfig.GetMode()
	modeString := ""
	switch mode {
	case RunModeDevelop:
		modeString = "dev"
	case RunModeRelease:
		modeString = "pro"
	default:
		modeString = "pro"
	}

	config := rest.RestConf{}

	config.Host = c.Host
	config.Port = int(c.Port)
	config.Mode = modeString
	config.MaxBytes = c.MaxBytes
	config.MaxConns = int(c.MaxConns)
	config.Log.ServiceName = c.LogServiceName
	config.Log.Mode = "console"
	config.Log.Encoding = "plain"
	config.Log.Level = c.LogLevel
	config.MetricsUrl = c.Dev.MetricsUrl
	config.DevServer.Enabled = true
	config.DevServer.EnableMetrics = true
	config.DevServer.EnablePprof = true
	config.DevServer.Host = c.Dev.Host
	config.DevServer.Port = int(c.Dev.Port)
	config.DevServer.MetricsPath = "/metrics"
	config.DevServer.HealthPath = "/healthz"

	config.Middlewares.Trace = false
	config.Middlewares.Log = false
	config.Middlewares.Recover = false
	config.Middlewares.Prometheus = c.Middlewares.Prometheus
	config.Middlewares.Gunzip = c.Middlewares.Gunzip
	config.Middlewares.Metrics = c.Middlewares.Metrics
	config.Middlewares.MaxConns = true
	config.Middlewares.MaxBytes = true
	config.Middlewares.Timeout = false
	config.Middlewares.Shedding = c.Middlewares.Shedding
	config.Middlewares.Breaker = c.Middlewares.Breaker

	return config
}

func GoZeroHttpSetDefaultValue(v *viper.Viper, prefix string, logServiceName string, serverPort int64, devPort int64) {
	if len(prefix) != 0 {
		prefix += "."
	}

	v.SetDefault(prefix+"host", "0.0.0.0")
	v.SetDefault(prefix+"port", serverPort)
	v.SetDefault(prefix+"maxConns", 10000)
	v.SetDefault(prefix+"maxBytes", 20971520)
	v.SetDefault(prefix+"logServiceName", logServiceName)
	v.SetDefault(prefix+"logLevel", "info")
	v.SetDefault(prefix+"dev.host", "127.0.0.1")
	v.SetDefault(prefix+"dev.port", devPort)
	v.SetDefault(prefix+"middlewares.prometheus", true)
	v.SetDefault(prefix+"middlewares.maxConns", false)
	v.SetDefault(prefix+"middlewares.breaker", false)
	v.SetDefault(prefix+"middlewares.shedding", false)
	v.SetDefault(prefix+"middlewares.timeout", true)
	v.SetDefault(prefix+"middlewares.metrics", true)
	v.SetDefault(prefix+"middlewares.maxBytes", true)
	v.SetDefault(prefix+"middlewares.gunzip", true)
}
