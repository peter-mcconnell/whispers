package main

import (
	"flag"

	"github.com/peter-mcconnell/whispers/pkg/config"
)

const (
	defaultBinPath = "/lib/x86_64-linux-gnu/libpam.so.0"
	defaultSymbol  = "pam_get_authtok"
)

func cfgFromFlags() *config.Config {
	cfg := &config.Config{}
	flag.StringVar(&cfg.BinPath, "binPath", defaultBinPath, "Path to the binary")
	flag.StringVar(&cfg.Symbol, "symbol", defaultSymbol, "Symbol to target")
	flag.Parse()

	return cfg
}
