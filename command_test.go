// Package shellx 命令执行测试模块
// 本文件包含 Command 结构体及其相关方法的单元测试，包括：
//   - 基本命令执行测试
//   - 超时和取消机制测试
//   - 错误处理测试
//   - 配置方法测试
//
// 确保 Command 的各项功能正常工作，提供回归测试保障。
package shellx

import (
	"bytes"
	"context"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNewCmd 测试NewCmd函数
func TestNewCmd(t *testing.T) {
	t.Run("正常创建", func(t *testing.T) {
		cmd := NewCmd("echo", "hello", "world")
		if cmd.Name() != "echo" {
			t.Errorf("期望命令名为 'echo', 实际为 '%s'", cmd.Name())
		}
		args := cmd.Args()
		if len(args) != 2 || args[0] != "hello" || args[1] != "world" {
			t.Errorf("期望参数为 ['hello', 'world'], 实际为 %v", args)
		}
	})

	t.Run("空命令名panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("期望panic，但没有发生")
			}
		}()
		NewCmd("")
	})

	t.Run("无参数", func(t *testing.T) {
		cmd := NewCmd("ls")
		if cmd.Name() != "ls" {
			t.Errorf("期望命令名为 'ls', 实际为 '%s'", cmd.Name())
		}
		args := cmd.Args()
		if len(args) != 0 {
			t.Errorf("期望无参数, 实际为 %v", args)
		}
	})
}

// TestNewCmds 测试NewCmds函数
func TestNewCmds(t *testing.T) {
	t.Run("正常创建", func(t *testing.T) {
		cmd := NewCmds([]string{"git", "status", "--porcelain"})
		if cmd.Name() != "git" {
			t.Errorf("期望命令名为 'git', 实际为 '%s'", cmd.Name())
		}
		args := cmd.Args()
		if len(args) != 2 || args[0] != "status" || args[1] != "--porcelain" {
			t.Errorf("期望参数为 ['status', '--porcelain'], 实际为 %v", args)
		}
	})

	t.Run("空切片panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("期望panic，但没有发生")
			}
		}()
		NewCmds([]string{})
	})

	t.Run("单个元素", func(t *testing.T) {
		cmd := NewCmds([]string{"pwd"})
		if cmd.Name() != "pwd" {
			t.Errorf("期望命令名为 'pwd', 实际为 '%s'", cmd.Name())
		}
		args := cmd.Args()
		if len(args) != 0 {
			t.Errorf("期望无参数, 实际为 %v", args)
		}
	})
}

// TestNewCmdStr 测试NewCmdStr函数
func TestNewCmdStr(t *testing.T) {
	t.Run("简单命令", func(t *testing.T) {
		cmd := NewCmdStr("echo hello world")
		if cmd.Name() != "echo" {
			t.Errorf("期望命令名为 'echo', 实际为 '%s'", cmd.Name())
		}
		args := cmd.Args()
		if len(args) != 2 || args[0] != "hello" || args[1] != "world" {
			t.Errorf("期望参数为 ['hello', 'world'], 实际为 %v", args)
		}
		if cmd.Raw() != "echo hello world" {
			t.Errorf("期望原始命令为 'echo hello world', 实际为 '%s'", cmd.Raw())
		}
	})

	t.Run("带引号的命令", func(t *testing.T) {
		cmd := NewCmdStr(`git commit -m "test message"`)
		if cmd.Name() != "git" {
			t.Errorf("期望命令名为 'git', 实际为 '%s'", cmd.Name())
		}
		args := cmd.Args()
		expected := []string{"commit", "-m", "test message"}
		if len(args) != len(expected) {
			t.Errorf("期望参数长度为 %d, 实际为 %d", len(expected), len(args))
		}
		for i, arg := range args {
			if arg != expected[i] {
				t.Errorf("期望参数[%d]为 '%s', 实际为 '%s'", i, expected[i], arg)
			}
		}
	})
}

// TestWithWorkDir 测试WithWorkDir方法
func TestWithWorkDir(t *testing.T) {
	cmd := NewCmd("pwd")

	t.Run("设置工作目录", func(t *testing.T) {
		cmd.WithWorkDir("/tmp")
		if cmd.WorkDir() != "/tmp" {
			t.Errorf("期望工作目录为 '/tmp', 实际为 '%s'", cmd.WorkDir())
		}
	})

	t.Run("空目录被忽略", func(t *testing.T) {
		originalDir := cmd.WorkDir()
		cmd.WithWorkDir("")
		if cmd.WorkDir() != originalDir {
			t.Errorf("空目录应该被忽略，工作目录不应该改变")
		}
	})

	t.Run("链式调用", func(t *testing.T) {
		result := cmd.WithWorkDir("/home")
		if result != cmd {
			t.Error("WithWorkDir应该返回自身以支持链式调用")
		}
	})
}

