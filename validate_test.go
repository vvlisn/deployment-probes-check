package main

import (
	"encoding/json"
	"testing"

	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

func TestValidateDeploymentProbes(t *testing.T) {
	tests := []struct {
		name        string
		settings    string
		deployment  string
		shouldAllow bool
	}{
		{
			name: "reject deployment with invalid probe period",
			settings: `{
				"liveness_probe": {"required": true, "min_period_seconds": 10},
				"readiness_probe": {"required": true}
			}`,
			deployment: `{
				"apiVersion": "apps/v1",
				"kind": "Deployment",
				"spec": {
					"template": {
						"spec": {
							"containers": [
								{
									"name": "test-container",
									"livenessProbe": {
										"httpGet": {
											"path": "/healthz",
											"port": 8080
										},
										"periodSeconds": 5
									},
									"readinessProbe": {
										"httpGet": {
											"path": "/ready",
											"port": 8080
										}
									}
								}
							]
						}
					}
				}
			}`,
			shouldAllow: false,
		},
		{
			name: "reject deployment with invalid probe timeout",
			settings: `{
				"liveness_probe": {"required": true, "max_timeout_seconds": 5},
				"readiness_probe": {"required": true}
			}`,
			deployment: `{
				"apiVersion": "apps/v1",
				"kind": "Deployment",
				"spec": {
					"template": {
						"spec": {
							"containers": [
								{
									"name": "test-container",
									"livenessProbe": {
										"httpGet": {
											"path": "/healthz",
											"port": 8080
										},
										"timeoutSeconds": 10
									},
									"readinessProbe": {
										"httpGet": {
											"path": "/ready",
											"port": 8080
										}
									}
								}
							]
						}
					}
				}
			}`,
			shouldAllow: false,
		},
		{
			name: "accept deployment with valid probe configurations",
			settings: `{
				"liveness_probe": {"required": true},
				"readiness_probe": {"required": true}
			}`,
			deployment: `{
				"apiVersion": "apps/v1",
				"kind": "Deployment",
				"spec": {
					"template": {
						"spec": {
							"containers": [
								{
									"name": "test-container",
									"livenessProbe": {
										"httpGet": {
											"path": "/healthz",
											"port": 8080
										}
									},
									"readinessProbe": {
										"httpGet": {
											"path": "/ready",
											"port": 8080
										}
									}
								}
							]
						}
					}
				}
			}`,
			shouldAllow: true,
		},
		{
			name: "reject deployment with missing required probes",
			settings: `{
				"liveness_probe": {"required": true},
				"readiness_probe": {"required": true}
			}`,
			deployment: `{
				"apiVersion": "apps/v1",
				"kind": "Deployment",
				"spec": {
					"template": {
						"spec": {
							"containers": [
								{
									"name": "test-container"
								}
							]
						}
					}
				}
			}`,
			shouldAllow: false,
		},
		{
			name: "accept deployment with optional probes",
			settings: `{
				"liveness_probe": {"required": false},
				"readiness_probe": {"required": false}
			}`,
			deployment: `{
				"apiVersion": "apps/v1",
				"kind": "Deployment",
				"spec": {
					"template": {
						"spec": {
							"containers": [
								{
									"name": "test-container"
								}
							]
						}
					}
				}
			}`,
			shouldAllow: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := kubewarden_protocol.ValidationRequest{
				Request: kubewarden_protocol.KubernetesAdmissionRequest{
					Object: json.RawMessage(test.deployment),
				},
				Settings: json.RawMessage(test.settings),
			}

			payload, err := json.Marshal(request)
			if err != nil {
				t.Errorf("Unexpected error: %+v", err)
			}

			responsePayload, err := validate(payload)
			if err != nil {
				t.Errorf("Unexpected error: %+v", err)
			}

			var response kubewarden_protocol.ValidationResponse
			if err := json.Unmarshal(responsePayload, &response); err != nil {
				t.Errorf("Unexpected error: %+v", err)
			}

			if response.Accepted != test.shouldAllow {
				t.Errorf("Expected validation to return %v, got %v. Message: %s",
					test.shouldAllow, response.Accepted, *response.Message)
			}
		})
	}
}
