package shellx

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNewCmd 测试 NewCmd 函数
func TestNewCmd(t *testing.T) {
	tests := []struct {
		name      string
		cmdName   string
		args      []string
		wantName  string
		wantArgs  []string
		wantPanic bool
	}{
		{
			name:      "正常创建命令",
			cmdName:   "ls",
			args:      []string{"-l", "-a"},
			wantName:  "ls",
			wantArgs:  []string{"-l", "-a"},
			wantPanic: false,
		},
		{
			name:      "无参数命令",
			cmdName:   "pwd",
			args:      []string{},
			wantName:  "pwd",
			wantArgs:  []string{},
			wantPanic: false,
		},
		{
			name:      "空命令名应该panic",
			cmdName:   "",
			args:      []string{},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewCmd() 应该panic但没有")
					}
				}()
			}

			builder := NewCmd(tt.cmdName, tt.args...)

			if !tt.wantPanic {
				if builder.Name() != tt.wantName {
					t.Errorf("NewCmd() name = %v, want %v", builder.Name(), tt.wantName)
				}

				if !reflect.DeepEqual(builder.Args(), tt.wantArgs) {
					t.Errorf("NewCmd() args = %v, want %v", builder.Args(), tt.wantArgs)
				}

				// 检查默认值
				if builder.ShellType() != ShellDefault {
					t.Errorf("NewCmd() shellType = %v, want %v", builder.ShellType(), ShellDefault)
				}

				if len(builder.Env()) == 0 {
					t.Errorf("NewCmd() 应该继承父进程环境变量")
				}
			}
		})
	}
}

// TestNewCmds 测试 NewCmds 函数
func TestNewCmds(t *testing.T) {
	tests := []struct {
		name     string
		cmdArgs  []string
		wantName string
		wantArgs []string
	}{
		{
			name:     "正常命令数组",
			cmdArgs:  []string{"git", "status", "--short"},
			wantName: "git",
			wantArgs: []string{"status", "--short"},
		},
		{
			name:     "单个命令",
			cmdArgs:  []string{"pwd"},
			wantName: "pwd",
			wantArgs: []string{},
		},
		{
			name:     "空数组",
			cmdArgs:  []string{},
			wantName: "",
			wantArgs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewCmds(tt.cmdArgs)

			if builder.Name() != tt.wantName {
				t.Errorf("NewCmds() name = %v, want %v", builder.Name(), tt.wantName)
			}

			if !reflect.DeepEqual(builder.Args(), tt.wantArgs) {
				t.Errorf("NewCmds() args = %v, want %v", builder.Args(), tt.wantArgs)
			}
		})
	}
}

// TestNewCmdStr 测试 NewCmdStr 函数
func TestNewCmdStr(t *testing.T) {
	tests := []struct {
		name     string
		cmdStr   string
		wantName string
		wantArgs []string
		wantRaw  string
	}{
		{
			name:     "简单命令字符串",
			cmdStr:   "ls -l -a",
			wantName: "ls",
			wantArgs: []string{"-l", "-a"},
			wantRaw:  "ls -l -a",
		},
		{
			name:     "带引号的命令",
			cmdStr:   `echo "hello world"`,
			wantName: "echo",
			wantArgs: []string{"hello world"},
			wantRaw:  `echo "hello world"`,
		},
		{
			name:     "空字符串",
			cmdStr:   "",
			wantName: "",
			wantArgs: []string{},
			wantRaw:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewCmdStr(tt.cmdStr)

			if builder.Name() != tt.wantName {
				t.Errorf("NewCmdStr() name = %v, want %v", builder.Name(), tt.wantName)
			}

			if !reflect.DeepEqual(builder.Args(), tt.wantArgs) {
				t.Errorf("NewCmdStr() args = %v, want %v", builder.Args(), tt.wantArgs)
			}

			if builder.Raw() != tt.wantRaw {
				t.Errorf("NewCmdStr() raw = %v, want %v", builder.Raw(), tt.wantRaw)
			}
		})
	}
}

// TestBuilderWithWorkDir 测试 WithWorkDir 方法
func TestBuilderWithWorkDir(t *testing.T) {
	builder := NewCmd("ls")

	// 测试设置工作目录
	result := builder.WithWorkDir("/tmp")
	if result != builder {
		t.Errorf("WithWorkDir() 应该返回同一个builder实例")
	}

	if builder.WorkDir() != "/tmp" {
		t.Errorf("WithWorkDir() workDir = %v, want %v", builder.WorkDir(), "/tmp")
	}

	// 测试空目录
	builder.WithWorkDir("")
	if builder.WorkDir() != "/tmp" {
		t.Errorf("WithWorkDir(\"\") 不应该改变工作目录")
	}
}

