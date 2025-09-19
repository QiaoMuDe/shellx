// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了Command结构体, 表示一个待执行的shell命令, 支持数组和字符串两种创建方式。
package shellx

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// Command 表示一个待执行的shell命令
type Command struct {
	// 基本命令信息 - 支持两种方式
	name string   // 命令名称, 如 "ls", "git", "docker" (当使用数组方式时)
	args []string // 命令参数列表 (当使用数组方式时)
	raw  string   // 原始命令字符串 (当使用字符串方式时)

	// 执行环境配置
	workDir string            // 工作目录
	env     map[string]string // 环境变量

	// 输入输出配置
	stdin  io.Reader // 标准输入
	stdout io.Writer // 标准输出重定向
	stderr io.Writer // 标准错误重定向

	// 执行控制
	timeout time.Duration   // 超时时间
	ctx     context.Context // 上下文控制

	// 执行选项
	shell ShellType // 指定shell类型
}

// 提供公共访问方法
func (c *Command) Name() string           { return c.name }
func (c *Command) Args() []string         { return c.args }
func (c *Command) Raw() string            { return c.raw }
func (c *Command) Dir() string            { return c.workDir }
func (c *Command) Env() map[string]string { return c.env }
func (c *Command) Input() io.Reader       { return c.stdin }
func (c *Command) Output() io.Writer      { return c.stdout }
func (c *Command) ErrOutput() io.Writer   { return c.stderr }
func (c *Command) Timeout() time.Duration { return c.timeout }
func (c *Command) Ctx() context.Context   { return c.ctx }
func (c *Command) Shell() ShellType       { return c.shell }

// validate 验证命令参数是否有效
//
// 验证规则:
// 1. 命令必须通过name或raw至少一种方式设置
// 2. 如果设置了timeout, 必须为正值
// 3. 如果设置了context, 不能为nil
// 4. 环境变量的键不能为空
//
// 返回:
//   - error: 如果验证失败, 返回ValidationError；如果验证通过, 返回nil
func (c *Command) validate() error {
	// 检查是否有命令设置
	if c.name == "" && c.raw == "" {
		return &ValidationError{
			Field:   "command",
			Message: "command must be set either by name or raw string",
		}
	}

	// 验证超时设置
	if c.timeout < 0 {
		return &ValidationError{
			Field:   "timeout",
			Message: "timeout must be a non-negative value",
		}
	}

	// 验证上下文
	if c.ctx != nil && c.ctx.Err() != nil {
		return &ValidationError{
			Field:   "ctx",
			Message: "context is already canceled or timeout: " + c.ctx.Err().Error(),
		}
	}

	// 验证环境变量
	for k := range c.env {
		if k == "" {
			return &ValidationError{
				Field:   "env",
				Message: "environment variable key cannot be empty",
			}
		}
	}

	return nil
}

// Exec 执行命令, 并返回结果
//
// 返回:
//   - *Result: 执行结果
//   - error: 如果执行失败, 返回错误信息
func (c *Command) Exec() (*Result, error) {
	// 验证命令参数
	if err := c.validate(); err != nil {
		return nil, err
	}

	return &Result{}, nil
}

// ExecOut 执行命令, 并返回合并后的输出和错误信息
//
// 返回:
//   - []byte: 合并后的输出和错误信息
//   - error: 如果执行失败, 返回错误信息
func (c *Command) ExecOutput() ([]byte, error) {
	// 验证命令参数
	if err := c.validate(); err != nil {
		return nil, err
	}

	return []byte{}, nil
}

// createCmd 用于创建内部使用的exec.Cmd对象
//
// 返回:
//   - *exec.Cmd: 内部使用的exec.Cmd对象
func (c *Command) createCmd() *exec.Cmd {
	// 根据上下文创建命令
	var cmd *exec.Cmd
	if c.ctx != nil {
		cmd = exec.CommandContext(c.ctx, c.name, c.args...)
	} else {
		cmd = exec.Command(c.name, c.args...)
	}

	// 设置工作目录
	if c.workDir != "" {
		cmd.Dir = c.workDir
	}

	// 设置环境变量
	if len(c.env) > 0 {
		envs := make([]string, 0, len(c.env))
		for k, v := range c.env {
			envs = append(envs, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = append(os.Environ(), envs...)
	}

	// 设置输入
	if c.stdin != nil {
		cmd.Stdin = c.stdin
	}

	// 设置输出
	if c.stdout != nil {
		cmd.Stdout = c.stdout
	}

	// 设置错误输出
	if c.stderr != nil {
		cmd.Stderr = c.stderr
	}

	// 设置超时
	if c.timeout > 0 {
		cmd.WaitDelay = c.timeout
	}

	// 返回创建的命令对象
	return cmd
}
