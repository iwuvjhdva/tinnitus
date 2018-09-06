package tinnitus

import (
	"flag"
	"path"
)

var Flags TinnitusFlags

type TinnitusFlags struct {
	Verbose *bool
	Debug   *bool
	Version *bool
	Log     *string
	Config  *string
	Live    *bool
	From    *int64
	To      *int64
}

func InitFlags() {
	defaultConfigPath := path.Join(PackagePath(), "config/config.yaml")

	Flags = TinnitusFlags{
		Verbose: flag.Bool("verbose", false, "force verbose output"),
		Debug:   flag.Bool("debug", false, "enable debug mode"),
		Version: flag.Bool("version", false, "print program version"),
		Log:     flag.String("log", "trader.log", "log file path"),
		Config:  flag.String("config", defaultConfigPath, "config file path"),
		Live:    flag.Bool("live", false, "go live"),
		From:    flag.Int64("from", 0, "start from the block number"),
		To:      flag.Int64("to", -1, "stop at the block number"),
	}

	flag.Parse()
}
