package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/tidwall/gjson"
)

// validate validates the deployment configuration。
func validate(payload []byte) ([]byte, error) {
	// Parse the validation request。
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

	// Parse the settings。
	settings, settingsErr := NewSettingsFromValidationReq(&validationRequest)
	if settingsErr != nil {
		Logger.ErrorWith("cannot unmarshal settings").
			Err("error", settingsErr).
			Write()
		return kubewarden.RejectRequest(
			kubewarden.Message(fmt.Sprintf("cannot unmarshal settings: %v", settingsErr)),
			kubewarden.Code(http.StatusBadRequest))
	}

	// Validate deployment。
	if err := validateDeployment(validationRequest.Request.Object, settings); err != nil {
		Logger.WarnWith("deployment validation failed").
			Err("error", err).
			Write()
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(http.StatusBadRequest))
	}

	Logger.InfoWith("deployment validation succeeded").Write()
	return kubewarden.AcceptRequest()
}

// validateDeployment validates the deployment configuration。
func validateDeployment(deploymentJSON []byte, settings Settings) error {
	// Validate containers
	containers := gjson.GetBytes(deploymentJSON, "spec.template.spec.containers")
	if !containers.Exists() {
		return errors.New("invalid deployment: missing containers")
	}

	if !containers.IsArray() {
		return errors.New("invalid deployment: containers must be an array")
	}

	if len(containers.Array()) == 0 {
		return errors.New("no containers found in deployment")
	}

	// Validate each container's probes。
	var validationErr error
	containers.ForEach(func(_, container gjson.Result) bool {
		if err := validateContainer(container, settings); err != nil {
			validationErr = err
			return false
		}
		return true
	})

	return validationErr
}

// validateContainer validates a single container's probe configurations。
func validateContainer(container gjson.Result, settings Settings) error {
	containerName := container.Get("name").String()
	if containerName == "" {
		return errors.New("container name is required")
	}

	// Validate liveness probe。
	if err := validateLivenessProbe(container, containerName, settings.LivenessProbe); err != nil {
		return err
	}

	// Validate readiness probe。
	if err := validateReadinessProbe(container, containerName, settings.ReadinessProbe); err != nil {
		return err
	}

	// Validate startup probe。
	if err := validateStartupProbe(container, containerName, settings.StartupProbe); err != nil {
		return err
	}

	return nil
}

// validateLivenessProbe validates the liveness probe configuration。
func validateLivenessProbe(container gjson.Result, containerName string, config ProbeConfig) error {
	if config.Required && !container.Get("livenessProbe").Exists() {
		return fmt.Errorf("container '%s': missing liveness probe", containerName)
	}

	if container.Get("livenessProbe").Exists() {
		periodSeconds := container.Get("livenessProbe.periodSeconds").Int()
		timeoutSeconds := container.Get("livenessProbe.timeoutSeconds").Int()

		if err := validateProbeTimings("liveness", containerName, periodSeconds, timeoutSeconds, config); err != nil {
			return err
		}
	}

	return nil
}

// validateReadinessProbe validates the readiness probe configuration。
func validateReadinessProbe(container gjson.Result, containerName string, config ProbeConfig) error {
	if config.Required && !container.Get("readinessProbe").Exists() {
		return fmt.Errorf("container '%s': missing readiness probe", containerName)
	}

	if container.Get("readinessProbe").Exists() {
		periodSeconds := container.Get("readinessProbe.periodSeconds").Int()
		timeoutSeconds := container.Get("readinessProbe.timeoutSeconds").Int()

		if err := validateProbeTimings("readiness", containerName, periodSeconds, timeoutSeconds, config); err != nil {
			return err
		}
	}

	return nil
}

// validateStartupProbe validates the startup probe configuration。
func validateStartupProbe(container gjson.Result, containerName string, config ProbeConfig) error {
	if config.Required && !container.Get("startupProbe").Exists() {
		return fmt.Errorf("container '%s': missing startup probe", containerName)
	}

	if container.Get("startupProbe").Exists() {
		periodSeconds := container.Get("startupProbe.periodSeconds").Int()
		timeoutSeconds := container.Get("startupProbe.timeoutSeconds").Int()

		if err := validateProbeTimings("startup", containerName, periodSeconds, timeoutSeconds, config); err != nil {
			return err
		}
	}

	return nil
}

// validateProbeTimings validates the timing parameters of a probe。
func validateProbeTimings(probeType string, containerName string, periodSeconds, timeoutSeconds int64, config ProbeConfig) error {
	if config.MinPeriodSeconds > 0 && periodSeconds < int64(config.MinPeriodSeconds) {
		return fmt.Errorf(
			"container '%s': %s probe period (%ds) is less than minimum required (%ds)",
			containerName,
			probeType,
			periodSeconds,
			config.MinPeriodSeconds,
		)
	}

	if config.MaxTimeoutSeconds > 0 && timeoutSeconds > int64(config.MaxTimeoutSeconds) {
		return fmt.Errorf("container '%s': %s probe timeout (%ds) exceeds maximum allowed (%ds)",
			containerName, probeType, timeoutSeconds, config.MaxTimeoutSeconds)
	}

	return nil
}
