package main

import (
	. "tinnitus"
)

func main() {
	Initialize(DefaultTinnitusMode)

	Logger.Info("Starting Tinnitus...")

	tinnitus := NewTinnitus()

	err := tinnitus.Run()

	if err != nil {
		Logger.Fatal("Fatal runtime error", err)
	}
}
