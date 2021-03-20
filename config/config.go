package config

import (
	"fmt"
	"net/http"
)

type Config struct {
	Version     string
	VersionDate string
	HTTPAddr    string `json:"http_addr"`
	DebugAddr   string `json:"debug_addr"`
	GRPCAddr    string `json:"grpc_addr"`
}

func (c *Config) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Version:%s\nDate:%s", c.Version, c.VersionDate)
}
