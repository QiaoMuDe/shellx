package shx

import (
	"context"
	"strings"
	"testing"
	"time"

	"mvdan.cc/sh/v3/syntax"
)

func TestNew(t *testing.T) {
	cmd := New("echo hello")
	if cmd == nil {
		t.Fatal("New returned nil")
	}

	if cmd.Raw() != "echo hello" {
		t.Fatalf("unexpected raw command: %s", cmd.Raw())
	}
}

func TestNewWithParser(t *testing.T) {
	parser := syntax.NewParser()
	cmd := NewWithParser("echo world", parser)
	if cmd == nil {
		t.Fatal("NewWithParser returned nil")
	}

	if cmd.Raw() != "echo world" {
		t.Fatalf("unexpected raw command: %s", cmd.Raw())
	}
}

func TestShxMethods(t *testing.T) {
	cmd := New("echo test")

	// 测试 getter 方法
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

func TestShxExecution(t *testing.T) {
	cmd := New("echo hello")

	// 测试执行
	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if !cmd.IsExecuted() {
		t.Fatal("IsExecuted should be true after execution")
	}

	// 测试重复执行
	err = cmd.Exec()
	if err == nil {
		t.Fatal("expected error for duplicate execution")
	}

	if !strings.Contains(err.Error(), "already been executed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShxExecutionWithOutput(t *testing.T) {
	cmd := New("echo world")

	// 测试执行并获取输出
	output, err := cmd.ExecOutput()
	if err != nil {
		t.Fatalf("ExecOutput failed: %v", err)
	}

	if !strings.Contains(string(output), "world") {
		t.Fatalf("unexpected output: %s", output)
	}

	if !cmd.IsExecuted() {
		t.Fatal("IsExecuted should be true after execution")
	}
}

func TestShxWithTimeout(t *testing.T) {
	cmd := New("echo test").WithTimeout(100 * time.Millisecond)

	if cmd.Timeout() != 100*time.Millisecond {
		t.Fatalf("expected timeout 100ms, got %v", cmd.Timeout())
	}

	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
}

func TestShxWithContext(t *testing.T) {
	ctx, cancel := createTestContext()
	defer cancel()

	cmd := New("echo test").WithContext(ctx)

	if cmd.Context() != ctx {
		t.Fatal("Context not set correctly")
	}

	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
}

func TestShxWithEnv(t *testing.T) {
	cmd := New("echo $TEST_VAR").WithEnv("TEST_VAR", "test_value")

	output, err := cmd.ExecOutput()
	if err != nil {
		t.Fatalf("ExecOutput failed: %v", err)
	}

	if !strings.Contains(string(output), "test_value") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestShxWithEnvs(t *testing.T) {
	envs := map[string]string{
		"VAR1": "value1",
		"VAR2": "value2",
	}
	cmd := New("echo $VAR1 $VAR2").WithEnvMap(envs)

	output, err := cmd.ExecOutput()
	if err != nil {
		t.Fatalf("ExecOutput failed: %v", err)
	}

	if !strings.Contains(string(output), "value1") || !strings.Contains(string(output), "value2") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestShxWithDir(t *testing.T) {
	// 使用当前目录
	cmd := New("echo $PWD").WithDir(".")

	output, err := cmd.ExecOutput()
	if err != nil {
		t.Fatalf("ExecOutput failed: %v", err)
	}

	// 应该包含当前目录路径
	if len(output) == 0 {
		t.Fatal("expected non-empty output")
	}
}

// 辅助函数
func createTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