// TestWithEnv 测试WithEnv方法
func TestWithEnv(t *testing.T) {
	cmd := NewCmd("env")

	t.Run("设置环境变量", func(t *testing.T) {
		cmd.WithEnv("TEST_VAR", "test_value")
		envs := cmd.Env()
		found := false
		for _, env := range envs {
			if env == "TEST_VAR=test_value" {
				found = true
				break
			}
		}
		if !found {
			t.Error("环境变量 TEST_VAR=test_value 未找到")
		}
	})

	t.Run("空key被忽略", func(t *testing.T) {
		originalEnvCount := len(cmd.Env())
		cmd.WithEnv("", "value")
		if len(cmd.Env()) != originalEnvCount {
			t.Error("空key的环境变量应该被忽略")
		}
	})

	t.Run("链式调用", func(t *testing.T) {
		result := cmd.WithEnv("CHAIN_TEST", "value")
		if result != cmd {
			t.Error("WithEnv应该返回自身以支持链式调用")
		}
	})
}

// TestWithEnvs 测试WithEnvs方法
func TestWithEnvs(t *testing.T) {
	cmd := NewCmd("env")

	t.Run("批量设置环境变量", func(t *testing.T) {
		envs := []string{"VAR1=value1", "VAR2=value2", "VAR3=value3"}
		cmd.WithEnvs(envs)

		cmdEnvs := cmd.Env()
		for _, expectedEnv := range envs {
			found := false
			for _, env := range cmdEnvs {
				if env == expectedEnv {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("环境变量 %s 未找到", expectedEnv)
			}
		}
	})

	t.Run("空切片被忽略", func(t *testing.T) {
		originalEnvCount := len(cmd.Env())
		cmd.WithEnvs([]string{})
		if len(cmd.Env()) != originalEnvCount {
			t.Error("空环境变量切片应该被忽略")
		}
	})

	t.Run("链式调用", func(t *testing.T) {
		result := cmd.WithEnvs([]string{"BATCH_TEST=value"})
		if result != cmd {
			t.Error("WithEnvs应该返回自身以支持链式调用")
		}
	})
}

// TestWithTimeout 测试WithTimeout方法
func TestWithTimeout(t *testing.T) {
	cmd := NewCmd("sleep", "1")

	t.Run("设置超时时间", func(t *testing.T) {
		timeout := 5 * time.Second
		cmd.WithTimeout(timeout)
		if cmd.Timeout() != timeout {
			t.Errorf("期望超时时间为 %v, 实际为 %v", timeout, cmd.Timeout())
		}
	})

	t.Run("零或负数超时被忽略", func(t *testing.T) {
		originalTimeout := cmd.Timeout()
		cmd.WithTimeout(0)
		if cmd.Timeout() != originalTimeout {
			t.Error("零超时应该被忽略")
		}

		cmd.WithTimeout(-1 * time.Second)
		if cmd.Timeout() != originalTimeout {
			t.Error("负数超时应该被忽略")
		}
	})

	t.Run("链式调用", func(t *testing.T) {
		result := cmd.WithTimeout(1 * time.Second)
		if result != cmd {
			t.Error("WithTimeout应该返回自身以支持链式调用")
		}
	})
}

// TestWithContext 测试WithContext方法
func TestWithContext(t *testing.T) {
	cmd := NewCmd("sleep", "1")

	t.Run("设置上下文", func(t *testing.T) {
		ctx := context.Background()
		cmd.WithContext(ctx)
		// 由于userCtx是私有字段，我们通过其他方式验证
		// 这里主要测试不会panic
	})

	t.Run("nil上下文panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("期望panic，但没有发生")
			}
		}()
		//nolint:all
		cmd.WithContext(nil)
	})

	t.Run("链式调用", func(t *testing.T) {
		ctx := context.Background()
		result := cmd.WithContext(ctx)
		if result != cmd {
			t.Error("WithContext应该返回自身以支持链式调用")
		}
	})
}

