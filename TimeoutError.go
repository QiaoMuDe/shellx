// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了TimeoutError结构体，表示命令执行超时时的错误信息。
package shellx

import (
	"fmt"
	"time"
)

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
