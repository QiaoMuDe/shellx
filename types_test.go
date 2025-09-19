// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件包含shellx包中所有类型的单元测试，验证类型定义和方法的正确性。
package shellx

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestShellType(t *testing.T) {
	tests := []struct {
		shell    ShellType
		expected string
	}{
		{ShellSh, "sh"},
		{ShellBash, "bash"},
		{ShellPwsh, "pwsh"},
		{ShellPowerShell, "powershell"},
		{ShellCmd, "cmd"},
	}

	for _, tt := range tests {
		if got := tt.shell.String(); got != tt.expected {
			t.Errorf("ShellType.String() = %v, want %v", got, tt.expected)
		}
	}
}

func TestBuilder(t *testing.T) {
	// 测试 NewCmd (可变参数方式)
	cmd1 := NewCmd("ls", "-la", "-h").
		WithWorkDir("/tmp").
		WithTimeout(30*time.Second).
		WithEnv("PATH", "/usr/bin").
		Build()

	if cmd1.Name() != "ls" {
		t.Errorf("Expected name 'ls', got '%s'", cmd1.Name())
	}
	if len(cmd1.Args()) != 2 {
		t.Errorf("Expected 2 args, got %d", len(cmd1.Args()))
	}
	if cmd1.Dir() != "/tmp" {
		t.Errorf("Expected workDir '/tmp', got '%s'", cmd1.Dir())
	}

	// 测试 NewCmds (切片方式)
	cmdArgs := []string{"git", "commit", "-m", "test"}
	cmd2 := NewCmds(cmdArgs).
		WithContext(context.Background()).
		Build()

	if cmd2.Name() != "git" {
		t.Errorf("Expected name 'git', got '%s'", cmd2.Name())
	}
	if len(cmd2.Args()) != 3 {
		t.Errorf("Expected 3 args, got %d", len(cmd2.Args()))
	}

	// 测试 NewCmdString (字符串方式)
	cmd3 := NewCmdString("ps aux | grep go").
		WithStdin(strings.NewReader("input")).
		Build()

	if cmd3.Raw() != "ps aux | grep go" {
		t.Errorf("Expected raw 'ps aux | grep go', got '%s'", cmd3.Raw())
	}
}

func TestExecuteOptions(t *testing.T) {
	opts := &ExecuteOptions{
		Shell:   ShellBash,
		Capture: true,
	}

	if opts.Shell != ShellBash {
		t.Errorf("Expected ShellBash, got %v", opts.Shell)
	}
	if !opts.Capture {
		t.Errorf("Expected Capture to be true")
	}
}

func TestErrors(t *testing.T) {
	cmd := NewCmd("test").Build()

	// 测试 ExecutionError
	execErr := &ExecutionError{
		Cmd:      cmd,
		ExitCode: 1,
		Stderr:   "command failed",
	}
	if !strings.Contains(execErr.Error(), "test") {
		t.Errorf("ExecutionError should contain command name")
	}

	// 测试 TimeoutError
	timeoutErr := &TimeoutError{
		Cmd:     cmd,
		Timeout: 30 * time.Second,
	}
	if !strings.Contains(timeoutErr.Error(), "timed out") {
		t.Errorf("TimeoutError should contain 'timed out'")
	}

	// 测试 ValidationError
	validErr := &ValidationError{
		Field:   "name",
		Message: "cannot be empty",
	}
	if !strings.Contains(validErr.Error(), "name") {
		t.Errorf("ValidationError should contain field name")
	}
}
