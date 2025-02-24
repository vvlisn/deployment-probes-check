[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)

# Deployment Probes Check

这个 [Kubewarden](https://kubewarden.io) 策略用于验证 Kubernetes Deployments 中的健康检查探针配置。它不仅可以验证必需的探针是否存在，还可以验证探针的时间参数是否合理。

## 功能特性

- 支持验证 liveness、readiness 和 startup 探针的配置
- 可以设置哪些探针是必需的
- 验证探针的时间参数是否合理，包括：
  - initialDelaySeconds
  - timeoutSeconds
  - periodSeconds
  - successThreshold
  - failureThreshold

## 配置说明

The policy settings allow you to specify which probes are required for containers in Deployments:

```yaml
# 所有字段都是可选的
liveness_probe:
  required: true  # 是否要求 liveness 探针
  min_period_seconds: 10  # 最小周期时间（秒）
  max_period_seconds: 300  # 最大周期时间（秒）
  min_timeout_seconds: 1  # 最小超时时间（秒）
  max_timeout_seconds: 60  # 最大超时时间（秒）
readiness_probe:
  required: true  # 是否要求 readiness 探针
  min_period_seconds: 10
  max_period_seconds: 300
  min_timeout_seconds: 1
  max_timeout_seconds: 60
startup_probe:
  required: false  # 是否要求 startup 探针
  min_period_seconds: 10
  max_period_seconds: 300
  min_timeout_seconds: 1
  max_timeout_seconds: 60
```

默认配置：
- Liveness 探针是必需的
- Readiness 探针是必需的
- Startup 探针是可选的

时间参数的默认限制：
- periodSeconds: 10-300 秒
- timeoutSeconds: 1-60 秒

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

111