// Package shellx 提供了一个功能完善、易于使用的Go语言shell命令执行库。
//
// 本库支持同步和异步命令执行、命令管道、输入输出重定向、超时控制、
// 上下文管理等功能，并提供类型安全的API和友好的链式调用接口。
//
// 主要特性：
//   - 支持数组和字符串两种命令创建方式
//   - 链式调用API，易于使用
//   - 完整的错误处理和类型安全
//   - 支持多种shell类型（bash、sh、powershell等）
//   - 异步执行和命令管道支持
//   - 跨平台兼容
//
// 基本用法：
//
//	import "gitee.com/MM-Q/shellx/shellx"
//
//	// 创建命令
//	cmd := shellx.NewCmd("ls", "-la").
//		WithWorkDir("/tmp").
//		WithTimeout(30 * time.Second).
//		Build()
//
//	// 执行命令
//	executor := NewExecutor()
//	result, err := executor.Exec(cmd)
package shellx
