package tinnitus

import (
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
)

type TinnitusMode int

const (
	DefaultTinnitusMode TinnitusMode = iota
	TestingTinnitusMode
)

var Mode TinnitusMode

var ctrlC chan os.Signal

var ShouldExitGracefully = func() bool {
	shouldExit := len(ctrlC) > 0

	if shouldExit {
		Logger.Info("Shutting down gracefully.")
	}

	return shouldExit
}

func PackagePath() string {
	_, fileName, _, ok := runtime.Caller(0)

	if !ok {
		panic("No caller information.")
	}

	return path.Dir(fileName)
}

func Initialize(tinnitusMode TinnitusMode) {
	if tinnitusMode == DefaultTinnitusMode {
		ctrlC = make(chan os.Signal, 1)
		signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)
	}

	InitFlags()
	InitConfig()
	InitLogger()
	InitSuperCollider()
	InitRPC()
}
