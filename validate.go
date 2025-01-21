package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/tidwall/gjson"
)

// validate validates the deployment configuration
func validate(payload []byte) ([]byte, error) {
	// Parse the validation request
	validationRequest := kubewarden_protocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		Logger.ErrorWith("cannot unmarshal validation request").
			Err("error", err).
			Write()
		return kubewarden.RejectRequest(
			kubewarden.Message(fmt.Sprintf("cannot unmarshal validation request: %v", err)),
			kubewarden.Code(http.StatusBadRequest))
	}

	// Parse the settings
	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		Logger.ErrorWith("cannot unmarshal settings").
			Err("error", err).
			Write()
		return kubewarden.RejectRequest(
			kubewarden.Message(fmt.Sprintf("cannot unmarshal settings: %v", err)),
			kubewarden.Code(http.StatusBadRequest))
	}

	// Access the raw JSON that describes the object
	deploymentJSON := validationRequest.Request.Object

	// Validate containers
	containers := gjson.GetBytes(deploymentJSON, "spec.template.spec.containers")
	if !containers.Exists() {
		return kubewarden.RejectRequest(
			kubewarden.Message("invalid deployment: missing containers"),
			kubewarden.Code(http.StatusBadRequest))
	}

	if !containers.IsArray() {
		return kubewarden.RejectRequest(
			kubewarden.Message("invalid deployment: containers must be an array"),
			kubewarden.Code(http.StatusBadRequest))
	}

	if len(containers.Array()) == 0 {
		return kubewarden.RejectRequest(
			kubewarden.Message("no containers found in deployment"),
			kubewarden.Code(http.StatusBadRequest))
	}

	// Validate each container's probes
	var validationErr error
	containers.ForEach(func(_, container gjson.Result) bool {
		containerName := container.Get("name").String()
		if containerName == "" {
			validationErr = fmt.Errorf("container name is required")
			return false
		}

		// Validate liveness probe
		if settings.LivenessProbe.Required && !container.Get("livenessProbe").Exists() {
			validationErr = fmt.Errorf("container '%s': missing liveness probe", containerName)
			return false
		}

		// Validate readiness probe
		if settings.ReadinessProbe.Required && !container.Get("readinessProbe").Exists() {
			validationErr = fmt.Errorf("container '%s': missing readiness probe", containerName)
			return false
		}

		// Validate startup probe
		if settings.StartupProbe.Required && !container.Get("startupProbe").Exists() {
			validationErr = fmt.Errorf("container '%s': missing startup probe", containerName)
			return false
		}

		return true
	})

	if validationErr != nil {
		Logger.WarnWith("deployment validation failed").
			Err("error", validationErr).
			Write()
		return kubewarden.RejectRequest(
			kubewarden.Message(validationErr.Error()),
			kubewarden.Code(http.StatusBadRequest))
	}

	Logger.InfoWith("deployment validation succeeded").Write()
	return kubewarden.AcceptRequest()
}
