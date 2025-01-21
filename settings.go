package main

import (
	"encoding/json"
	"fmt"

	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

// Settings represents the policy settings
type Settings struct {
	LivenessProbe  ProbeConfig `json:"liveness_probe"`
	ReadinessProbe ProbeConfig `json:"readiness_probe"`
	StartupProbe   ProbeConfig `json:"startup_probe"`
}

// ProbeConfig represents the configuration for a probe
type ProbeConfig struct {
	Required bool `json:"required"`
}

// DefaultSettings returns default settings
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

// UnmarshalJSON unmarshals the settings with defaults
func (s *Settings) UnmarshalJSON(data []byte) error {
	// Set defaults first
	defaults := DefaultSettings()
	*s = *defaults

	// Define a type alias to avoid recursion
	type SettingsAlias Settings
	alias := (*SettingsAlias)(s)

	// Unmarshal into the alias
	if err := json.Unmarshal(data, alias); err != nil {
		return err
	}

	return nil
}

// Validate validates the Settings
func (s *Settings) Validate() error {
	// All probes can be optional
	return nil
}

// ValidateSettings validates the settings
func ValidateSettings(payload []byte) ([]byte, error) {
	// Parse the settings
	settings := Settings{}
	err := json.Unmarshal(payload, &settings)
	if err != nil {
		Logger.ErrorWith("cannot unmarshal settings").
			Err("error", err).
			Write()
		return kubewarden.RejectSettings(
			kubewarden.Message(fmt.Sprintf("cannot unmarshal settings: %v", err)))
	}

	// Validate the settings
	err = settings.Validate()
	if err != nil {
		Logger.ErrorWith("settings validation failed").
			Err("error", err).
			Write()
		return kubewarden.RejectSettings(
			kubewarden.Message(err.Error()))
	}

	Logger.InfoWith("settings validation succeeded").Write()
	return kubewarden.AcceptSettings()
}

// NewSettingsFromValidationReq creates Settings from ValidationRequest
func NewSettingsFromValidationReq(validationReq *kubewarden_protocol.ValidationRequest) (Settings, error) {
	settings := Settings{}
	err := json.Unmarshal(validationReq.Settings, &settings)
	if err != nil {
		return Settings{}, err
	}

	return settings, nil
}
