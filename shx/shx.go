package shx

import (
	"context"
	"os"
	"strings"
	"time"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

// New 从字符串创建命令
//
// 参数:
//   - cmdStr: 命令字符串
//
// 返回:
//   - *Shx: 命令对象
//
// 示例:
//
//	cmd := shx.New("echo hello world")
//	cmd := shx.New("ls -la | grep .go")
func New(cmdStr string) *Shx {
	return &Shx{
		raw:    cmdStr,
		parser: syntax.NewParser(syntax.Variant(syntax.LangBash)),
		env:    expand.ListEnviron(os.Environ()...),
		dir:    mustGetwd(),
	}
}

// NewArgs 从命令名和可变参数创建命令
//
// 参数:
//   - cmd: 命令名
//   - args: 可变参数列表
//
// 返回:
//   - *Shx: 命令对象
//
// 示例:
//
//	cmd := shx.NewArgs("ls", "-la", "/tmp")
//	cmd := shx.NewArgs("git", "commit", "-m", "message")
func NewArgs(cmd string, args ...string) *Shx {
	// 构建命令字符串
	cmdStr := cmd
	for _, arg := range args {
		cmdStr += " " + arg
	}

	return &Shx{
		raw:    cmdStr,
		parser: syntax.NewParser(syntax.Variant(syntax.LangBash)),
		env:    expand.ListEnviron(os.Environ()...),
		dir:    mustGetwd(),
	}
}

// NewCmds 从命令切片创建命令
//
// 参数:
//   - cmds: 命令切片，每个元素是一个完整的命令部分
//
// 返回:
//   - *Shx: 命令对象
//
// 示例:
//
//	cmd := shx.NewCmds([]string{"ls", "-la", "|", "grep", ".go"})
//	cmd := shx.NewCmds([]string{"echo", "hello", ">", "output.txt"})
func NewCmds(cmds []string) *Shx {
	// 用空格连接命令部分
	cmdStr := strings.Join(cmds, " ")

	return &Shx{
		raw:    cmdStr,
		parser: syntax.NewParser(syntax.Variant(syntax.LangBash)),
		env:    expand.ListEnviron(os.Environ()...),
		dir:    mustGetwd(),
	}
}

// NewScript 从 bash 脚本文件创建命令
//
// 参数:
//   - filePath: 脚本文件路径
//
// 返回:
//   - *Shx: 命令对象
//
// 示例:
//
//	cmd := shx.NewScript("deploy.sh")
//	cmd := shx.NewScript("/path/to/script.sh")
//
// 注意:
//   - 如果 filePath 为空, 会 panic
func NewScript(filePath string) *Shx {
	if filePath == "" {
		panic("script file path must not be empty")
	}

	return &Shx{
		raw:        filePath,
		parser:     syntax.NewParser(syntax.Variant(syntax.LangBash)),
		scriptFile: filePath,
		env:        expand.ListEnviron(os.Environ()...),
		dir:        mustGetwd(),
	}
}

// mustGetwd 获取当前工作目录, 如果失败则返回空字符串
func mustGetwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// markExecuted 标记命令为已执行
// 使用 CompareAndSwap 确保线程安全
func (s *Shx) markExecuted() bool {
	return s.executed.CompareAndSwap(false, true)
}

// Raw 获取原始命令字符串
//
// 返回:
//   - string: 原始命令字符串
func (s *Shx) Raw() string {
	return s.raw
}

// Dir 获取工作目录
//
// 返回:
//   - string: 工作目录
func (s *Shx) Dir() string {
	return s.dir
}

// Env 获取环境变量
//
// 返回:
//   - expand.Environ: 环境变量
func (s *Shx) Env() expand.Environ {
	return s.env
}

// Timeout 获取超时时间
//
// 返回:
//   - time.Duration: 超时时间
func (s *Shx) Timeout() time.Duration {
	return s.timeout
}

// Context 获取上下文
//
// 返回:
//   - context.Context: 上下文
func (s *Shx) Context() context.Context {
	return s.ctx
}

// IsExecuted 检查命令是否已经执行过
//
// 返回:
//   - bool: 是否已执行
func (s *Shx) IsExecuted() bool {
	return s.executed.Load()
}
