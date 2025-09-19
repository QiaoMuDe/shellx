// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了Builder结构体和相关构造函数，提供链式调用API来构建命令对象。
package shellx

import (
	"context"
	"io"
	"time"
)

// Builder 命令构建器，提供链式调用
type Builder struct {
	cmd *Command // 命令对象
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
	return &Builder{
		cmd: &Command{
			name: name,
			args: args,
			env:  make(map[string]string),
		},
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
		return &Builder{
			cmd: &Command{
				env: make(map[string]string),
			},
		}
	}

	name := cmdArgs[0]
	args := []string{}
	if len(cmdArgs) > 1 {
		args = cmdArgs[1:]
	}

	return &Builder{
		cmd: &Command{
			name: name,
			args: args,
			env:  make(map[string]string),
		},
	}
}

// NewCmdString 创建新的命令构建器 (字符串方式)
//
// 参数：
//   - cmdStr: 命令字符串
//
// 返回：
//   - *Builder: 命令构建器对象
func NewCmdString(cmdStr string) *Builder {
	return &Builder{
		cmd: &Command{
			raw: cmdStr,
			env: make(map[string]string),
		},
	}
}

// 链式方法
func (b *Builder) WithArgs(args ...string) *Builder {
	b.cmd.args = append(b.cmd.args, args...)
	return b
}

func (b *Builder) WithWorkDir(dir string) *Builder {
	b.cmd.workDir = dir
	return b
}

func (b *Builder) WithEnv(key, value string) *Builder {
	if b.cmd.env == nil {
		b.cmd.env = make(map[string]string)
	}
	b.cmd.env[key] = value
	return b
}

func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	b.cmd.timeout = timeout
	return b
}

func (b *Builder) WithContext(ctx context.Context) *Builder {
	b.cmd.context = ctx
	return b
}

func (b *Builder) WithStdin(stdin io.Reader) *Builder {
	b.cmd.stdin = stdin
	return b
}

func (b *Builder) WithStdout(stdout io.Writer) *Builder {
	b.cmd.stdout = stdout
	return b
}

func (b *Builder) WithStderr(stderr io.Writer) *Builder {
	b.cmd.stderr = stderr
	return b
}

func (b *Builder) WithOptions(opts *ExecuteOptions) *Builder {
	b.cmd.options = opts
	return b
}

func (b *Builder) Build() *Command {
	return b.cmd
}
