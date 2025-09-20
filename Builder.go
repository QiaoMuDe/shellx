// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了Builder结构体和相关构造函数，提供链式调用API来构建命令对象。
package shellx

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// Builder 命令构建器，提供链式调用
type Builder struct {
	shellType ShellType       // shell类型, 默认根据操作系统自动选择(windows使用cmd, 其他系统使用sh)
	raw       string          // 原始命令字符串
	path      string          // 命令路径
	args      []string        // 命令参数
	ctx       context.Context // 上下文
	dir       string          // 工作目录
	env       []string        // 环境变量
	stdin     io.Reader       // 标准输入
	stdout    io.Writer       // 标准输出
	stderr    io.Writer       // 标准错误输出
	timeout   time.Duration   // 超时时间
}

// NewCmd 创建新的命令构建器 (数组方式 - 可变参数)
//
// 参数：
//   - name: 命令名
//   - args: 命令参数列表
//
// 返回：
//   - *Builder: 命令构建器对象
func NewCmd(name string, args ...string) *Builder {
	if name == "" {
		panic("name cannot be empty")
	}

	return &Builder{
		path:      name,         // 命令路径
		args:      args,         // 命令参数
		env:       os.Environ(), // 默认继承父进程的环境变量
		shellType: ShellDefault, // 默认根据操作系统自动选择shell
	}
}

// NewCmds 创建新的命令构建器 (数组方式 - 切片参数，第一个元素为命令名)
//
// 参数：
//   - cmdArgs: 命令参数列表，第一个元素为命令名，后续元素为参数
//
// 返回：
//   - *Builder: 命令构建器对象
func NewCmds(cmdArgs []string) *Builder {
	if len(cmdArgs) == 0 {
		return &Builder{}
	}

	name := cmdArgs[0] // 第一个元素为命令名
	args := []string{} // 后续元素为参数
	if len(cmdArgs) > 1 {
		args = cmdArgs[1:]
	}

	return NewCmd(name, args...)
}

// NewCmdStr 创建新的命令构建器 (字符串方式)
//
// 参数：
//   - cmdStr: 命令字符串
//
// 返回：
//   - *Builder: 命令构建器对象
func NewCmdStr(cmdStr string) *Builder {
	// 使用命令解析器解析命令字符串
	cmds := ParseCmd(cmdStr)
	b := NewCmds(cmds)
	// 保存原始命令字符串
	b.raw = cmdStr
	return b
}

// WithWorkDir 设置命令的工作目录
//
// 参数：
//   - dir: 命令的工作目录
//
// 返回：
//   - *Builder: 命令构建器对象
func (b *Builder) WithWorkDir(dir string) *Builder {
	if dir == "" {
		return b
	}

	b.dir = dir
	return b
}

// WithEnv 设置命令的环境变量
//
// 参数：
//   - key: 环境变量的键
//   - value: 环境变量的值
//
// 返回：
//   - *Builder: 命令构建器对象
func (b *Builder) WithEnv(key, value string) *Builder {
	if b.env == nil {
		b.env = os.Environ()
	}

	if key == "" {
		return b
	}

	// 只有当key不为空时才添加环境变量
	b.env = append(b.env, fmt.Sprintf("%s=%s", key, value))
	return b
}

// WithTimeout 设置命令的超时时间
//
// 参数：
//   - timeout: time.Duration类型，命令执行的超时时间
//
// 返回：
//   - *Builder: 命令构建器对象
func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	if timeout <= 0 {
		return b
	}

	b.timeout = timeout
	return b
}

// WithContext 设置命令的上下文
//
// 参数：
//   - ctx: context.Context类型，用于取消命令执行和超时控制
//
// 返回：
//   - *Builder: 命令构建器对象
func (b *Builder) WithContext(ctx context.Context) *Builder {
	if ctx == nil {
		panic("context cannot be nil")
	}
	b.ctx = ctx
	return b
}

// WithStdin 设置命令的标准输入
//
// 参数：
//   - stdin: io.Reader类型，用于提供命令的标准输入
//
// 返回：
//   - *Builder: 命令构建器对象
func (b *Builder) WithStdin(stdin io.Reader) *Builder {
	if stdin == nil {
		panic("stdin cannot be nil")
	}

	b.stdin = stdin
	return b
}

// WithStdout 设置命令的标准输出
//
// 参数：
//   - stdout: io.Writer类型，用于接收命令的标准输出
//
// 返回：
//   - *Builder: 命令构建器对象
func (b *Builder) WithStdout(stdout io.Writer) *Builder {
	if stdout == nil {
		panic("stdout cannot be nil")
	}

	b.stdout = stdout
	return b
}

// WithStderr 设置命令的标准错误输出
//
// 参数：
//   - stderr: io.Writer类型，用于接收命令的标准错误输出
//
// 返回：
//   - *Builder: 命令构建器对象
func (b *Builder) WithStderr(stderr io.Writer) *Builder {
	if stderr == nil {
		panic("stderr cannot be nil")
	}

	b.stderr = stderr
	return b
}

// WithShell 设置命令的shell类型
//
// 参数：
//   - shell: ShellType类型，表示要使用的shell类型
//
// 返回：
//   - *Builder: 命令构建器对象
func (b *Builder) WithShell(shell ShellType) *Builder {
	b.shellType = shell
	return b
}

// Build 构建并返回命令对象
//
// 返回:
//   - *Command: 构建的命令对象
func (b *Builder) Build() *Command {
	var cmd *exec.Cmd
	if b.ctx != nil {
		// 使用上下文创建命令对象
		if b.shellType != ShellNone {
			// 使用shell执行命令
			cmd = exec.CommandContext(b.ctx, b.shellType.String(), b.shellType.shellFlags(), getCmdStr(b))
		} else {
			// 不使用shell执行命令
			cmd = exec.CommandContext(b.ctx, b.path, b.args...)
		}

	} else {
		// 创建命令对象
		if b.shellType != ShellNone {
			// 使用shell执行命令
			cmd = exec.Command(b.shellType.String(), b.shellType.shellFlags(), getCmdStr(b))
		} else {
			// 不使用shell执行命令
			cmd = exec.Command(b.path, b.args...)
		}
	}

	// 设置工作目录
	if b.dir != "" {
		cmd.Dir = b.dir
	}

	// 设置环境变量
	if b.env != nil {
		cmd.Env = b.env
	}

	// 设置标准输入、输出和错误输出
	if b.stdin != nil {
		cmd.Stdin = b.stdin
	}
	if b.stdout != nil {
		cmd.Stdout = b.stdout
	}
	if b.stderr != nil {
		cmd.Stderr = b.stderr
	}

	// 设置超时时间
	if b.timeout > 0 {
		cmd.WaitDelay = b.timeout
	}

	// 返回包装的命令对象
	return &Command{cmd: cmd}
}
