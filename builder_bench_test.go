package shellx

import (
	"context"
	"strings"
	"testing"
	"time"
)

// BenchmarkNewCmd 基准测试 NewCmd 函数
func BenchmarkNewCmd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewCmd("echo", "test", "arg1", "arg2")
	}
}

// BenchmarkNewCmds 基准测试 NewCmds 函数
func BenchmarkNewCmds(b *testing.B) {
	cmdArgs := []string{"echo", "test", "arg1", "arg2"}
	for i := 0; i < b.N; i++ {
		_ = NewCmds(cmdArgs)
	}
}

// BenchmarkNewCmdStr 基准测试 NewCmdStr 函数
func BenchmarkNewCmdStr(b *testing.B) {
	cmdStr := `echo "hello world" test`
	for i := 0; i < b.N; i++ {
		_ = NewCmdStr(cmdStr)
	}
}

// BenchmarkBuilderChaining 基准测试链式调用
func BenchmarkBuilderChaining(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_ = NewCmd("echo", "test").
			WithWorkDir("/tmp").
			WithEnv("TEST", "value").
			WithTimeout(5 * time.Second).
			WithContext(ctx).
			WithShell(ShellBash)
	}
}

// BenchmarkBuilderBuild 基准测试 Build 方法
func BenchmarkBuilderBuild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := NewCmd("echo", "test").
			WithWorkDir("/tmp").
			WithEnv("TEST", "value").
			WithTimeout(5 * time.Second).
			WithShell(ShellBash)
		_ = builder.Build()
	}
}

// BenchmarkBuilderGetters 基准测试所有getter方法
func BenchmarkBuilderGetters(b *testing.B) {
	builder := NewCmd("echo", "test", "arg1", "arg2").
		WithWorkDir("/tmp").
		WithEnv("TEST", "value").
		WithTimeout(5 * time.Second).
		WithShell(ShellBash)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Name()
		_ = builder.Args()
		_ = builder.WorkDir()
		_ = builder.Env()
		_ = builder.Timeout()
		_ = builder.ShellType()
		_ = builder.Raw()
	}
}

// BenchmarkBuilderConcurrentRead 基准测试并发读取
func BenchmarkBuilderConcurrentRead(b *testing.B) {
	builder := NewCmd("echo", "test").
		WithWorkDir("/tmp").
		WithEnv("TEST", "value").
		WithTimeout(5 * time.Second).
		WithShell(ShellBash)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = builder.Name()
			_ = builder.Args()
			_ = builder.WorkDir()
			_ = builder.Env()
			_ = builder.Timeout()
			_ = builder.ShellType()
		}
	})
}

// BenchmarkBuilderWithEnv 基准测试环境变量设置
func BenchmarkBuilderWithEnv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := NewCmd("echo", "test")
		for j := 0; j < 10; j++ {
			builder.WithEnv("TEST"+string(rune(j)), "value")
		}
	}
}

// BenchmarkBuilderWithStdio 基准测试标准输入输出设置
func BenchmarkBuilderWithStdio(b *testing.B) {
	stdin := strings.NewReader("test input")

	for i := 0; i < b.N; i++ {
		builder := NewCmd("cat")
		builder.WithStdin(stdin)
	}
}
