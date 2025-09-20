package shellx

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

// Command 命令对象
type Command struct {
	cmd  *exec.Cmd // 底层的 exec.Cmd 对象
	once sync.Once // 用于确保每个命令只执行一次
}

// Cmd 获取底层的 exec.Cmd 对象
//
// 返回:
//   - *exec.Cmd: 底层的 exec.Cmd 对象
func (c *Command) Cmd() *exec.Cmd {
	return c.cmd
}

// Exec 执行命令(阻塞)
//
// 返回:
//   - error: 错误信息
func (c *Command) Exec() error {
	return c.cmd.Run()
}

// ExecOutput 执行命令并返回合并后的输出(阻塞)
//
// 返回:
//   - []byte: 命令输出
//   - error: 错误信息
func (c *Command) ExecOutput() ([]byte, error) {
	return c.cmd.CombinedOutput()
}

// ExecStdout 执行命令并返回标准输出(阻塞)
//
// 返回:
//   - []byte: 标准输出
//   - error: 错误信息
func (c *Command) ExecStdout() ([]byte, error) {
	return c.cmd.Output()
}

// ExecResult 执行命令并返回完整的执行结果(阻塞)
//
// 返回:
//   - *Result: 执行结果对象
func (c *Command) ExecResult() *Result {
	// 命令执行开始时间
	startTime := time.Now()

	// 执行命令
	output, err := c.cmd.CombinedOutput()

	// 命令执行结束时间
	endTime := time.Now()

	// 获取命令的退出码
	var exitCode int
	if err != nil {
		exitCode = -1
	}

	// 返回Result对象
	return &Result{
		startTime: startTime,              // 命令开始执行时间
		endTime:   endTime,                // 命令执行结束时间
		duration:  endTime.Sub(startTime), // 命令执行耗时
		err:       err,                    // 命令执行错误信息
		output:    output,                 // 命令输出
		success:   err == nil,             // 命令执行是否成功
		exitCode:  exitCode,               // 命令退出码
	}
}

// ExecAsync 异步执行命令(非阻塞)
//
// 返回:
//   - error: 错误信息
func (c *Command) ExecAsync() error {
	return c.cmd.Start()
}

// Wait 等待命令执行完成(仅在异步执行时有效)
//
// 返回:
//   - error: 错误信息
func (c *Command) Wait() error {
	return c.cmd.Wait()
}

// Kill 杀死当前命令的进程
//
// 返回:
//   - error: 错误信息
func (c *Command) Kill() error {
	if c.cmd.Process == nil {
		return fmt.Errorf("no process to kill")
	}
	return c.cmd.Process.Kill()
}

// Signal 向当前进程发送信号
//
// 参数:
//   - sig: 信号类型
//
// 返回:
//   - error: 错误信息
func (c *Command) Signal(sig os.Signal) error {
	if c.cmd.Process == nil {
		return fmt.Errorf("no process to signal")
	}
	return c.cmd.Process.Signal(sig)
}

// IsRunning 检查进程是否还在运行
//
// 返回:
//   - bool: 是否在运行
func (c *Command) IsRunning() bool {
	if c.cmd.Process == nil {
		return false
	}

	if c.cmd.ProcessState != nil {
		return false // 进程已结束
	}

	// 尝试发送信号0检查进程是否存在
	err := c.cmd.Process.Signal(syscall.Signal(0))
	return err == nil
}

// GetPID 获取进程ID
//
// 返回:
//   - int: 进程ID，如果进程不存在返回0
func (c *Command) GetPID() int {
	if c.cmd.Process == nil {
		return 0
	}
	return c.cmd.Process.Pid
}
