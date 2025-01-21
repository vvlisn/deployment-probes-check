[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)

# Deployment Probes Check

This [Kubewarden](https://kubewarden.io) policy validates the health check probe configurations in Kubernetes Deployments.

## Settings

The policy settings allow you to specify which probes are required for containers in Deployments:

```yaml
# All fields are optional
liveness_probe:
  required: true  # Whether liveness probe is required
readiness_probe:
  required: true  # Whether readiness probe is required
startup_probe:
  required: false # Whether startup probe is required
```

By default:
- Liveness probe is required
- Readiness probe is required
- Startup probe is optional

## Examples

### Accept a Deployment with valid probe configurations

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        livenessProbe:
          httpGet:
            path: /healthz
            port: 80
        readinessProbe:
          httpGet:
            path: /ready
            port: 80
```

### Reject a Deployment with missing required probes

The policy will reject Deployments that are missing required probes:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        # Missing required liveness and readiness probes
```

## Installation

You can install the policy using `kwctl`:

```console
kwctl pull ghcr.io/vvhuang-ll/policies/deployment-probes-check:v0.1.0
```

Then, you can generate the policy manifest:

```console
kwctl scaffold manifest -t ClusterAdmissionPolicy registry://ghcr.io/vvhuang-ll/policies/deployment-probes-check:v0.1.0
```

## License

Apache-2.0

## 开发

### 构建

```bash
make build
```

### 测试

运行单元测试：
```bash
make test
```

运行端到端测试：
```bash
make e2e-tests
```
