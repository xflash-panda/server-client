# 错误处理指南

## 概述

本项目使用自定义的 `APIError` 类型来处理API调用过程中的各种错误。`APIError` 提供了清晰的错误分类，便于快速定位问题。

## 错误类型

### HTTP错误

- `ErrorTypeServerError` (4xx/5xx) - 服务端错误（所有HTTP 4xx和5xx错误都被归类为服务端错误）

### 特殊错误类型

- `ErrorTypeNetworkError` - 网络连接错误（无法连接到服务器）
- `ErrorTypeParseError` - 响应解析错误（服务器返回了无法解析的数据）
- `ErrorTypeNotModified` (304) - 内容未修改（缓存有效）
- `ErrorTypeUnknown` - 未知错误

## 使用方法

### 基本使用

```go
package main

import (
    "fmt"
    "errors"
    
    "github.com/xflash-panda/server-client/pkg"
)

func main() {
    config := &pkg.Config{
        APIHost: "https://api.example.com",
        Token:   "your-token",
    }
    
    client := pkg.New(config)
    
    // 调用API
    users, err := client.Users(1, pkg.VMess)
    if err != nil {
        handleError(err)
        return
    }
    
    fmt.Printf("获取到 %d 个用户\n", len(*users))
}

func handleError(err error) {
    // 检查是否为APIError
    var apiErr *pkg.APIError
    if errors.As(err, &apiErr) {
        // 判断错误类型
        if apiErr.IsServerError() {
            fmt.Printf("服务器错误 [%d]: %s\n", apiErr.StatusCode, apiErr.Message)
            fmt.Println("服务器出现问题，请稍后重试或联系管理员")
        } else if apiErr.IsNetworkError() {
            fmt.Printf("网络错误: %s\n", apiErr.Message)
            fmt.Println("请检查网络连接")
        } else if apiErr.IsParseError() {
            fmt.Printf("解析错误: %s\n", apiErr.Message)
            fmt.Println("服务器返回了无效的数据")
        } else if apiErr.IsNotModified() {
            fmt.Println("数据未修改，可以使用缓存")
        }
        
        // 打印完整错误信息（包含URL）
        fmt.Printf("详细信息: %s\n", apiErr.Error())
    } else {
        // 非APIError类型
        fmt.Printf("未知错误: %v\n", err)
    }
}
```

### 快速判断错误类型

#### 判断是否为服务端错误

```go
var apiErr *pkg.APIError
if errors.As(err, &apiErr) {
    if apiErr.IsServerError() {
        // 4xx/5xx错误 - 服务端问题
        // 所有HTTP错误都被归类为服务端错误
        fmt.Println("服务器出现问题，请稍后重试")
    }
}
```

#### 判断特殊错误类型

```go
if apiErr.IsNetworkError() {
    // 网络连接错误
    fmt.Println("无法连接到服务器，请检查网络")
}

if apiErr.IsParseError() {
    // 响应解析错误
    fmt.Println("服务器返回了无效数据，可能需要更新客户端版本")
}

if apiErr.IsNotModified() {
    // 304 Not Modified - 这不是错误，表示可以使用缓存
    fmt.Println("数据未修改，使用缓存")
}
```

#### 检查具体状态码

```go
switch apiErr.StatusCode {
case 401:
    fmt.Println("认证失败，请检查Token")
case 403:
    fmt.Println("权限不足")
case 404:
    fmt.Println("资源不存在")
case 500:
    fmt.Println("服务器内部错误")
case 503:
    fmt.Println("服务暂时不可用")
}
```

### 重试策略示例

根据错误类型决定是否重试：

```go
import "time"

func callAPIWithRetry(client *pkg.Client) error {
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        users, err := client.Users(1, pkg.VMess)
        if err == nil {
            // 成功
            return nil
        }
        
        var apiErr *pkg.APIError
        if errors.As(err, &apiErr) {
            // 服务端错误或网络错误可以重试
            if apiErr.IsServerError() || apiErr.IsNetworkError() {
                fmt.Printf("请求失败，重试 %d/%d: %s\n", i+1, maxRetries, err)
                time.Sleep(time.Second * time.Duration(i+1)) // 指数退避
                continue
            }
            
            // 304 Not Modified 不是错误，不需要重试
            if apiErr.IsNotModified() {
                return nil // 使用缓存
            }
        }
        
        // 其他未知错误，不重试
        return err
    }
    
    return fmt.Errorf("达到最大重试次数")
}
```

