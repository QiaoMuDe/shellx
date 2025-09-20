package shellx

import (
	"os/exec"
	"time"
)

// Command 命令对象
type Command struct {
	cmd *exec.Cmd // 底层的 exec.Cmd 对象
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

	// 返回Result对象
	return &Result{
		startTime:    startTime,                     // 命令开始执行时间
		endTime:      endTime,                       // 命令执行结束时间
		duration:     endTime.Sub(startTime),        // 命令执行耗时
		pid:          c.cmd.Process.Pid,             // 命令进程ID
		err:          err,                           // 命令执行错误信息
		output:       output,                        // 命令输出
		success:      err == nil,                    // 命令执行是否成功
		exitCode:     c.cmd.ProcessState.ExitCode(), // 命令退出码
		processState: c.cmd.ProcessState,            // 命令进程状态
	}
}
