package main

import (
	onelog "github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	wapc "github.com/wapc/wapc-guest-tinygo"
)

// Logger is the logger instance used throughout the policy。
var Logger *onelog.Logger

func init() {
	// Initialize the logger
	logWriter := kubewarden.KubewardenLogWriter{}
	Logger = onelog.New(
		&logWriter,
		onelog.ALL, // shortcut for onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL
	)
}

func main() {
	wapc.RegisterFunctions(wapc.Functions{
		"validate":          validate,
		"validate_settings": ValidateSettings,
	})
}
