// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了ValidationError结构体，表示参数验证失败时的错误信息。
package shellx

import "fmt"

// ValidationError 验证错误
type ValidationError struct {
	Field   string // 字段名
	Message string // 错误信息
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s",
		e.Field, e.Message)
}
