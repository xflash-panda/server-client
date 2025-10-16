package pkg

import (
	"fmt"
	"net/http"
)

// ErrorType 定义错误类型
type ErrorType string

const (
	// HTTP错误
	ErrorTypeServerError ErrorType = "ServerError" // 5xx 服务端错误

	// 特殊错误类型
	ErrorTypeNetworkError ErrorType = "NetworkError" // 网络连接错误
	ErrorTypeParseError   ErrorType = "ParseError"   // 响应解析错误
	ErrorTypeNotModified  ErrorType = "NotModified"  // 304 Not Modified
	ErrorTypeUnknown      ErrorType = "Unknown"      // 未知错误
)

// APIError 自定义API错误类型
type APIError struct {
	StatusCode int       // HTTP状态码
	Type       ErrorType // 错误类型
	Message    string    // 错误消息
	URL        string    // 请求的URL
	Err        error     // 原始错误
}

// Error 实现error接口
func (e *APIError) Error() string {
	if e.URL != "" {
		return fmt.Sprintf("[%d] %s: %s (URL: %s)", e.StatusCode, e.Type, e.Message, e.URL)
	}
	return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.Type, e.Message)
}

// Unwrap 实现errors.Unwrap接口
func (e *APIError) Unwrap() error {
	return e.Err
}

// IsServerError 判断是否为服务端错误 (4xx/5xx)
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 600
}

// IsNetworkError 判断是否为网络错误
func (e *APIError) IsNetworkError() bool {
	return e.Type == ErrorTypeNetworkError
}

// IsParseError 判断是否为解析错误
func (e *APIError) IsParseError() bool {
	return e.Type == ErrorTypeParseError
}

// IsNotModified 判断是否为304未修改
func (e *APIError) IsNotModified() bool {
	return e.StatusCode == http.StatusNotModified || e.Type == ErrorTypeNotModified
}

// NewAPIError 创建一个新的API错误
func NewAPIError(statusCode int, errorType ErrorType, message string, url string, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Type:       errorType,
		Message:    message,
		URL:        url,
		Err:        err,
	}
}

// NewAPIErrorFromStatusCode 根据状态码自动推断错误类型
func NewAPIErrorFromStatusCode(statusCode int, message string, url string, err error) *APIError {
	errorType := getErrorTypeFromStatusCode(statusCode)
	return NewAPIError(statusCode, errorType, message, url, err)
}

// getErrorTypeFromStatusCode 根据HTTP状态码获取对应的错误类型
func getErrorTypeFromStatusCode(statusCode int) ErrorType {
	if statusCode == http.StatusNotModified {
		return ErrorTypeNotModified
	}
	if statusCode >= 400 && statusCode < 600 {
		return ErrorTypeServerError
	}
	return ErrorTypeUnknown
}

// 预定义的常见错误

// NewNetworkError 创建网络错误
func NewNetworkError(message string, url string, err error) *APIError {
	return NewAPIError(0, ErrorTypeNetworkError, message, url, err)
}

// NewParseError 创建解析错误
func NewParseError(message string, err error) *APIError {
	return NewAPIError(0, ErrorTypeParseError, message, "", err)
}

// NewNotModifiedError 创建304未修改错误
func NewNotModifiedError() *APIError {
	return NewAPIError(http.StatusNotModified, ErrorTypeNotModified, "content not modified", "", nil)
}

// NewBusinessLogicError 创建业务逻辑错误
// 业务逻辑错误通常来自API响应中的Message字段，默认视为服务端错误(500)
func NewBusinessLogicError(message string, url string) *APIError {
	return NewAPIError(http.StatusInternalServerError, ErrorTypeServerError, message, url, nil)
}
