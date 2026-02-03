package shx

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWithDir(t *testing.T) {
	cmd := New("echo test")

	// 测试设置目录
	result := cmd.WithDir(".")
	if result == nil {
		t.Fatal("WithDir returned nil")
	}

	if result.Dir() != cmd.Dir() {
		t.Fatal("WithDir should return the same instance")
	}
}

func TestWithDirInvalid(t *testing.T) {
	cmd := New("echo test")

	// 测试无效目录
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid directory")
		}
	}()

	cmd.WithDir("/nonexistent/directory/path")
}

func TestWithEnv(t *testing.T) {
	cmd := New("echo $TEST")

	// 测试设置环境变量
	result := cmd.WithEnv("TEST", "value")
	if result == nil {
		t.Fatal("WithEnv returned nil")
	}

	if result != cmd {
		t.Fatal("WithEnv should return the same instance")
	}

	// 验证环境变量设置
	output, err := result.ExecOutput()
	if err != nil {
		t.Fatalf("ExecOutput failed: %v", err)
	}

	if !strings.Contains(string(output), "value") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWithEnvs(t *testing.T) {
	cmd := New("echo $VAR1 $VAR2")

	envs := map[string]string{
		"VAR1": "value1",
		"VAR2": "value2",
	}

	// 测试批量设置环境变量
	result := cmd.WithEnvMap(envs)
	if result == nil {
		t.Fatal("WithEnvs returned nil")
	}

	if result != cmd {
		t.Fatal("WithEnvs should return the same instance")
	}

	// 验证环境变量设置
	output, err := result.ExecOutput()
	if err != nil {
		t.Fatalf("ExecOutput failed: %v", err)
	}

	if !strings.Contains(string(output), "value1") || !strings.Contains(string(output), "value2") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWithStdin(t *testing.T) {
	cmd := New("echo").WithStdin(strings.NewReader("test input"))

	var stdout bytes.Buffer

	// 测试设置标准输入
	result := cmd.WithStdout(&stdout)
	if result == nil {
		t.Fatal("WithStdout returned nil")
	}

	if result != cmd {
		t.Fatal("WithStdout should return the same instance")
	}

	// 验证输出
	err := result.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	// echo 不读取 stdin，但应该有输出
	if stdout.String() == "" {
		t.Fatal("expected some output from echo")
	}
}

func TestWithStdout(t *testing.T) {
	cmd := New("echo hello")

	var stdout bytes.Buffer

	// 测试设置标准输出
	result := cmd.WithStdout(&stdout)
	if result == nil {
		t.Fatal("WithStdout returned nil")
	}

	if result != cmd {
		t.Fatal("WithStdout should return the same instance")
	}

	// 验证输出
	err := result.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if !strings.Contains(stdout.String(), "hello") {
		t.Fatalf("unexpected stdout: %s", stdout.String())
	}
}

func TestWithStderr(t *testing.T) {
	cmd := New("echo error 1>&2")

	var stderr bytes.Buffer

	// 测试设置标准错误
	result := cmd.WithStderr(&stderr)
	if result == nil {
		t.Fatal("WithStderr returned nil")
	}

	if result != cmd {
		t.Fatal("WithStderr should return the same instance")
	}

	// 验证错误输出
	err := result.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if !strings.Contains(stderr.String(), "error") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestWithTimeout(t *testing.T) {
	cmd := New("echo test")

	// 测试设置超时
	result := cmd.WithTimeout(5 * time.Second)
	if result == nil {
		t.Fatal("WithTimeout returned nil")
	}

	if result != cmd {
		t.Fatal("WithTimeout should return the same instance")
	}

	if result.Timeout() != 5*time.Second {
		t.Fatalf("expected timeout 5s, got %v", result.Timeout())
	}
}

func TestWithTimeoutZero(t *testing.T) {
	cmd := New("echo test")
	originalTimeout := cmd.Timeout()

	// 测试设置零超时
	result := cmd.WithTimeout(0)
	if result == nil {
		t.Fatal("WithTimeout returned nil")
	}

	if result.Timeout() != originalTimeout {
		t.Fatal("WithTimeout with 0 should not change timeout")
	}
}

func TestWithContext(t *testing.T) {
	cmd := New("echo test")
	ctx, cancel := createTestContext()
	defer cancel()

	// 测试设置上下文
	result := cmd.WithContext(ctx)
	if result == nil {
		t.Fatal("WithContext returned nil")
	}

	if result != cmd {
		t.Fatal("WithContext should return the same instance")
	}

	if result.Context() != ctx {
		t.Fatal("Context not set correctly")
	}
}

func TestWithNilContext(t *testing.T) {
	cmd := New("echo test")

	// 测试设置 nil 上下文应该 panic
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil context")
		}
	}()

	//nolint:all
	cmd.WithContext(nil)
}

func TestWithAfterExecution(t *testing.T) {
	cmd := New("echo test")

	// 先执行
	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	// 测试执行后设置选项应该 panic
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for option after execution")
		}
	}()

	cmd.WithTimeout(5 * time.Second)
}