### 错误链支持

`APIError` 实现了 `errors.Unwrap()` 接口，支持Go标准库的错误链功能：

```go
import "syscall"

// 示例：检查是否是特定的底层错误
var apiErr *pkg.APIError
if errors.As(err, &apiErr) {
    // 获取原始错误
    if apiErr.Err != nil {
        fmt.Printf("原始错误: %v\n", apiErr.Err)
    }
    
    // 使用 errors.Is 检查是否是特定错误
    if errors.Is(apiErr, syscall.ECONNREFUSED) {
        fmt.Println("连接被拒绝")
    }
}
```

### 日志记录示例

```go
func logError(err error) {
    var apiErr *pkg.APIError
    if errors.As(err, &apiErr) {
        fields := map[string]interface{}{
            "status_code": apiErr.StatusCode,
            "error_type":  apiErr.Type,
            "message":     apiErr.Message,
            "url":         apiErr.URL,
        }
        
        if apiErr.IsServerError() {
            log.WithFields(fields).Error("服务器错误")
        } else if apiErr.IsNetworkError() {
            log.WithFields(fields).Error("网络错误")
        } else {
            log.WithFields(fields).Error("其他错误")
        }
        
        // 记录原始错误（如果有）
        if apiErr.Err != nil {
            log.WithError(apiErr.Err).Debug("原始错误")
        }
    } else {
        log.WithError(err).Error("未知错误")
    }
}
```

## 错误信息格式

`APIError` 的 `Error()` 方法返回格式化的错误信息：

- 包含URL: `[404] ServerError: resource not found (URL: http://api.example.com/users)`
- 不包含URL: `[500] ServerError: database connection failed`
- 网络错误: `[0] NetworkError: connection timeout (URL: http://api.example.com)`

## 创建自定义错误

### 方式1: 根据状态码自动推断错误类型（推荐）

```go
// 所有HTTP 4xx/5xx错误都被归类为服务端错误
err := pkg.NewAPIErrorFromStatusCode(404, "user not found", "http://api.example.com/users/123", nil)
// 结果: Type = ErrorTypeServerError

err := pkg.NewAPIErrorFromStatusCode(500, "database error", "http://api.example.com/data", nil)
// 结果: Type = ErrorTypeServerError
```

### 方式2: 使用预定义工厂函数

```go
// 创建网络错误
err := pkg.NewNetworkError("connection timeout", "http://api.example.com", originalErr)

// 创建解析错误
err := pkg.NewParseError("invalid JSON response", originalErr)

// 创建304未修改错误
err := pkg.NewNotModifiedError()
```

### 方式3: 完全自定义

```go
err := pkg.NewAPIError(
    418,                        // 状态码
    pkg.ErrorTypeServerError,   // 错误类型
    "I'm a teapot",            // 消息
    "http://api.example.com",  // URL
    nil,                       // 原始错误
)
```

## 最佳实践

1. **区分错误类型来决定处理策略**:
   - **服务端错误 (4xx/5xx)**: 所有HTTP错误都被归类为服务端错误，可以重试
   - **网络错误**: 网络问题，可以重试
   - **解析错误**: 数据格式问题，记录日志并检查API版本

2. **使用 `errors.As()` 安全地转换错误**:
   ```go
   var apiErr *pkg.APIError
   if errors.As(err, &apiErr) {
       // 安全地访问 APIError 的字段和方法
   }
   ```

3. **记录详细日志**: 利用 `APIError` 的字段记录完整的错误上下文

4. **处理 304 Not Modified**: 这不是真正的错误，表示可以使用缓存数据

5. **保留错误链**: 创建错误时传递原始错误，便于深层调试

## API错误字段说明

```go
type APIError struct {
    StatusCode int       // HTTP状态码（0表示非HTTP错误，如网络错误）
    Type       ErrorType // 错误类型（ClientError/ServerError等）
    Message    string    // 人类可读的错误消息
    URL        string    // 发生错误的请求URL
    Err        error     // 原始错误（可选，用于错误链）
}
```

## 错误类型决策树

```
是否能连接到服务器？
├─ 否 → ErrorTypeNetworkError
└─ 是
    └─ 是否能解析响应？
        ├─ 否 → ErrorTypeParseError
        └─ 是
            └─ HTTP状态码
                ├─ 304 → ErrorTypeNotModified
                ├─ 4xx/5xx → ErrorTypeServerError
                └─ 其他 → ErrorTypeUnknown
```

