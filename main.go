package main

import (
	onelog "github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	wapc "github.com/wapc/wapc-guest-tinygo"
)

// PolicyHandler handles the policy validation requests
type PolicyHandler struct {
	logger *onelog.Logger
}

// NewPolicyHandler creates a new PolicyHandler instance
func NewPolicyHandler() *PolicyHandler {
	logWriter := kubewarden.KubewardenLogWriter{}
	logger := onelog.New(
		&logWriter,
		onelog.ALL, // shortcut for onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL
	)
	return &PolicyHandler{logger: logger}
}

// Validate handles the validation request
func (h *PolicyHandler) Validate(payload []byte) ([]byte, error) {
	return validate(payload, h.logger)
}

// ValidateSettings handles the settings validation request
func (h *PolicyHandler) ValidateSettings(payload []byte) ([]byte, error) {
	return ValidateSettings(payload, h.logger)
}

func main() {
	handler := NewPolicyHandler()
	wapc.RegisterFunctions(wapc.Functions{
		"validate":          handler.Validate,
		"validate_settings": handler.ValidateSettings,
	})
}
