// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了Command结构体，表示一个待执行的shell命令，支持数组和字符串两种创建方式。
package shellx

import (
	"context"
	"io"
	"time"
)

// Command 表示一个待执行的shell命令
type Command struct {
	// 基本命令信息 - 支持两种方式
	name string   // 命令名称，如 "ls", "git", "docker" (当使用数组方式时)
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
	context context.Context // 上下文控制

	// 执行选项
	options *ExecuteOptions
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
func (c *Command) Ctx() context.Context   { return c.context }
func (c *Command) Opts() *ExecuteOptions  { return c.options }

// Validate 验证命令参数是否有效
//
// 验证规则:
// 1. 命令必须通过name或raw至少一种方式设置
//
// 返回:
//   - error: 如果验证失败，返回ValidationError；如果验证通过，返回nil
func (c *Command) Validate() error {
	// 检查是否有命令设置
	if c.name == "" && c.raw == "" {
		return &ValidationError{
			Field:   "command",
			Message: "command must be set either by name or raw string",
		}
	}

	return nil
}
