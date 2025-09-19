// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了Executor接口，提供命令执行、异步执行、管道执行等核心功能。
package shellx

// Executor 命令执行器接口
type Executor interface {
	// Exec 执行单个命令
	//
	// 参数：
	//   - cmd: 待执行的命令
	//
	// 返回：
	//   - *Result: 命令执行结果
	//   - error: 错误信息
	Exec(cmd *Command) (*Result, error)

	// ExecAsync 异步执行命令
	//
	// 参数：
	//   - cmd: 待执行的命令
	//
	// 返回：
	//   - <-chan *Result: 命令执行结果通道
	//   - error: 错误信息
	ExecAsync(cmd *Command) (<-chan *Result, error)

	// ExecPipe 执行命令管道 (可变参数方式)
	//
	// 参数：
	//   - commands: 命令管道中的命令列表
	//
	// 返回：
	//   - *Result: 命令执行结果
	//   - error: 错误信息
	//
	// 注意：
	//   - 命令管道中的命令将按顺序执行，前一个命令的输出将作为下一个命令的输入
	ExecPipe(commands ...*Command) (*Result, error)

	// ExecPipes 执行命令管道 (切片方式)
	//
	// 参数：
	//   - commands: 命令管道中的命令列表
	//
	// 返回：
	//   - *Result: 命令执行结果
	//   - error: 错误信息
	//
	// 注意：
	//   - 命令管道中的命令将按顺序执行，前一个命令的输出将作为下一个命令的输入
	ExecPipes(commands []*Command) (*Result, error)

	// Kill 终止正在执行的命令
	//
	// 参数：
	//   - pid: 待终止的命令进程ID
	//
	// 返回：
	//   - error: 错误信息
	Kill(pid int) error

	// Running 检查命令是否正在运行
	//
	// 参数：
	//   - pid: 待检查的命令进程ID
	//
	// 返回：
	//   - bool: 命令是否正在运行
	Running(pid int) bool
}
