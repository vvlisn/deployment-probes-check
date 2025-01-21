[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)

# Deployment Probes Check Policy

这是一个 Kubewarden 策略，用于验证 Kubernetes Deployment 中的容器探针配置。该策略可以确保容器配置了适当的健康检查机制。

## 功能特点

- 支持三种类型的探针检查：
  - 存活探针（Liveness Probe）
  - 就绪探针（Readiness Probe）
  - 启动探针（Startup Probe）
- 可以为每种探针类型配置：
  - 是否必需
  - 最小探测周期
  - 最大超时时间
- 提供详细的验证错误信息
- 支持灵活的配置选项

## 配置说明

策略配置使用 JSON 格式，支持以下选项：

```json
{
  "liveness_probe": {
    "required": true,
    "minimum_period_seconds": 5,
    "maximum_timeout_seconds": 3
  },
  "readiness_probe": {
    "required": true,
    "minimum_period_seconds": 3,
    "maximum_timeout_seconds": 2
  },
  "startup_probe": {
    "required": false,
    "minimum_period_seconds": 20,
    "maximum_timeout_seconds": 5
  }
}
```

### 配置选项说明

每种探针类型（`liveness_probe`、`readiness_probe`、`startup_probe`）都支持以下配置：

- `required`：布尔值，指定是否要求配置该类型的探针
  - 默认值：
    - `liveness_probe`: `true`
    - `readiness_probe`: `true`
    - `startup_probe`: `false`

- `minimum_period_seconds`：整数，指定探测周期的最小值（可选）
  - 如果设置，必须大于 0
  - 必须大于 `maximum_timeout_seconds`

- `maximum_timeout_seconds`：整数，指定探测超时时间的最大值（可选）
  - 如果设置，必须大于 0
  - 必须小于 `minimum_period_seconds`

### 配置规则

1. 至少需要启用一种探针检查（`required` 为 `true`）
2. 如果设置了 `minimum_period_seconds`，容器的探针周期必须大于或等于这个值
3. 如果设置了 `maximum_timeout_seconds`，容器的探针超时时间必须小于或等于这个值
4. 探针的周期必须大于其超时时间

## 配置示例

### 基本配置

只要求存活探针：

```json
{
  "liveness_probe": {
    "required": true
  },
  "readiness_probe": {
    "required": false
  }
}
```

### 完整配置

要求所有探针，并设置时间限制：

```json
{
  "liveness_probe": {
    "required": true,
    "minimum_period_seconds": 5,
    "maximum_timeout_seconds": 3
  },
  "readiness_probe": {
    "required": true,
    "minimum_period_seconds": 3,
    "maximum_timeout_seconds": 2
  },
  "startup_probe": {
    "required": true,
    "minimum_period_seconds": 20,
    "maximum_timeout_seconds": 5
  }
}
```

## 验证示例

### 有效的 Deployment 配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example
spec:
  template:
    spec:
      containers:
      - name: example
        image: example:latest
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          periodSeconds: 10
          timeoutSeconds: 1
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          periodSeconds: 5
          timeoutSeconds: 1
        startupProbe:
          httpGet:
            path: /startup
            port: 8080
          periodSeconds: 30
          timeoutSeconds: 1
```

### 常见错误

1. 缺少必需的探针：
   ```
   container 'example': missing liveness probe
   ```

2. 探针周期太短：
   ```
   container 'example': invalid liveness probe: period seconds 3 is less than minimum required 5
   ```

3. 探针超时时间太长：
   ```
   container 'example': invalid liveness probe: timeout seconds 5 is greater than maximum allowed 3
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

## 许可证

Apache License 2.0 - 查看 [LICENSE](LICENSE) 文件了解详情。
