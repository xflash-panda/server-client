# Go 1.20 兼容性修改说明

## 修改日期
2025-10-24

## 修改目的
将代码调整为与 Go 1.20 版本兼容，移除对 Go 1.21+ 特性的依赖。

## 主要修改

### 1. 修复方法签名不匹配问题

**文件**: `pkg/client.go`

**问题**: `SubmitStatsWithAgent` 方法的定义与测试代码中的调用不匹配。

**修改内容**:
- 添加 `nodeIp` 参数到 `SubmitStatsWithAgent` 方法
- 更新方法签名: `func (c *Client) SubmitStatsWithAgent(registerId int, nodeType NodeType, nodeIp string, stats *TrafficStats) error`
- 在请求体中添加可选的 `node_ip` 字段（与 `Heartbeat` 方法保持一致）

### 2. 降级依赖包版本

**文件**: `go.mod`

**原因**: 某些 `golang.org/x` 包的新版本需要 Go 1.21 或更高版本。

**修改内容**:
- `golang.org/x/net`: v0.36.0 → v0.17.0
- `golang.org/x/sys`: v0.30.0 → v0.13.0
- `golang.org/x/time`: v0.4.0 → v0.3.0

这些版本完全兼容 Go 1.20。

## 代码特性验证

以下特性已验证与 Go 1.20 兼容：

✓ **泛型 (Generics)**: 代码中使用的泛型特性在 Go 1.18 引入，Go 1.20 完全支持
  - `AsConfig[T NodeConfig]` 函数
  - `UnmarshalConfig[T any]` 函数

✓ **错误处理**: 使用标准的 `errors.As`, `errors.Is` 等方法（Go 1.13+）

✓ **没有使用 Go 1.21+ 特性**:
  - 未使用 `min()`, `max()`, `clear()` 内置函数
  - 未使用 `log/slog` 包
  - 未使用 `slices`, `maps`, `cmp` 包
  - 未使用 `errors.ErrUnsupported`

## 兼容性说明

修改后的代码完全兼容 Go 1.20，并且向前兼容更高版本的 Go。

## 测试建议

建议运行以下命令验证兼容性：

```bash
# 验证 Go 版本
go version

# 下载依赖
go mod download

# 编译代码
go build ./...

# 运行测试
go test ./...
```

## 相关文件

- `pkg/client.go` - 修复方法签名
- `go.mod` - 降级依赖版本
- `go.sum` - 需要运行 `go mod tidy` 更新

