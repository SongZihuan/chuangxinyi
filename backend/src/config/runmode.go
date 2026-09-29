//go:build !release

package config

import "net/http"

func (c *Config) GetMode() RunMode {
	switch c.Mode {
	case "Dev", "dev", "Develop", "develop":
		return RunModeDevelop
	case "Rel", "rel", "Release", "release", "prop", "Prop", "Product", "product":
		return RunModeRelease
	default:
		return RunModeDevelop
	}
}

func (c *Config) GetModeWithHeaderTrue(r *http.Request, header string) bool {
	mode := c.GetMode()
	if mode == RunModeRelease {
		return false
	}

	headerMode := r.Header.Get(header)
	if len(headerMode) == 0 {
		return false
	}

	switch headerMode {
	case "1", "#t", "#T", "T", "t", "True", "true", "on", "On", "ON":
		return true
	default:
		return false
	}
}

func (c *Config) GetModeFromRequests(r *http.Request) RunMode {
	mode := c.GetMode()
	if mode == RunModeRelease {
		return mode
	}

	headerMode := r.Header.Get("X-RunMode")
	if len(headerMode) == 0 {
		headerMode = r.URL.Query().Get("xrunmode")
	}

	switch headerMode {
	case "Dev", "dev", "Develop", "develop":
		return RunModeDevelop
	case "Rel", "rel", "Release", "release", "prop", "Prop", "Product", "product":
		return RunModeRelease
	default:
		return mode
	}
}
