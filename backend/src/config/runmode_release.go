//go:build release

package config

import "net/http"

func (c *Config) GetMode() RunMode {
	return RunModeRelease
}

func (c *Config) GetModeWithHeaderTrue(r *http.Request, header string) bool {
	return false
}

func (c *Config) GetModeFromRequests(r *http.Request) RunMode {
	return RunModeRelease
}