// TestBuilderWithEnv 测试 WithEnv 方法
func TestBuilderWithEnv(t *testing.T) {
	builder := NewCmd("ls")
	originalEnvCount := len(builder.Env())

	// 测试添加环境变量
	result := builder.WithEnv("TEST_VAR", "test_value")
	if result != builder {
		t.Errorf("WithEnv() 应该返回同一个builder实例")
	}

	env := builder.Env()
	if len(env) != originalEnvCount+1 {
		t.Errorf("WithEnv() 环境变量数量 = %v, want %v", len(env), originalEnvCount+1)
	}

	found := false
	for _, e := range env {
		if e == "TEST_VAR=test_value" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("WithEnv() 没有找到设置的环境变量")
	}

	// 测试空key
	builder.WithEnv("", "value")
	if len(builder.Env()) != originalEnvCount+1 {
		t.Errorf("WithEnv(\"\", \"value\") 不应该添加环境变量")
	}
}

// TestBuilderWithTimeout 测试 WithTimeout 方法
func TestBuilderWithTimeout(t *testing.T) {
	builder := NewCmd("ls")

	// 测试设置超时时间
	timeout := 5 * time.Second
	result := builder.WithTimeout(timeout)
	if result != builder {
		t.Errorf("WithTimeout() 应该返回同一个builder实例")
	}

	if builder.Timeout() != timeout {
		t.Errorf("WithTimeout() timeout = %v, want %v", builder.Timeout(), timeout)
	}

	// 测试负数或零超时
	builder.WithTimeout(-1 * time.Second)
	if builder.Timeout() != timeout {
		t.Errorf("WithTimeout(-1) 不应该改变超时时间")
	}

	builder.WithTimeout(0)
	if builder.Timeout() != timeout {
		t.Errorf("WithTimeout(0) 不应该改变超时时间")
	}
}

// TestBuilderWithContext 测试 WithContext 方法
func TestBuilderWithContext(t *testing.T) {
	builder := NewCmd("ls")
	ctx := context.Background()

	// 测试设置上下文
	result := builder.WithContext(ctx)
	if result != builder {
		t.Errorf("WithContext() 应该返回同一个builder实例")
	}

	// 测试nil上下文应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("WithContext(nil) 应该panic但没有")
		}
	}()
	//nolint:all
	builder.WithContext(nil)
}

// TestBuilderWithStdin 测试 WithStdin 方法
func TestBuilderWithStdin(t *testing.T) {
	builder := NewCmd("cat")
	stdin := strings.NewReader("test input")

	// 测试设置标准输入
	result := builder.WithStdin(stdin)
	if result != builder {
		t.Errorf("WithStdin() 应该返回同一个builder实例")
	}

	// 测试nil输入应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("WithStdin(nil) 应该panic但没有")
		}
	}()
	builder.WithStdin(nil)
}

// TestBuilderWithStdout 测试 WithStdout 方法
func TestBuilderWithStdout(t *testing.T) {
	builder := NewCmd("echo", "test")
	var stdout bytes.Buffer

	// 测试设置标准输出
	result := builder.WithStdout(&stdout)
	if result != builder {
		t.Errorf("WithStdout() 应该返回同一个builder实例")
	}

	// 测试nil输出应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("WithStdout(nil) 应该panic但没有")
		}
	}()
	builder.WithStdout(nil)
}

// TestBuilderWithStderr 测试 WithStderr 方法
func TestBuilderWithStderr(t *testing.T) {
	builder := NewCmd("ls", "/nonexistent")
	var stderr bytes.Buffer

	// 测试设置标准错误输出
	result := builder.WithStderr(&stderr)
	if result != builder {
		t.Errorf("WithStderr() 应该返回同一个builder实例")
	}

	// 测试nil输出应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("WithStderr(nil) 应该panic但没有")
		}
	}()
	builder.WithStderr(nil)
}

// TestBuilderWithShell 测试 WithShell 方法
func TestBuilderWithShell(t *testing.T) {
	builder := NewCmd("ls")

	// 测试设置shell类型
	result := builder.WithShell(ShellBash)
	if result != builder {
		t.Errorf("WithShell() 应该返回同一个builder实例")
	}

	if builder.ShellType() != ShellBash {
		t.Errorf("WithShell() shellType = %v, want %v", builder.ShellType(), ShellBash)
	}
}

// TestBuilderChaining 测试链式调用
func TestBuilderChaining(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ctx := context.Background()

	builder := NewCmd("echo", "test").
		WithWorkDir("/tmp").
		WithEnv("TEST", "value").
		WithTimeout(5 * time.Second).
		WithContext(ctx).
		WithStdout(&stdout).
		WithStderr(&stderr).
		WithShell(ShellBash)

	// 验证所有设置都生效
	if builder.Name() != "echo" {
		t.Errorf("链式调用后 name = %v, want %v", builder.Name(), "echo")
	}

	if !reflect.DeepEqual(builder.Args(), []string{"test"}) {
		t.Errorf("链式调用后 args = %v, want %v", builder.Args(), []string{"test"})
	}

	if builder.WorkDir() != "/tmp" {
		t.Errorf("链式调用后 workDir = %v, want %v", builder.WorkDir(), "/tmp")
	}

	if builder.Timeout() != 5*time.Second {
		t.Errorf("链式调用后 timeout = %v, want %v", builder.Timeout(), 5*time.Second)
	}

	if builder.ShellType() != ShellBash {
		t.Errorf("链式调用后 shellType = %v, want %v", builder.ShellType(), ShellBash)
	}

	// 检查环境变量
	env := builder.Env()
	found := false
	for _, e := range env {
		if e == "TEST=value" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("链式调用后没有找到设置的环境变量")
	}
}

