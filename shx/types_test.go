package shx

import (
	"strings"
	"testing"
	"time"
)

func TestExitStatus(t *testing.T) {
	// 测试 ExitStatus 结构体
	err := ExitStatus{Code: 42}

	expected := "exit status 42"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestExitStatusZero(t *testing.T) {
	// 测试退出码为 0
	err := ExitStatus{Code: 0}

	expected := "exit status 0"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestExitStatusMax(t *testing.T) {
	// 测试最大退出码
	err := ExitStatus{Code: 255}

	expected := "exit status 255"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestShxInitialValues(t *testing.T) {
	cmd := New("echo test")

	// 测试初始值
	if cmd.Raw() != "echo test" {
		t.Fatalf("expected raw 'echo test', got %q", cmd.Raw())
	}

	if cmd.Dir() == "" {
		t.Fatal("Dir should not be empty")
	}

	if cmd.Env() == nil {
		t.Fatal("Env should not be nil")
	}

	if cmd.Timeout() != 0 {
		t.Fatalf("expected timeout 0, got %v", cmd.Timeout())
	}

	if cmd.Context() != nil {
		t.Fatal("Context should be nil by default")
	}

	if cmd.IsExecuted() {
		t.Fatal("IsExecuted should be false initially")
	}
}

func TestShxWithValues(t *testing.T) {
	cmd := New("echo $TEST").
		WithTimeout(5*time.Second).
		WithEnv("TEST", "value")

	// 测试设置的值
	if cmd.Timeout() != 5*time.Second {
		t.Fatalf("expected timeout 5s, got %v", cmd.Timeout())
	}

	if cmd.IsExecuted() {
		t.Fatal("IsExecuted should be false before execution")
	}

	// 验证环境变量设置
	output, err := cmd.ExecOutput()
	if err != nil {
		t.Fatalf("ExecOutput failed: %v", err)
	}

	// 环境变量应该被展开
	if !strings.Contains(string(output), "value") {
		t.Fatalf("unexpected output: %s", output)
	}

	if !cmd.IsExecuted() {
		t.Fatal("IsExecuted should be true after execution")
	}
}
