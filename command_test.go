package shellx

import (
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestCommandCmd 测试 Cmd 方法
func TestCommandCmd(t *testing.T) {
	builder := NewCmd("echo", "test")
	cmd := builder.Build()

	if cmd.Cmd() == nil {
		t.Errorf("Cmd() 应该返回非nil的exec.Cmd对象")
	}

	if cmd.Cmd().Path == "" && cmd.Cmd().Args == nil {
		t.Errorf("Cmd() 返回的exec.Cmd对象应该包含命令信息")
	}
}

// TestCommandIsExecuted 测试 IsExecuted 方法
func TestCommandIsExecuted(t *testing.T) {
	builder := NewCmd("echo", "test")
	cmd := builder.Build()

	// 初始状态应该是未执行
	if cmd.IsExecuted() {
		t.Errorf("新创建的Command应该是未执行状态")
	}

	// 执行后应该是已执行状态
	_ = cmd.Exec()
	if !cmd.IsExecuted() {
		t.Errorf("执行后的Command应该是已执行状态")
	}
}

// TestCommandGetPID 测试 GetPID 方法
func TestCommandGetPID(t *testing.T) {
	builder := NewCmd("echo", "test")
	cmd := builder.Build()

	// 未执行时PID应该为0
	if cmd.GetPID() != 0 {
		t.Errorf("未执行的Command的PID应该为0，实际为%d", cmd.GetPID())
	}
}

// TestCommandExecOnce 测试命令只能执行一次的限制
func TestCommandExecOnce(t *testing.T) {
	tests := []struct {
		name     string
		execFunc func(*Command) error
	}{
		{
			name: "Exec方法",
			execFunc: func(c *Command) error {
				return c.Exec()
			},
		},
		{
			name: "ExecOutput方法",
			execFunc: func(c *Command) error {
				_, err := c.ExecOutput()
				return err
			},
		},
		{
			name: "ExecStdout方法",
			execFunc: func(c *Command) error {
				_, err := c.ExecStdout()
				return err
			},
		},
		{
			name: "ExecAsync方法",
			execFunc: func(c *Command) error {
				return c.ExecAsync()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewCmd("echo", "test")
			cmd := builder.Build()

			// 第一次执行应该成功（或至少不是因为重复执行而失败）
			err1 := tt.execFunc(cmd)

			// 第二次执行应该返回"已执行"错误
			err2 := tt.execFunc(cmd)
			if err2 == nil {
				t.Errorf("%s 第二次执行应该返回错误", tt.name)
			}

			expectedMsg := "command has already been executed"
			if !strings.Contains(err2.Error(), expectedMsg) {
				t.Errorf("%s 第二次执行错误信息应该包含'%s'，实际为'%s'", tt.name, expectedMsg, err2.Error())
			}

			// 验证命令状态
			if !cmd.IsExecuted() {
				t.Errorf("%s 执行后命令应该标记为已执行", tt.name)
			}

			// 对于异步执行，需要等待完成
			if tt.name == "ExecAsync方法" && err1 == nil {
				_ = cmd.Wait()
			}
		})
	}
}

// TestCommandExecResult 测试 ExecResult 方法
func TestCommandExecResult(t *testing.T) {
	builder := NewCmd("echo", "test")
	cmd := builder.Build()

	result := cmd.ExecResult()

	// 验证Result对象的基本属性
	if result == nil {
		t.Fatalf("ExecResult() 应该返回非nil的Result对象")
	}

	// 验证时间相关字段
	if result.Start().IsZero() {
		t.Errorf("Result.Start() 不应该是零时间")
	}

	if result.End().IsZero() {
		t.Errorf("Result.End() 不应该是零时间")
	}

	if result.Duration() <= 0 {
		t.Errorf("Result.Duration() 应该大于0，实际为%v", result.Duration())
	}

	if result.End().Before(result.Start()) {
		t.Errorf("结束时间应该在开始时间之后")
	}

	// 验证命令已执行
	if !cmd.IsExecuted() {
		t.Errorf("ExecResult() 执行后命令应该标记为已执行")
	}

	// 第二次调用应该返回错误
	result2 := cmd.ExecResult()
	if result2.Error() == nil {
		t.Errorf("第二次调用ExecResult()应该返回错误")
	}

	if result2.Success() {
		t.Errorf("第二次调用ExecResult()应该标记为失败")
	}

	if result2.Code() != -1 {
		t.Errorf("第二次调用ExecResult()的退出码应该为-1，实际为%d", result2.Code())
	}
}

// TestCommandKillAndSignal 测试 Kill 和 Signal 方法
func TestCommandKillAndSignal(t *testing.T) {
	builder := NewCmd("echo", "test")
	cmd := builder.Build()

	// 未执行时应该返回错误
	err := cmd.Kill()
	if err == nil {
		t.Errorf("未执行的命令调用Kill()应该返回错误")
	}

	expectedMsg := "no process to kill"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Kill()错误信息应该包含'%s'，实际为'%s'", expectedMsg, err.Error())
	}

	// 测试Signal方法
	err = cmd.Signal(syscall.SIGTERM)
	if err == nil {
		t.Errorf("未执行的命令调用Signal()应该返回错误")
	}

	expectedMsg = "no process to signal"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Signal()错误信息应该包含'%s'，实际为'%s'", expectedMsg, err.Error())
	}
}