// TestBuilderBuild 测试 Build 方法
func TestBuilderBuild(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() *Builder
		wantNil   bool
	}{
		{
			name: "基本构建",
			setupFunc: func() *Builder {
				return NewCmd("echo", "test")
			},
			wantNil: false,
		},
		{
			name: "带上下文构建",
			setupFunc: func() *Builder {
				return NewCmd("echo", "test").WithContext(context.Background())
			},
			wantNil: false,
		},
		{
			name: "使用shell构建",
			setupFunc: func() *Builder {
				return NewCmd("echo", "test").WithShell(ShellBash)
			},
			wantNil: false,
		},
		{
			name: "不使用shell构建",
			setupFunc: func() *Builder {
				return NewCmd("echo", "test").WithShell(ShellNone)
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.setupFunc()
			cmd := builder.Build()

			if (cmd == nil) != tt.wantNil {
				t.Errorf("Build() = %v, wantNil %v", cmd, tt.wantNil)
			}

			if cmd != nil {
				if cmd.Cmd() == nil {
					t.Errorf("Build() 返回的Command对象的Cmd()不应该为nil")
				}

				if cmd.IsExecuted() {
					t.Errorf("Build() 返回的Command对象不应该已执行")
				}
			}
		})
	}
}

// TestBuilderConcurrency 测试并发安全性
func TestBuilderConcurrency(t *testing.T) {
	builder := NewCmd("echo", "test")

	var wg sync.WaitGroup
	const numGoroutines = 100

	// 并发读取
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_ = builder.Name()
			_ = builder.Args()
			_ = builder.WorkDir()
			_ = builder.Env()
			_ = builder.Timeout()
			_ = builder.ShellType()
			_ = builder.Raw()
		}()
	}

	// 并发写入
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			builder.WithEnv("TEST"+string(rune(i)), "value")
			builder.WithWorkDir("/tmp")
			builder.WithTimeout(time.Duration(i) * time.Second)
		}(i)
	}

	wg.Wait()

	// 验证最终状态是一致的
	if builder.WorkDir() != "/tmp" {
		t.Errorf("并发测试后 workDir = %v, want %v", builder.WorkDir(), "/tmp")
	}
}

// TestBuilderGetters 测试所有getter方法
func TestBuilderGetters(t *testing.T) {
	var stdout, stderr bytes.Buffer
	builder := NewCmd("test", "arg1", "arg2").
		WithWorkDir("/test").
		WithEnv("TEST", "value").
		WithTimeout(10 * time.Second).
		WithStdout(&stdout).
		WithStderr(&stderr).
		WithShell(ShellPwsh)

	// 测试Name()
	if builder.Name() != "test" {
		t.Errorf("Name() = %v, want %v", builder.Name(), "test")
	}

	// 测试Args() - 应该返回副本
	args1 := builder.Args()
	args2 := builder.Args()
	if &args1[0] == &args2[0] {
		t.Errorf("Args() 应该返回副本而不是原始切片")
	}

	// 测试WorkDir()
	if builder.WorkDir() != "/test" {
		t.Errorf("WorkDir() = %v, want %v", builder.WorkDir(), "/test")
	}

	// 测试Env() - 应该返回副本
	env1 := builder.Env()
	env2 := builder.Env()
	if len(env1) > 0 && len(env2) > 0 && &env1[0] == &env2[0] {
		t.Errorf("Env() 应该返回副本而不是原始切片")
	}

	// 测试Timeout()
	if builder.Timeout() != 10*time.Second {
		t.Errorf("Timeout() = %v, want %v", builder.Timeout(), 10*time.Second)
	}

	// 测试ShellType()
	if builder.ShellType() != ShellPwsh {
		t.Errorf("ShellType() = %v, want %v", builder.ShellType(), ShellPwsh)
	}

	// 测试Raw()
	if builder.Raw() != "" {
		t.Errorf("Raw() = %v, want empty string", builder.Raw())
	}
}

// TestBuilderEdgeCases 测试边界情况
func TestBuilderEdgeCases(t *testing.T) {
	// 测试空参数的Args()
	builder := NewCmd("test")
	args := builder.Args()
	if len(args) != 0 {
		t.Errorf("无参数命令的Args() = %v, want empty slice", args)
	}

	// 测试修改返回的Args切片不影响原始数据
	builder = NewCmd("test", "arg1", "arg2")
	args = builder.Args()
	args[0] = "modified"
	if builder.Args()[0] == "modified" {
		t.Errorf("修改Args()返回的切片不应该影响原始数据")
	}

	// 测试修改返回的Env切片不影响原始数据
	builder.WithEnv("TEST", "value")
	env := builder.Env()
	if len(env) > 0 {
		original := env[0]
		env[0] = "modified"
		if builder.Env()[0] == "modified" {
			t.Errorf("修改Env()返回的切片不应该影响原始数据")
		}
		// 恢复原值以便后续测试
		env[0] = original
	}
}
