package main

import (
	"errors"
	"fmt"

	"github.com/xflash-panda/server-client/pkg"
)

func main() {
	// 模拟不同类型的错误
	demonstrateErrors()
}

func demonstrateErrors() {
	fmt.Println("=== 错误处理示例 ===\n")

	// 示例1: 客户端错误 (401)
	err1 := pkg.NewAPIErrorFromStatusCode(401, "invalid token", "http://api.example.com/users", nil)
	handleError("示例1: 401 Unauthorized", err1)

	// 示例2: 客户端错误 (404)
	err2 := pkg.NewAPIErrorFromStatusCode(404, "user not found", "http://api.example.com/users/123", nil)
	handleError("示例2: 404 Not Found", err2)

	// 示例3: 服务端错误 (500)
	err3 := pkg.NewAPIErrorFromStatusCode(500, "database connection failed", "http://api.example.com/data", nil)
	handleError("示例3: 500 Internal Server Error", err3)

	// 示例4: 服务端错误 (503)
	err4 := pkg.NewAPIErrorFromStatusCode(503, "service temporarily unavailable", "http://api.example.com/service", nil)
	handleError("示例4: 503 Service Unavailable", err4)

	// 示例5: 网络错误
	err5 := pkg.NewNetworkError("connection timeout", "http://api.example.com", nil)
	handleError("示例5: 网络错误", err5)

	// 示例6: 解析错误
	err6 := pkg.NewParseError("invalid JSON response", nil)
	handleError("示例6: 解析错误", err6)

	// 示例7: 304 Not Modified
	err7 := pkg.NewNotModifiedError()
	handleError("示例7: 304 Not Modified", err7)

	// 示例8: 错误链
	originalErr := errors.New("connection refused")
	err8 := pkg.NewNetworkError("无法连接到服务器", "http://api.example.com", originalErr)
	handleErrorWithChain("示例8: 带错误链的网络错误", err8)
}

func handleError(title string, err error) {
	fmt.Printf("--- %s ---\n", title)

	var apiErr *pkg.APIError
	if errors.As(err, &apiErr) {
		// 输出基本信息
		fmt.Printf("错误信息: %s\n", apiErr.Error())
		fmt.Printf("状态码: %d\n", apiErr.StatusCode)
		fmt.Printf("错误类型: %s\n", apiErr.Type)
		fmt.Printf("错误消息: %s\n", apiErr.Message)

		// 判断错误类别
		if apiErr.IsClientError() {
			fmt.Println("✗ 这是客户端错误 (4xx)")
			fmt.Println("  建议: 检查请求参数、认证信息或权限")
			fmt.Println("  是否重试: 否")
		} else if apiErr.IsServerError() {
			fmt.Println("✗ 这是服务端错误 (5xx)")
			fmt.Println("  建议: 稍后重试或联系管理员")
			fmt.Println("  是否重试: 是")
		} else if apiErr.IsNetworkError() {
			fmt.Println("✗ 这是网络错误")
			fmt.Println("  建议: 检查网络连接")
			fmt.Println("  是否重试: 是")
		} else if apiErr.IsParseError() {
			fmt.Println("✗ 这是解析错误")
			fmt.Println("  建议: 检查API版本兼容性")
			fmt.Println("  是否重试: 否")
		} else if apiErr.IsNotModified() {
			fmt.Println("ℹ 数据未修改")
			fmt.Println("  建议: 使用缓存数据")
			fmt.Println("  是否重试: 否")
		}
	}

	fmt.Println()
}

func handleErrorWithChain(title string, err error) {
	fmt.Printf("--- %s ---\n", title)

	var apiErr *pkg.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("错误信息: %s\n", apiErr.Error())

		// 检查错误链
		if apiErr.Err != nil {
			fmt.Printf("原始错误: %v\n", apiErr.Err)
			fmt.Println("✓ 包含错误链，便于深层调试")
		}
	}

	fmt.Println()
}
