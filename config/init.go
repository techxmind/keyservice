package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	DefaultConfig *Config
)

func init() {
	DefaultConfig = &Config{}

	flag.StringVar(&DefaultConfig.DebugAddr, "debug.addr", ":5060", "Debug and metrics listen address")
	flag.StringVar(&DefaultConfig.HTTPAddr, "http.addr", ":5050", "HTTP listen address")
	flag.StringVar(&DefaultConfig.GRPCAddr, "grpc.addr", ":5040", "gRPC (HTTP) listen address")

	// Use environment variables, if set. Flags have priority over Env vars.
	if addr := os.Getenv("DEBUG_ADDR"); addr != "" {
		DefaultConfig.DebugAddr = addr
	}
	if port := os.Getenv("PORT"); port != "" {
		DefaultConfig.HTTPAddr = fmt.Sprintf(":%s", port)
	}
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		DefaultConfig.HTTPAddr = addr
	}
	if addr := os.Getenv("GRPC_ADDR"); addr != "" {
		DefaultConfig.GRPCAddr = addr
	}
}
