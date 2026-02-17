package shx

import (
	"context"
	"io"
	"os"
	"time"
)

// ParseCmd 将命令字符串解析为命令切片
//
// 使用 mvdan.cc/sh/v3 的 Words 方法解析，支持完整的 shell 语法：
//   - 环境变量：${VAR}, $VAR
//   - 通配符：*.go, test?.txt
//   - 命令替换：$(cmd), `cmd`
//   - 转义字符：\", \', \`
//   - 引号嵌套：支持不同类型引号嵌套
//
// 注意：
//   - 与主包 ParseCmd 行为一致（解析错误时返回空切片）
//   - 不返回错误信息
//
// 参数:
//   - cmdStr: 要解析的命令字符串
//
// 返回值:
//   - []string: 解析后的命令切片，解析错误时返回空切片
func ParseCmd(cmdStr string) []string {
	result, err := parseCmdInternal(cmdStr)
	if err != nil {
		return []string{}
	}
	return result
}

// ParseCmdE 将命令字符串解析为命令切片（带错误信息）
//
// 使用 mvdan.cc/sh/v3 的 Words 方法解析，支持完整的 shell 语法：
//   - 环境变量：${VAR}, $VAR
//   - 通配符：*.go, test?.txt
//   - 命令替换：$(cmd), `cmd`
//   - 转义字符：\", \', \`
//   - 引号嵌套：支持不同类型引号嵌套
//
// 注意：
//   - 返回详细的错误信息，便于调试和错误处理
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

// Run 执行命令
//
// 参数：
//   - cmd: 命令字符串
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	err := shx.Run("echo hello")
func Run(cmd string) error {
	return New(cmd).Exec()
}

// RunToTerminal 执行命令并输出到终端
//
// 参数：
//   - cmd: 命令字符串
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	err := shx.RunToTerminal("echo hello")
func RunToTerminal(cmd string) error {
	return New(cmd).
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Exec()
}

// Out 执行并获取输出
//
// 参数：
//   - cmd: 命令字符串
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	output, err := shx.Out("ls -la")
func Out(cmd string) ([]byte, error) {
	return New(cmd).ExecOutput()
}

// RunWith 超时执行
//
// 参数：
//   - cmd: 命令字符串
//   - timeout: 超时时间
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	err := shx.RunWith("sleep 10", 5*time.Second)
func RunWith(cmd string, timeout time.Duration) error {
	return New(cmd).WithTimeout(timeout).Exec()
}

// OutWith 超时执行并获取输出
//
// 参数：
//   - cmd: 命令字符串
//   - timeout: 超时时间
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	output, err := shx.OutWith("sleep 5", 10*time.Second)
func OutWith(cmd string, timeout time.Duration) ([]byte, error) {
	return New(cmd).WithTimeout(timeout).ExecOutput()
}

// RunWithIO 使用自定义输入输出执行
//
// 参数：
//   - cmd: 命令字符串
//   - stdin: 标准输入
//   - stdout: 标准输出
//   - stderr: 标准错误
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	var buf bytes.Buffer
//	err := shx.RunWithIO("cat", strings.NewReader("hello"), &buf, os.Stderr)
func RunWithIO(cmd string, stdin io.Reader, stdout, stderr io.Writer) error {
	return New(cmd).WithStdin(stdin).WithStdout(stdout).WithStderr(stderr).Exec()
}

// OutWithIO 使用自定义输入输出执行并获取输出
//
// 参数：
//   - cmd: 命令字符串
//   - stdin: 标准输入
//   - stdout: 标准输出
//   - stderr: 标准错误
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	var buf bytes.Buffer
//	output, err := shx.OutWithIO("cat", strings.NewReader("hello"), &buf, os.Stderr)
func OutWithIO(cmd string, stdin io.Reader, stdout, stderr io.Writer) ([]byte, error) {
	return New(cmd).WithStdin(stdin).WithStdout(stdout).WithStderr(stderr).ExecOutput()
}

// RunCtx 使用上下文执行
//
// 参数：
//   - ctx: 上下文
//   - cmd: 命令字符串
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	err := shx.RunCtx(ctx, "sleep 10")
func RunCtx(ctx context.Context, cmd string) error {
	return New(cmd).WithContext(ctx).Exec()
}

// OutCtx 使用上下文执行并获取输出
//
// 参数：
//   - ctx: 上下文
//   - cmd: 命令字符串
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	output, err := shx.OutCtx(ctx, "ls -la")
func OutCtx(ctx context.Context, cmd string) ([]byte, error) {
	return New(cmd).WithContext(ctx).ExecOutput()
}
