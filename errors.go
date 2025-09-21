// Package shellx 定义了命令执行相关的错误类型和处理函数
package shellx

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// ErrorType 定义错误类型枚举
type ErrorType int

const (
	ErrorUnknown   ErrorType = iota // 未知错误
	ErrorTimeout                    // 超时错误
	ErrorCanceled                   // 取消错误
	ErrorExecution                  // 执行错误（命令执行失败）
	ErrorSystem                     // 系统错误
)

// 预定义的错误变量
var (
	// ErrTimeout 表示命令执行超时
	ErrTimeout = errors.New("command execution timeout")
	// ErrCanceled 表示命令被取消
	ErrCanceled = errors.New("command execution canceled")
	// ErrAlreadyExecuted 表示命令已经执行过
	ErrAlreadyExecuted = errors.New("command has already been executed")
	// ErrNotStarted 表示命令尚未启动
	ErrNotStarted = errors.New("command has not been started")
	// ErrNoProcess 表示没有进程可操作
	ErrNoProcess = errors.New("no process to operate")
)

// CommandError 包装命令执行错误，提供详细的错误信息和类型判断
type CommandError struct {
	Err             error         // 原始错误
	Type            ErrorType     // 错误类型
	ExitCode        int           // 命令退出码
	TimeoutDuration time.Duration // 设置的超时时间
}

// Error 实现 error 接口，返回格式化的错误信息
func (e *CommandError) Error() string {
	switch e.Type {
	case ErrorTimeout:
		return fmt.Sprintf("command execution timeout: exceeded %v", e.TimeoutDuration)
	case ErrorCanceled:
		return "command execution canceled"
	case ErrorExecution:
		return fmt.Sprintf("command execution failed with exit code: %d", e.ExitCode)
	case ErrorSystem:
		return fmt.Sprintf("system error: %v", e.Err)
	default:
		return e.Err.Error()
	}
}

// Unwrap 实现错误解包，支持 errors.Is 和 errors.As
func (e *CommandError) Unwrap() error {
	return e.Err
}

// IsTimeoutError 判断是否为超时错误
func IsTimeoutError(err error) bool {
	var cmdErr *CommandError
	if errors.As(err, &cmdErr) {
		return cmdErr.Type == ErrorTimeout
	}
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrTimeout)
}

// IsCanceledError 判断是否为取消错误
func IsCanceledError(err error) bool {
	var cmdErr *CommandError
	if errors.As(err, &cmdErr) {
		return cmdErr.Type == ErrorCanceled
	}
	return errors.Is(err, context.Canceled) || errors.Is(err, ErrCanceled)
}

// IsExecutionError 判断是否为命令执行错误
func IsExecutionError(err error) bool {
	var cmdErr *CommandError
	if errors.As(err, &cmdErr) {
		return cmdErr.Type == ErrorExecution
	}
	_, ok := err.(*exec.ExitError)
	return ok
}

// GetExitCode 获取命令退出码，如果不是执行错误则返回 -1
func GetExitCode(err error) int {
	var cmdErr *CommandError
	if errors.As(err, &cmdErr) {
		return cmdErr.ExitCode
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}

	return -1
}

// classifyError 分类错误并包装为 CommandError
func classifyError(err error, timeoutDuration time.Duration) error {
	if err == nil {
		return nil
	}

	cmdErr := &CommandError{
		Err:             err,
		Type:            ErrorUnknown,
		ExitCode:        -1,
		TimeoutDuration: timeoutDuration,
	}

	// 检查是否为超时错误
	if errors.Is(err, context.DeadlineExceeded) {
		cmdErr.Type = ErrorTimeout
		return cmdErr
	}

	// 检查是否为取消错误
	if errors.Is(err, context.Canceled) {
		cmdErr.Type = ErrorCanceled
		return cmdErr
	}

	// 检查是否为命令执行错误
	if exitErr, ok := err.(*exec.ExitError); ok {
		cmdErr.Type = ErrorExecution
		cmdErr.ExitCode = exitErr.ExitCode()
		return cmdErr
	}

	// 其他系统错误
	cmdErr.Type = ErrorSystem
	return cmdErr
}