// TestWithStdin 测试WithStdin方法
func TestWithStdin(t *testing.T) {
	cmd := NewCmd("cat")

	t.Run("设置标准输入", func(t *testing.T) {
		stdin := strings.NewReader("test input")
		cmd.WithStdin(stdin)
		// 由于stdin是私有字段，我们通过其他方式验证
		// 这里主要测试不会panic
	})

	t.Run("nil输入panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("期望panic，但没有发生")
			}
		}()
		cmd.WithStdin(nil)
	})

	t.Run("链式调用", func(t *testing.T) {
		stdin := strings.NewReader("test")
		result := cmd.WithStdin(stdin)
		if result != cmd {
			t.Error("WithStdin应该返回自身以支持链式调用")
		}
	})
}

// TestWithStdout 测试WithStdout方法
func TestWithStdout(t *testing.T) {
	cmd := NewCmd("echo", "test")

	t.Run("设置标准输出", func(t *testing.T) {
		var stdout bytes.Buffer
		cmd.WithStdout(&stdout)
		// 由于stdout是私有字段，我们通过其他方式验证
		// 这里主要测试不会panic
	})

	t.Run("nil输出panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("期望panic，但没有发生")
			}
		}()
		cmd.WithStdout(nil)
	})

	t.Run("链式调用", func(t *testing.T) {
		var stdout bytes.Buffer
		result := cmd.WithStdout(&stdout)
		if result != cmd {
			t.Error("WithStdout应该返回自身以支持链式调用")
		}
	})
}

// TestWithStderr 测试WithStderr方法
func TestWithStderr(t *testing.T) {
	cmd := NewCmd("echo", "test")

	t.Run("设置标准错误输出", func(t *testing.T) {
		var stderr bytes.Buffer
		cmd.WithStderr(&stderr)
		// 由于stderr是私有字段，我们通过其他方式验证
		// 这里主要测试不会panic
	})

	t.Run("nil错误输出panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("期望panic，但没有发生")
			}
		}()
		cmd.WithStderr(nil)
	})

	t.Run("链式调用", func(t *testing.T) {
		var stderr bytes.Buffer
		result := cmd.WithStderr(&stderr)
		if result != cmd {
			t.Error("WithStderr应该返回自身以支持链式调用")
		}
	})
}

// TestWithShell 测试WithShell方法
func TestWithShell(t *testing.T) {
	cmd := NewCmd("echo", "test")

	t.Run("设置shell类型", func(t *testing.T) {
		cmd.WithShell(ShellBash)
		if cmd.ShellType() != ShellBash {
			t.Errorf("期望shell类型为 %v, 实际为 %v", ShellBash, cmd.ShellType())
		}
	})

	t.Run("链式调用", func(t *testing.T) {
		result := cmd.WithShell(ShellCmd)
		if result != cmd {
			t.Error("WithShell应该返回自身以支持链式调用")
		}
	})
}

// TestIsExecuted 测试IsExecuted方法
func TestIsExecuted(t *testing.T) {
	cmd := NewCmd("echo", "test")

	t.Run("初始状态未执行", func(t *testing.T) {
		if cmd.IsExecuted() {
			t.Error("新创建的命令应该是未执行状态")
		}
	})

	t.Run("执行后状态改变", func(t *testing.T) {
		// 使用一个简单的命令来测试
		testCmd := NewCmd("echo", "hello")
		err := testCmd.Exec()
		if err != nil {
			t.Fatalf("命令执行失败: %v", err)
		}

		if !testCmd.IsExecuted() {
			t.Error("执行后的命令应该是已执行状态")
		}
	})
}

// TestExecOnce 测试命令只能执行一次
func TestExecOnce(t *testing.T) {
	cmd := NewCmd("echo", "test")

	// 第一次执行应该成功
	err := cmd.Exec()
	if err != nil {
		t.Fatalf("第一次执行失败: %v", err)
	}

	// 第二次执行应该失败
	err = cmd.Exec()
	if err == nil {
		t.Error("第二次执行应该失败")
	}

	if !strings.Contains(err.Error(), "already been executed") {
		t.Errorf("错误信息应该包含 'already been executed', 实际为: %v", err)
	}
}

// TestConcurrentAccess 测试并发访问安全性
func TestConcurrentAccess(t *testing.T) {
	cmd := NewCmd("echo", "test")

	var wg sync.WaitGroup
	const goroutines = 10

	// 并发设置配置
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			cmd.WithEnv("TEST_VAR", "value").
				WithWorkDir("/tmp").
				WithTimeout(5 * time.Second)
		}(i)
	}

	wg.Wait()

	// 验证最终状态
	if cmd.WorkDir() != "/tmp" {
		t.Errorf("期望工作目录为 '/tmp', 实际为 '%s'", cmd.WorkDir())
	}

	if cmd.Timeout() != 5*time.Second {
		t.Errorf("期望超时时间为 5s, 实际为 %v", cmd.Timeout())
	}
}

