package main

import (
	onelog "github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	wapc "github.com/wapc/wapc-guest-tinygo"
)

// Logger is the global logger instance
var (
	logWriter = kubewarden.KubewardenLogWriter{}
	Logger    = onelog.New(
		&logWriter,
		onelog.ALL, // shortcut for onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL
	)
)

func main() {
	wapc.RegisterFunctions(wapc.Functions{
		"validate":          validate,
		"validate_settings": ValidateSettings,
	})
}
