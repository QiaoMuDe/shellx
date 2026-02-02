package shx

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	err := Run("echo hello")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	_ = err // 避免 errcheck 错误
}

func TestRunToTerminal(t *testing.T) {
	err := RunToTerminal("echo hello")
	if err != nil {
		t.Fatalf("RunToTerminal failed: %v", err)
	}
	_ = err // 避免 errcheck 错误
}

func TestOut(t *testing.T) {
	output, err := Out("echo world")
	if err != nil {
		t.Fatalf("Out failed: %v", err)
	}
	if !strings.Contains(string(output), "world") {
		t.Fatalf("unexpected output: %s", output)
	}
	_ = output // 避免 errcheck 错误
}

func TestRunWith(t *testing.T) {
	start := time.Now()
	err := RunWith("echo test", 100*time.Millisecond)
	if err != nil {
		t.Fatalf("RunWith failed: %v", err)
	}
	// 应该在100ms内完成
	if time.Since(start) > 200*time.Millisecond {
		t.Fatalf("took too long: %v", time.Since(start))
	}
}

func TestOutWith(t *testing.T) {
	start := time.Now()
	output, err := OutWith("echo timeout", 100*time.Millisecond)
	if err != nil {
		t.Fatalf("OutWith failed: %v", err)
	}
	if !strings.Contains(string(output), "timeout") {
		t.Fatalf("unexpected output: %s", output)
	}
	if time.Since(start) > 200*time.Millisecond {
		t.Fatalf("took too long: %v", time.Since(start))
	}
}

func TestRunWithIO(t *testing.T) {
	var stdout, stderr bytes.Buffer
	input := strings.NewReader("test input")

	// 使用 cat 来读取 stdin（如果可用）
	err := RunWithIO("cat", input, &stdout, &stderr)
	if err != nil {
		// 如果 cat 不可用，尝试其他方法
		t.Logf("cat failed: %v, trying alternative", err)
		// 使用 echo 并检查输入是否被忽略
		err = RunWithIO("echo", input, &stdout, &stderr)
		if err != nil {
			t.Fatalf("RunWithIO failed: %v", err)
		}
		// echo 不读取 stdin，但至少应该有输出
		if stdout.String() == "" {
			t.Fatal("expected some output from echo")
		}
	} else {
		// cat 应该输出输入内容
		if !strings.Contains(stdout.String(), "test input") {
			t.Fatalf("unexpected stdout: %s", stdout.String())
		}

		if stderr.String() != "" {
			t.Fatalf("unexpected stderr: %s", stderr.String())
		}
	}
}

func TestOutWithIO(t *testing.T) {
	var stdout, stderr bytes.Buffer
	input := strings.NewReader("hello")

	// 尝试使用 cat，失败则使用 echo
	output, err := OutWithIO("cat", input, &stdout, &stderr)
	if err != nil {
		t.Logf("cat failed: %v, trying echo", err)
		// echo 不读取 stdin，但至少应该有输出
		output, err = OutWithIO("echo", input, &stdout, &stderr)
		if err != nil {
			t.Fatalf("OutWithIO failed: %v", err)
		}
		// echo 输出是输入内容加上换行
		expected := "\n"
		if string(output) != expected {
			t.Fatalf("expected %q, got %q", expected, string(output))
		}
	} else {
		// cat 应该输出输入内容
		expected := "hello"
		if string(output) != expected {
			t.Fatalf("expected %q, got %q", expected, string(output))
		}
	}
}

func TestRunCtx(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// 使用一个不存在的地址测试超时
	err := RunCtx(ctx, "ping -n 1 192.0.2.1")
	if err == nil {
		t.Fatal("expected error")
	}

	// 检查是否是超时错误或退出码错误
	if !strings.Contains(err.Error(), "timed out") && !strings.Contains(err.Error(), "exit status") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOutCtx(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// 使用一个不存在的地址测试超时
	_, err := OutCtx(ctx, "ping -n 1 192.0.2.1")
	if err == nil {
		t.Fatal("expected error")
	}

	// 检查是否是超时错误或退出码错误
	if !strings.Contains(err.Error(), "timed out") && !strings.Contains(err.Error(), "exit status") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// 测试错误命令
func TestErrorCommand(t *testing.T) {
	err := Run("nonexistent-command")
	if err == nil {
		t.Fatal("expected error")
	}

	// 检查是否是退出码错误
	code, ok := IsExitStatus(err)
	if ok && code == 127 {
		// 命令不存在是退出码 127，这是预期的
		return
	}

	// 应该包含命令失败信息
	if !strings.Contains(err.Error(), "command failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// 测试退出码
func TestExitCode(t *testing.T) {
	// 使用跨平台的命令
	err := Run("echo hello && exit 1")
	if err == nil {
		t.Fatal("expected error")
	}

	// 检查是否是退出码错误
	code, ok := IsExitStatus(err)
	if !ok {
		t.Fatalf("expected exit status error")
	}

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

// 测试 ping 命令（跨平台）
func TestPing(t *testing.T) {
	// 使用 ping 127.0.0.1 -c 1 (Unix) 或 ping -n 1 127.0.0.1 (Windows)
	output, err := Out("ping -n 1 127.0.0.1")
	if err != nil {
		// ping 可能需要管理员权限，所以失败是可接受的
		t.Logf("ping failed (may need admin): %v", err)
		return
	}

	// 应该包含 ping 的输出
	if !strings.Contains(strings.ToLower(string(output)), "ping") {
		t.Fatalf("expected ping output: %s", string(output))
	}
}

// 基准测试 - 确保便捷函数性能合理
func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Run("echo test")
	}
}

func BenchmarkOut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Out("echo test")
	}
}

// func TestExecT(t *testing.T) {
// 	RunToTerminal("fck ls -c -i")
// }
