package shellx

import (
	"fmt"
	"time"
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string // 字段名
	Message string // 错误信息
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s",
		e.Field, e.Message)
}

// ExecutionError 执行错误类型
type ExecutionError struct {
	Cmd       *Command
	ExitCode  int
	Stderr    string
	Err       error
	Timestamp time.Time
}

func (e *ExecutionError) Error() string {
	cmdStr := e.Cmd.Name()
	if e.Cmd.Raw() != "" {
		cmdStr = e.Cmd.Raw()
	}
	return fmt.Sprintf("command '%s' failed with exit code %d: %s",
		cmdStr, e.ExitCode, e.Stderr)
}

// TimeoutError 超时错误
type TimeoutError struct {
	Cmd     *Command
	Timeout time.Duration
}

func (e *TimeoutError) Error() string {
	cmdStr := e.Cmd.Name()
	if e.Cmd.Raw() != "" {
		cmdStr = e.Cmd.Raw()
	}
	return fmt.Sprintf("command '%s' timed out after %v",
		cmdStr, e.Timeout)
}
