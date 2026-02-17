// Package shellx 内部实现模块
// 本文件包含 Command 结构体的内部实现方法，包括：
//   - buildExecCmd: 延迟构建 exec.Cmd 对象，支持上下文和超时控制
//   - cleanup: 资源清理函数，确保上下文取消函数被正确调用
//   - getCmdStr: 命令字符串获取函数，支持原始字符串和参数拼接
//
// 这些方法为 Command 的核心功能提供底层支持。
package shellx

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// buildExecCmd 在执行时构建真正的exec.Cmd对象
//
// 注意:
//   - 该方法会根据上下文和超时时间来创建exec.Cmd对象.
//   - 如果上下文设置了超时时间, 则会忽略timeout参数.
//   - 此方法不是并发安全的，不要在多个goroutine中并发调用
//
// 返回:
//   - error: 构建错误，成功时为 nil
func (c *Command) buildExecCmd() error {
	if c.execCmd != nil {
		return nil // 已经构建过了
	}

	// 统一验证所有参数
	if err := c.validateAllParameters(); err != nil {
		return err // 返回错误，让上层处理
	}

	// 根据实际情况选择创建方式，避免不必要的上下文使用
	if c.userCtx != nil {
		// 用户设置了上下文，使用CommandContext(忽略timeout)
		if c.shellType != ShellNone {
			cmdStr := c.getCmdStr()
			c.execCmd = exec.CommandContext(c.userCtx, c.shellType.String(), c.shellType.shellFlags(), cmdStr)
		} else {
			c.execCmd = exec.CommandContext(c.userCtx, c.name, c.args...)
		}

	} else if c.timeout > 0 {
		// 只设置了超时，创建超时上下文
		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		c.cancel = cancel // 保存cancel函数用于资源清理
		c.userCtx = ctx   // 将内部创建的上下文保存到userCtx，方便错误判断

		if c.shellType != ShellNone {
			cmdStr := c.getCmdStr()
			c.execCmd = exec.CommandContext(ctx, c.shellType.String(), c.shellType.shellFlags(), cmdStr)
		} else {
			c.execCmd = exec.CommandContext(ctx, c.name, c.args...)
		}

	} else {
		// 都没有设置，使用普通的Command(不带上下文)
		if c.shellType != ShellNone {
			cmdStr := c.getCmdStr()
			c.execCmd = exec.Command(c.shellType.String(), c.shellType.shellFlags(), cmdStr)
		} else {
			c.execCmd = exec.Command(c.name, c.args...)
		}
	}

	// 设置exec.Cmd的其他属性
	c.execCmd.Dir = c.dir       // 设置工作目录
	c.execCmd.Env = c.envs      // 设置环境变量
	c.execCmd.Stdin = c.stdin   // 设置标准输入
	c.execCmd.Stdout = c.stdout // 设置标准输出
	c.execCmd.Stderr = c.stderr // 设置标准错误输出

	return nil
}

// cleanup 清理资源
func (c *Command) cleanup() {
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
}

// getCmdStr 获取命令字符串
//
// 参数:
//   - c: 命令对象
//
// 返回:
//   - string: 命令字符串
func (c *Command) getCmdStr() string {
	if c == nil {
		return ""
	}

	// 构建命令字符串
	if c.raw != "" {
		return c.raw
	} else if len(c.args) == 0 {
		return c.name
	} else {
		return fmt.Sprintf("%s %s", c.name, strings.Join(c.args, " "))
	}
}

// extractExitCode 从错误中提取退出码
//
// 参数:
//   - err: 错误对象
//
// 返回:
//   - int: 退出码(0表示成功，-1表示无法提取的执行错误，其他值表示命令返回的退出码)
func extractExitCode(err error) int {
	if err == nil {
		return 0
	}

	// 尝试从ExitError中提取真实的退出码
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}

	// 其他类型的错误（如命令不存在、超时等）返回-1
	return -1
}

// validateEnvVar 验证环境变量格式
//
// 参数:
//   - env: 环境变量字符串，格式为 "key=value"
//
// 返回:
//   - error: 错误信息
func validateEnvVar(env string) error {
	if strings.TrimSpace(env) == "" {
		return fmt.Errorf("environment variable cannot be empty")
	}

	parts := strings.SplitN(env, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid environment variable format, expected 'key=value': %s", env)
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return fmt.Errorf("environment variable key cannot be empty: %s", env)
	}

	return nil
}

// validateAllParameters 统一验证命令对象的所有参数
//
// 该方法在执行命令前调用，验证所有配置参数的有效性。
// 相比于在各个配置方法中直接panic，这里返回错误信息，
// 让调用者决定如何处理错误，避免程序意外退出。
//
// 返回:
//   - error: 验证错误，验证通过时为 nil
func (c *Command) validateAllParameters() error {
	// 验证基本命令参数
	if c.name == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	// 验证工作目录
	if c.dir != "" {
		// 检查目录是否存在
		info, err := os.Lstat(c.dir)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("working directory does not exist: %s", c.dir)
			}
			return fmt.Errorf("failed to check working directory: %s, error: %v", c.dir, err)
		}

		// 检查是否为目录
		if !info.IsDir() {
			return fmt.Errorf("specified path is not a directory: %s", c.dir)
		}
	}

	// 验证环境变量格式
	if len(c.envs) > 0 {
		for _, env := range c.envs {
			if err := validateEnvVar(env); err != nil {
				return fmt.Errorf("environment variable format error: %w", err)
			}
		}
	}

	return nil
}
