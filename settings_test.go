package main

import (
	"encoding/json"
	"testing"
)

func TestParsingSettingsWithNoValueProvided(t *testing.T) {
	rawSettings := []byte(`{}`)
	settings := Settings{}
	err := json.Unmarshal(rawSettings, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if settings.LivenessProbe.Required {
		t.Error("Expected LivenessProbe.Required to be false")
	}
	if !settings.ReadinessProbe.Required {
		t.Error("Expected ReadinessProbe.Required to be true")
	}
	if settings.StartupProbe.Required {
		t.Error("Expected StartupProbe.Required to be false")
	}
}

func TestParsingSettingsWithProbeConfig(t *testing.T) {
	tests := []struct {
		name     string
		settings string
		expected Settings
	}{
		{
			name: "both probes required",
			settings: `{
				"liveness_probe": {"required": true},
				"readiness_probe": {"required": true}
			}`,
			expected: Settings{
				LivenessProbe: ProbeConfig{
					Required: true,
				},
				ReadinessProbe: ProbeConfig{
					Required: true,
				},
				StartupProbe: ProbeConfig{
					Required: false,
				},
			},
		},
		{
			name: "only liveness probe required",
			settings: `{
				"liveness_probe": {"required": true},
				"readiness_probe": {"required": false}
			}`,
			expected: Settings{
				LivenessProbe: ProbeConfig{
					Required: true,
				},
				ReadinessProbe: ProbeConfig{
					Required: false,
				},
				StartupProbe: ProbeConfig{
					Required: false,
				},
			},
		},
		{
			name: "no probes required",
			settings: `{
				"liveness_probe": {"required": false},
				"readiness_probe": {"required": false}
			}`,
			expected: Settings{
				LivenessProbe: ProbeConfig{
					Required: false,
				},
				ReadinessProbe: ProbeConfig{
					Required: false,
				},
				StartupProbe: ProbeConfig{
					Required: false,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			settings := Settings{}
			err := json.Unmarshal([]byte(test.settings), &settings)
			if err != nil {
				t.Errorf("Unexpected error: %+v", err)
			}

			if settings.LivenessProbe.Required != test.expected.LivenessProbe.Required {
				t.Errorf("Expected LivenessProbe.Required to be %v, got %v",
					test.expected.LivenessProbe.Required, settings.LivenessProbe.Required)
			}

			if settings.ReadinessProbe.Required != test.expected.ReadinessProbe.Required {
				t.Errorf("Expected ReadinessProbe.Required to be %v, got %v",
					test.expected.ReadinessProbe.Required, settings.ReadinessProbe.Required)
			}

			if settings.StartupProbe.Required != test.expected.StartupProbe.Required {
				t.Errorf("Expected StartupProbe.Required to be %v, got %v",
					test.expected.StartupProbe.Required, settings.StartupProbe.Required)
			}
		})
	}
}

func TestValidateMethod(t *testing.T) {
	tests := []struct {
		name     string
		settings Settings
		isValid  bool
	}{
		{
			name: "default settings",
			settings: Settings{
				LivenessProbe: ProbeConfig{
					Required: true,
				},
				ReadinessProbe: ProbeConfig{
					Required: true,
				},
			},
			isValid: true,
		},
		{
			name: "no probes required",
			settings: Settings{
				LivenessProbe: ProbeConfig{
					Required: false,
				},
				ReadinessProbe: ProbeConfig{
					Required: false,
				},
			},
			isValid: true,
		},
		{
			name: "only liveness probe required",
			settings: Settings{
				LivenessProbe: ProbeConfig{
					Required: true,
				},
				ReadinessProbe: ProbeConfig{
					Required: false,
				},
			},
			isValid: true,
		},
		{
			name: "only readiness probe required",
			settings: Settings{
				LivenessProbe: ProbeConfig{
					Required: false,
				},
				ReadinessProbe: ProbeConfig{
					Required: true,
				},
			},
			isValid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.settings.Validate()
			if test.isValid && err != nil {
				t.Errorf("Expected settings to be valid, got error: %v", err)
			}
			if !test.isValid && err == nil {
				t.Error("Expected settings to be invalid, but got no error")
			}
		})
	}
}

func TestValidateProbeTimeSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings Settings
		isValid  bool
	}{
		{
			name: "valid time settings",
			settings: Settings{
				LivenessProbe: ProbeConfig{
					Required:          true,
					MinPeriodSeconds:  30,
					MaxTimeoutSeconds: 5,
				},
			},
			isValid: true,
		},
		{
			name: "negative min period seconds",
			settings: Settings{
				LivenessProbe: ProbeConfig{
					Required:          true,
					MinPeriodSeconds:  -10,
					MaxTimeoutSeconds: 5,
				},
			},
			isValid: false,
		},
		{
			name: "negative max timeout seconds",
			settings: Settings{
				LivenessProbe: ProbeConfig{
					Required:          true,
					MinPeriodSeconds:  30,
					MaxTimeoutSeconds: -5,
				},
			},
			isValid: false,
		},
		{
			name: "min period less than max timeout",
			settings: Settings{
				LivenessProbe: ProbeConfig{
					Required:          true,
					MinPeriodSeconds:  5,
					MaxTimeoutSeconds: 10,
				},
			},
			isValid: false,
		},
		{
			name: "zero time settings",
			settings: Settings{
				LivenessProbe: ProbeConfig{
					Required:          true,
					MinPeriodSeconds:  0,
					MaxTimeoutSeconds: 0,
				},
			},
			isValid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.settings.Validate()
			if test.isValid && err != nil {
				t.Errorf("Expected settings to be valid, got error: %v", err)
			}
			if !test.isValid && err == nil {
				t.Error("Expected settings to be invalid, but got no error")
			}
		})
	}
}
