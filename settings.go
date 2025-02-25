package main

import (
	"encoding/json"
	"fmt"

	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

// Settings represents the policy settings for validating Kubernetes deployment probes。
type Settings struct {
	// LivenessProbe specifies the requirements for liveness probe configuration。
	LivenessProbe ProbeConfig `json:"liveness_probe"`
	// ReadinessProbe specifies the requirements for readiness probe configuration。
	ReadinessProbe ProbeConfig `json:"readiness_probe"`
	// StartupProbe specifies the requirements for startup probe configuration。
	StartupProbe ProbeConfig `json:"startup_probe"`
}

// ProbeConfig represents the configuration requirements for a probe。
type ProbeConfig struct {
	// Required indicates whether the probe must be configured in the deployment。
	Required bool `json:"required"`
	// MinPeriodSeconds specifies the minimum allowed period between probe executions (in seconds)。
	MinPeriodSeconds int32 `json:"min_period_seconds,omitempty"`
	// MaxTimeoutSeconds specifies the maximum allowed timeout for probe execution (in seconds)。
	MaxTimeoutSeconds int32 `json:"max_timeout_seconds,omitempty"`
}

// DefaultSettings returns default settings。
func DefaultSettings() *Settings {
	return &Settings{
		LivenessProbe: ProbeConfig{
			Required: true,
		},
		ReadinessProbe: ProbeConfig{
			Required: true,
		},
		StartupProbe: ProbeConfig{
			Required: false,
		},
	}
}

// UnmarshalJSON unmarshals the settings with defaults。
func (s *Settings) UnmarshalJSON(data []byte) error {
	// Set defaults first。
	defaults := DefaultSettings()
	*s = *defaults

	// Define a type alias to avoid recursion。
	type SettingsAlias Settings
	alias := (*SettingsAlias)(s)

	// Unmarshal into the alias。
	if err := json.Unmarshal(data, alias); err != nil {
		return err
	}

	return nil
}

// Validate validates the Settings configuration。
func (s *Settings) Validate() error {
	// Validate liveness probe configuration。
	if err := s.validateProbeConfig("liveness probe", s.LivenessProbe); err != nil {
		return err
	}

	// Validate readiness probe configuration。
	if err := s.validateProbeConfig("readiness probe", s.ReadinessProbe); err != nil {
		return err
	}

	// Validate startup probe configuration。
	if err := s.validateProbeConfig("startup probe", s.StartupProbe); err != nil {
		return err
	}

	return nil
}

// validateProbeConfig validates individual probe configuration。
func (s *Settings) validateProbeConfig(probeName string, config ProbeConfig) error {
	if config.MinPeriodSeconds < 0 {
		return fmt.Errorf("%s: min_period_seconds must be non-negative", probeName)
	}
	if config.MaxTimeoutSeconds < 0 {
		return fmt.Errorf("%s: max_timeout_seconds must be non-negative", probeName)
	}
	if config.MinPeriodSeconds > 0 && config.MaxTimeoutSeconds > 0 &&
		config.MinPeriodSeconds <= config.MaxTimeoutSeconds {
		return fmt.Errorf("%s: min_period_seconds must be greater than max_timeout_seconds", probeName)
	}
	return nil
}

// validateSettings validates the settings。
func validateSettings(payload []byte) ([]byte, error) {
	// Parse the settings。
	settings := Settings{}
	err := json.Unmarshal(payload, &settings)
	if err != nil {
		logger.ErrorWith("cannot unmarshal settings").
			Err("error", err).
			Write()
		return kubewarden.RejectSettings(
			kubewarden.Message(fmt.Sprintf("cannot unmarshal settings: %v", err)))
	}

	// Validate the settings。
	err = settings.Validate()
	if err != nil {
		logger.ErrorWith("settings validation failed").
			Err("error", err).
			Write()
		return kubewarden.RejectSettings(
			kubewarden.Message(err.Error()))
	}

	logger.InfoWith("settings validation succeeded").Write()
	return kubewarden.AcceptSettings()
}

// NewSettingsFromValidationReq creates Settings from ValidationRequest。
func NewSettingsFromValidationReq(validationReq *kubewarden_protocol.ValidationRequest) (Settings, error) {
	settings := Settings{}
	err := json.Unmarshal(validationReq.Settings, &settings)
	if err != nil {
		return Settings{}, err
	}

	return settings, nil
}