// TestCommandIsRunning 测试 IsRunning 方法
func TestCommandIsRunning(t *testing.T) {
	builder := NewCmd("echo", "test")
	cmd := builder.Build()

	// 未执行时应该返回false
	if cmd.IsRunning() {
		t.Errorf("未执行的命令IsRunning()应该返回false")
	}
}

// TestCommandWait 测试 Wait 方法
func TestCommandWait(t *testing.T) {
	builder := NewCmd("echo", "test")
	cmd := builder.Build()

	// 未启动的命令调用Wait应该返回错误
	err := cmd.Wait()
	if err == nil {
		t.Errorf("未启动的命令调用Wait()应该返回错误")
	}
}

// TestCommandWithDifferentShells 测试不同Shell类型的命令
func TestCommandWithDifferentShells(t *testing.T) {
	shells := []ShellType{
		ShellNone,
		ShellDefault,
		ShellSh,
		ShellBash,
		ShellCmd,
		ShellPwsh,
		ShellPowerShell,
	}

	for _, shell := range shells {
		t.Run(shell.String(), func(t *testing.T) {
			builder := NewCmd("echo", "test").WithShell(shell)
			cmd := builder.Build()

			if cmd == nil {
				t.Errorf("使用%s shell构建的命令不应该为nil", shell.String())
			}

			if cmd.Cmd() == nil {
				t.Errorf("使用%s shell构建的命令的Cmd()不应该为nil", shell.String())
			}

			// 验证命令未执行
			if cmd.IsExecuted() {
				t.Errorf("新构建的命令不应该标记为已执行")
			}
		})
	}
}

// TestCommandWithTimeout 测试带超时的命令
func TestCommandWithTimeout(t *testing.T) {
	timeout := 100 * time.Millisecond
	builder := NewCmd("echo", "test").WithTimeout(timeout)
	cmd := builder.Build()

	if cmd == nil {
		t.Fatalf("带超时的命令构建失败")
	}

	// 验证超时设置（通过检查底层exec.Cmd的WaitDelay字段）
	if cmd.Cmd().WaitDelay != timeout {
		t.Errorf("命令超时设置不正确，期望%v，实际%v", timeout, cmd.Cmd().WaitDelay)
	}
}

// TestCommandWithWorkDir 测试带工作目录的命令
func TestCommandWithWorkDir(t *testing.T) {
	workDir := "/tmp"
	if runtime.GOOS == "windows" {
		workDir = "C:\\temp"
	}

	builder := NewCmd("echo", "test").WithWorkDir(workDir)
	cmd := builder.Build()

	if cmd == nil {
		t.Fatalf("带工作目录的命令构建失败")
	}

	// 验证工作目录设置
	if cmd.Cmd().Dir != workDir {
		t.Errorf("命令工作目录设置不正确，期望%s，实际%s", workDir, cmd.Cmd().Dir)
	}
}

// TestCommandWithEnv 测试带环境变量的命令
func TestCommandWithEnv(t *testing.T) {
	builder := NewCmd("echo", "test").WithEnv("TEST_VAR", "test_value")
	cmd := builder.Build()

	if cmd == nil {
		t.Fatalf("带环境变量的命令构建失败")
	}

	// 验证环境变量设置
	env := cmd.Cmd().Env
	if env == nil {
		t.Errorf("命令环境变量不应该为nil")
		return
	}

	found := false
	for _, e := range env {
		if e == "TEST_VAR=test_value" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("命令环境变量中没有找到TEST_VAR=test_value")
	}
}

// TestCommandConcurrency 测试命令的并发安全性
func TestCommandConcurrency(t *testing.T) {
	builder := NewCmd("echo", "test")
	cmd := builder.Build()

	// 并发调用IsExecuted方法
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			_ = cmd.IsExecuted()
			_ = cmd.GetPID()
			_ = cmd.IsRunning()
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证命令状态仍然正确
	if cmd.IsExecuted() {
		t.Errorf("并发读取后命令不应该标记为已执行")
	}
}

// TestCommandEdgeCases 测试边界情况
func TestCommandEdgeCases(t *testing.T) {
	// 测试空命令名
	builder := NewCmds([]string{})
	cmd := builder.Build()

	if cmd == nil {
		t.Errorf("空命令构建不应该返回nil")
	}

	// 测试命令执行状态的原子性
	builder2 := NewCmd("echo", "test")
	cmd2 := builder2.Build()

	// 验证初始状态
	if cmd2.IsExecuted() {
		t.Errorf("新命令不应该标记为已执行")
	}

	// 验证PID初始值
	if cmd2.GetPID() != 0 {
		t.Errorf("新命令的PID应该为0")
	}

	// 验证IsRunning初始值
	if cmd2.IsRunning() {
		t.Errorf("新命令不应该处于运行状态")
	}
}
