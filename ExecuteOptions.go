// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了ExecuteOptions结构体，用于配置命令执行选项。
package shellx

// ExecuteOptions 执行选项配置
type ExecuteOptions struct {
	Shell   ShellType // 指定shell类型
	Capture bool      // 是否捕获输出到Result
}