// TestGetters 测试各种getter方法
func TestGetters(t *testing.T) {
	cmd := NewCmd("git", "status", "--porcelain").
		WithWorkDir("/home").
		WithEnv("GIT_CONFIG", "test").
		WithTimeout(10 * time.Second).
		WithShell(ShellBash)

	t.Run("Name getter", func(t *testing.T) {
		if cmd.Name() != "git" {
			t.Errorf("期望命令名为 'git', 实际为 '%s'", cmd.Name())
		}
	})

	t.Run("Args getter", func(t *testing.T) {
		args := cmd.Args()
		expected := []string{"status", "--porcelain"}
		if len(args) != len(expected) {
			t.Errorf("期望参数长度为 %d, 实际为 %d", len(expected), len(args))
		}
		for i, arg := range args {
			if arg != expected[i] {
				t.Errorf("期望参数[%d]为 '%s', 实际为 '%s'", i, expected[i], arg)
			}
		}
	})

	t.Run("WorkDir getter", func(t *testing.T) {
		if cmd.WorkDir() != "/home" {
			t.Errorf("期望工作目录为 '/home', 实际为 '%s'", cmd.WorkDir())
		}
	})

	t.Run("Env getter", func(t *testing.T) {
		envs := cmd.Env()
		found := false
		for _, env := range envs {
			if env == "GIT_CONFIG=test" {
				found = true
				break
			}
		}
		if !found {
			t.Error("环境变量 GIT_CONFIG=test 未找到")
		}
	})

	t.Run("Timeout getter", func(t *testing.T) {
		if cmd.Timeout() != 10*time.Second {
			t.Errorf("期望超时时间为 10s, 实际为 %v", cmd.Timeout())
		}
	})

	t.Run("ShellType getter", func(t *testing.T) {
		if cmd.ShellType() != ShellBash {
			t.Errorf("期望shell类型为 %v, 实际为 %v", ShellBash, cmd.ShellType())
		}
	})
}

// TestArgsImmutability 测试Args返回的切片是不可变的
func TestArgsImmutability(t *testing.T) {
	cmd := NewCmd("echo", "hello", "world")

	args1 := cmd.Args()
	args2 := cmd.Args()

	// 修改返回的切片不应该影响原始数据
	args1[0] = "modified"

	if args2[0] == "modified" {
		t.Error("Args()返回的切片应该是独立的副本")
	}

	// 原始命令的参数也不应该被修改
	originalArgs := cmd.Args()
	if originalArgs[0] != "hello" {
		t.Error("修改返回的切片不应该影响原始命令参数")
	}
}

// TestEnvImmutability 测试Env返回的切片是不可变的
func TestEnvImmutability(t *testing.T) {
	cmd := NewCmd("echo", "test").WithEnv("TEST_VAR", "original")

	env1 := cmd.Env()
	env2 := cmd.Env()

	// 修改返回的切片不应该影响原始数据
	if len(env1) > 0 {
		env1[0] = "MODIFIED=value"
	}

	// 检查第二次获取的环境变量是否受影响
	found := false
	for _, env := range env2 {
		if env == "MODIFIED=value" {
			found = true
			break
		}
	}

	if found {
		t.Error("Env()返回的切片应该是独立的副本")
	}
}

// TestProcessControl 测试进程控制方法
func TestProcessControl(t *testing.T) {
	t.Run("未启动进程的控制方法", func(t *testing.T) {
		cmd := NewCmd("sleep", "10")

		// 未启动的进程
		if cmd.IsRunning() {
			t.Error("未启动的进程不应该在运行")
		}

		if cmd.GetPID() != 0 {
			t.Error("未启动的进程PID应该为0")
		}

		err := cmd.Kill()
		if err == nil {
			t.Error("杀死未启动的进程应该返回错误")
		}

		err = cmd.Signal(os.Interrupt)
		if err == nil {
			t.Error("向未启动的进程发送信号应该返回错误")
		}
	})
}

// TestCmd 测试Cmd方法
func TestCmd(t *testing.T) {
	cmd := NewCmd("echo", "test")

	execCmd := cmd.Cmd()
	if execCmd == nil {
		t.Error("Cmd()应该返回非nil的exec.Cmd对象")
	}

	// 第二次调用应该返回同一个对象
	execCmd2 := cmd.Cmd()
	if execCmd != execCmd2 {
		t.Error("多次调用Cmd()应该返回同一个对象")
	}
}
