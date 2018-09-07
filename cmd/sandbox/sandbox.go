package main

import (
	. "tinnitus"
)

func main() {
	Initialize(DefaultTinnitusMode)

	Logger.Info("Starting Tinnitus Sandbox...")

	tinnitus := NewTinnitus()

	err := tinnitus.Sandbox()

	if err != nil {
		Logger.Fatal("Fatal runtime error", err)
	}
}
