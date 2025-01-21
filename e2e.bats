#!/usr/bin/env bats

@test "accept deployment with valid probe configurations" {
  run kwctl run annotated-policy.wasm \
    -r test_data/deployment-valid.json \
    --settings-json '{"liveness_probe": {"required": true}, "readiness_probe": {"required": true}}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request is accepted
  [ "$status" -eq 0 ]
  [[ "$output" =~ "deployment validation succeeded" ]]
}

@test "reject deployment with missing required probes" {
  run kwctl run annotated-policy.wasm \
    -r test_data/deployment-missing-probes.json \
    --settings-json '{"liveness_probe": {"required": true}, "readiness_probe": {"required": true}}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request is rejected
  [ "$status" -eq 0 ]
  [[ "$output" =~ "container 'test-container': missing liveness probe" ]]
}

@test "accept deployment with optional probes" {
  run kwctl run annotated-policy.wasm \
    -r test_data/deployment-missing-probes.json \
    --settings-json '{"liveness_probe": {"required": false}, "readiness_probe": {"required": false}}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request is accepted
  [ "$status" -eq 0 ]
  [[ "$output" =~ "deployment validation succeeded" ]]
}

@test "accept settings with no probes required" {
  run kwctl run annotated-policy.wasm \
    -r test_data/deployment-valid.json \
    --settings-json '{"liveness_probe": {"required": false}, "readiness_probe": {"required": false}, "startup_probe": {"required": false}}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # settings validation succeeded
  [ "$status" -eq 0 ]
  [[ "$output" =~ "deployment validation succeeded" ]]
}

@test "accept settings validation" {
  run kwctl run annotated-policy.wasm \
    -r test_data/deployment-valid.json \
    --settings-json '{"liveness_probe": {"required": true}}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # settings validation succeeded
  [ "$status" -eq 0 ]
  [[ "$output" =~ "deployment validation succeeded" ]]
}
