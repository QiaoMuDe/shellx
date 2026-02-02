package shx

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	cmd := New("echo hello")

	// 测试基本执行
	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if !cmd.IsExecuted() {
		t.Fatal("IsExecuted should be true after execution")
	}
}

func TestExecOutput(t *testing.T) {
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

func TestExecContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	cmd := New("echo test")

	// 测试上下文执行
	err := cmd.ExecContext(ctx)
	if err != nil {
		t.Fatalf("ExecContext failed: %v", err)
	}

	if !cmd.IsExecuted() {
		t.Fatal("IsExecuted should be true after execution")
	}
}

func TestExecContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	cmd := New("ping -n 1 127.0.0.1")

	// 测试上下文超时
	err := cmd.ExecContext(ctx)
	if err == nil {
		t.Fatal("expected timeout error")
	}

	if !strings.Contains(err.Error(), "timed out") && !strings.Contains(err.Error(), "exit status") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecContextOutput(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	cmd := New("echo test")

	// 测试上下文执行并获取输出
	output, err := cmd.ExecContextOutput(ctx)
	if err != nil {
		t.Fatalf("ExecContextOutput failed: %v", err)
	}

	if !strings.Contains(string(output), "test") {
		t.Fatalf("unexpected output: %s", output)
	}

	if !cmd.IsExecuted() {
		t.Fatal("IsExecuted should be true after execution")
	}
}

func TestExecContextOutputTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	cmd := New("ping -n 1 127.0.0.1")

	// 测试上下文执行并获取输出超时
	_, err := cmd.ExecContextOutput(ctx)
	if err == nil {
		t.Fatal("expected timeout error")
	}

	if !strings.Contains(err.Error(), "timed out") && !strings.Contains(err.Error(), "exit status") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecWithTimeout(t *testing.T) {
	cmd := New("echo test").WithTimeout(100 * time.Millisecond)

	// 测试超时执行
	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if !cmd.IsExecuted() {
		t.Fatal("IsExecuted should be true after execution")
	}
}

func TestExecWithTimeoutActual(t *testing.T) {
	cmd := New("ping -n 1 192.0.2.1").WithTimeout(10 * time.Millisecond)

	// 测试实际超时
	err := cmd.Exec()
	if err == nil {
		t.Fatal("expected error")
	}

	// 检查是否是超时错误或退出码错误
	if !strings.Contains(err.Error(), "timed out") && !strings.Contains(err.Error(), "exit status") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecWithIO(t *testing.T) {
	var stdout, stderr bytes.Buffer
	input := strings.NewReader("test input")

	// 使用 echo 来测试 IO（因为 cat 可能不存在）
	cmd := New("echo").WithStdin(input).WithStdout(&stdout).WithStderr(&stderr)

	// 测试 IO 重定向执行
	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	// echo 不读取 stdin，但应该有输出
	if stdout.String() == "" {
		t.Fatal("expected some output from echo")
	}

	if !cmd.IsExecuted() {
		t.Fatal("IsExecuted should be true after execution")
	}
}

func TestExecMultipleTimes(t *testing.T) {
	cmd1 := New("echo test1")
	cmd2 := New("echo test2")

	// 测试多个命令执行
	err1 := cmd1.Exec()
	err2 := cmd2.Exec()

	if err1 != nil {
		t.Fatalf("cmd1.Exec failed: %v", err1)
	}

	if err2 != nil {
		t.Fatalf("cmd2.Exec failed: %v", err2)
	}

	if !cmd1.IsExecuted() {
		t.Fatal("cmd1.IsExecuted should be true")
	}

	if !cmd2.IsExecuted() {
		t.Fatal("cmd2.IsExecuted should be true")
	}
}

func TestExecAfterExecuted(t *testing.T) {
	cmd := New("echo test")

	// 先执行一次
	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
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

func TestExecOutputAfterExecuted(t *testing.T) {
	cmd := New("echo test")

	// 先执行一次
	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	// 测试重复执行
	_, Err := cmd.ExecOutput()
	if Err == nil {
		t.Fatal("expected error for duplicate execution")
	}

	if !strings.Contains(Err.Error(), "already been executed") {
		t.Fatalf("unexpected error: %v", Err)
	}
}
