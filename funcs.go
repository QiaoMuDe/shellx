package shellx

import (
	"os"
	"os/exec"
	"time"
)

// ParseCmd 将命令字符串解析为命令切片，支持引号处理(单引号、双引号、反引号)
//
// 实现原理：
//  1. 去除首尾空白
//  2. 遍历每个字符
//  3. 处理引号状态切换
//  4. 在非引号状态下遇到空格时分割
//  5. 检测未闭合的引号
//
// 参数:
//   - cmdStr: 要解析的命令字符串
//
// 返回值:
//   - []string: 解析后的命令切片
func ParseCmd(cmdStr string) []string {
	result, _ := parseCmdInternal(cmdStr)
	return result
}

// ParseCmdE 将命令字符串解析为命令切片（带错误信息），支持引号处理(单引号、双引号、反引号)
//
// 实现原理：
//  1. 去除首尾空白
//  2. 遍历每个字符
//  3. 处理引号状态切换
//  4. 在非引号状态下遇到空格时分割
//  5. 检测未闭合的引号并返回错误
//
// 参数:
//   - cmdStr: 要解析的命令字符串
//
// 返回值:
//   - []string: 解析后的命令切片
//   - error: 解析错误，成功时为 nil
func ParseCmdE(cmdStr string) ([]string, error) {
	return parseCmdInternal(cmdStr)
}

// FindCmd 查找命令
//
// 参数:
//   - name: 命令名称
//
// 返回:
//   - string: 命令路径
//   - error: 错误信息
func FindCmd(name string) (string, error) {
	return exec.LookPath(name)
}

// ########################################
// 便捷函数
// ########################################

// ExecStr 执行命令(阻塞)
//
// 参数:
//   - cmdStr: 命令字符串
//
// 返回:
//   - error: 错误信息
func ExecStr(cmdStr string) error {
	return NewCmdStr(cmdStr).WithStdout(os.Stdout).WithStderr(os.Stderr).Exec()
}

// Exec 执行命令(阻塞)
//
// 函数:
//   - name: 命令名
//   - args: 命令参数
//
// 返回:
//   - error: 错误信息
func Exec(name string, args ...string) error {
	return NewCmd(name, args...).WithStdout(os.Stdout).WithStderr(os.Stderr).Exec()
}

// ExecOutStr 执行命令并返回合并后的输出(阻塞)
//
// 参数:
//   - cmdStr: 命令字符串
//
// 返回:
//   - []byte: 输出
//   - error: 错误信息
//
// 注意:
//   - 由于需要捕获默认的stdout和stderr合并输出, 内部已经设置了WithStdout(os.Stdout)和WithStderr(os.Stderr)
func ExecOutStr(cmdStr string) ([]byte, error) {
	return NewCmdStr(cmdStr).ExecOutput()
}

// ExecOut 执行命令并返回合并后的输出(阻塞)
//
// 函数:
//   - name: 命令名
//   - args: 命令参数
//
// 返回:
//   - []byte: 输出
//   - error: 错误信息
//
// 注意:
//   - 由于需要捕获默认的stdout和stderr合并输出, 内部已经设置了WithStdout(os.Stdout)和WithStderr(os.Stderr)
func ExecOut(name string, args ...string) ([]byte, error) {
	return NewCmd(name, args...).ExecOutput()
}

// ExecStrT 执行命令(阻塞，带超时)
//
// 参数:
//   - timeout: 超时时间，如果为0则不设置超时
//   - cmdStr: 命令字符串
//
// 返回:
//   - error: 错误信息
func ExecStrT(timeout time.Duration, cmdStr string) error {
	cmd := NewCmdStr(cmdStr).WithStdout(os.Stdout).WithStderr(os.Stderr)

	// 设置超时
	if timeout > 0 {
		cmd = cmd.WithTimeout(timeout)
	}

	return cmd.Exec()
}

// ExecT 执行命令(阻塞，带超时)
//
// 参数:
//   - timeout: 超时时间，如果为0则不设置超时
//   - name: 命令名
//   - args: 命令参数
//
// 返回:
//   - error: 错误信息
func ExecT(timeout time.Duration, name string, args ...string) error {
	cmd := NewCmd(name, args...).WithStdout(os.Stdout).WithStderr(os.Stderr)

	// 设置超时
	if timeout > 0 {
		cmd = cmd.WithTimeout(timeout)
	}

	return cmd.Exec()
}

// ExecOutStrT 执行命令并返回合并后的输出(阻塞，带超时)
//
// 参数:
//   - timeout: 超时时间，如果为0则不设置超时
//   - cmdStr: 命令字符串
//
// 返回:
//   - []byte: 合并后的输出
//   - error: 错误信息
func ExecOutStrT(timeout time.Duration, cmdStr string) ([]byte, error) {
	cmd := NewCmdStr(cmdStr)

	// 设置超时
	if timeout > 0 {
		cmd = cmd.WithTimeout(timeout)
	}

	return cmd.ExecOutput()
}

// ExecOutT 执行命令并返回合并后的输出(阻塞，带超时)
//
// 参数:
//   - timeout: 超时时间，如果为0则不设置超时
//   - name: 命令名
//   - args: 命令参数
//
// 返回:
//   - []byte: 合并后的输出
//   - error: 错误信息
func ExecOutT(timeout time.Duration, name string, args ...string) ([]byte, error) {
	cmd := NewCmd(name, args...)

	// 设置超时
	if timeout > 0 {
		cmd = cmd.WithTimeout(timeout)
	}

	return cmd.ExecOutput()
}

// ExecCode 执行命令并返回退出码(阻塞)
//
// 参数:
//   - name: 命令名
//   - args: 命令参数
//
// 返回:
//   - int: 退出码
//   - error: 错误信息
func ExecCode(name string, args ...string) (int, error) {
	err := NewCmd(name, args...).Exec()
	if err != nil {
		return -1, err
	}
	return extractExitCode(err), nil
}

// ExecCodeStr 字符串方式执行命令并返回退出码(阻塞)
//
// 参数:
//   - cmdStr: 命令字符串
//
// 返回:
//   - int: 退出码
//   - error: 错误信息
func ExecCodeStr(cmdStr string) (int, error) {
	err := NewCmdStr(cmdStr).Exec()
	if err != nil {
		return -1, err
	}
	return extractExitCode(err), nil
}
