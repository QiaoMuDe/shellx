// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了ExecutionError结构体，表示命令执行失败时的错误信息。
package shellx

import (
	"fmt"
	"time"
)

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
