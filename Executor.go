// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了Executor接口，提供命令执行、异步执行、管道执行等核心功能。
package shellx

// Executor 命令执行器接口
type Executor interface {
	// Exec 执行单个命令
	Exec(cmd *Command) (*Result, error)

	// ExecAsync 异步执行命令
	ExecAsync(cmd *Command) (<-chan *Result, error)

	// ExecPipe 执行命令管道 (可变参数方式)
	ExecPipe(commands ...*Command) (*Result, error)

	// ExecPipes 执行命令管道 (切片方式)
	ExecPipes(commands []*Command) (*Result, error)

	// Kill 终止正在执行的命令
	Kill(pid int) error

	// Running 检查命令是否正在运行
	Running(pid int) bool
}
