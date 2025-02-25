[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)

# Deployment Probes Check

这个 [Kubewarden](https://kubewarden.io) 策略用于验证 Kubernetes Deployments 中的健康检查探针配置。它不仅可以验证必需的探针是否存在，还可以验证探针的时间参数是否合理。

## 功能特性

- 支持验证 liveness、readiness 和 startup 探针的配置
- 可以设置哪些探针是必需的
- 验证探针的时间参数是否合理，包括：
  - periodSeconds（探测间隔）
  - timeoutSeconds（探测超时）

## 配置说明

策略配置允许你指定哪些探针是必需的，以及它们的时间参数限制：

```yaml
# 配置示例
liveness_probe:
  required: true  # 是否要求 liveness 探针
  min_period_seconds: 10  # 最小探测间隔（秒）
  max_timeout_seconds: 5  # 最大探测超时（秒）
readiness_probe:
  required: true  # 是否要求 readiness 探针
  min_period_seconds: 10  # 最小探测间隔（秒）
  max_timeout_seconds: 5  # 最大探测超时（秒）
startup_probe:
  required: false  # 是否要求 startup 探针
  min_period_seconds: 10  # 最小探测间隔（秒）
  max_timeout_seconds: 30  # 最大探测超时（秒）
```

默认配置：
- Liveness 探针是必需的
- Readiness 探针是必需的
- Startup 探针是可选的

## 示例

### 接受的 Deployment 配置

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
          periodSeconds: 10    # 符合最小探测间隔要求
          timeoutSeconds: 5    # 符合最大超时要求
        readinessProbe:
          httpGet:
            path: /ready
            port: 80
          periodSeconds: 10
          timeoutSeconds: 5
```

### 拒绝的 Deployment 配置

以下配置会被拒绝，因为缺少必需的探针：

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
        # 缺少必需的 liveness 和 readiness 探针
```

以下配置会被拒绝，因为探针参数不符合要求：

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
          periodSeconds: 5     # 小于最小探测间隔要求
          timeoutSeconds: 10   # 超过最大超时要求
```

## 安装

使用 `kwctl` 安装策略：

```console
kwctl pull ghcr.io/vvhuang-ll/policies/deployment-probes-check:v0.1.0
```

生成策略清单：

```console
kwctl scaffold manifest -t ClusterAdmissionPolicy registry://ghcr.io/vvhuang-ll/policies/deployment-probes-check:v0.1.0
```

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

## License

Apache-2.0
