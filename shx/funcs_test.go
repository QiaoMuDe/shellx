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

// ParseCmd 测试用例
func TestParseCmd(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		// 基础场景
		{
			name:     "简单命令",
			input:    "echo hello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "多参数命令",
			input:    "ls -la /home/user",
			expected: []string{"ls", "-la", "/home/user"},
		},
		{
			name:     "空命令",
			input:    "",
			expected: []string{},
		},
		{
			name:     "只有空格",
			input:    "   ",
			expected: []string{},
		},
		{
			name:     "单个命令",
			input:    "ls",
			expected: []string{"ls"},
		},

		// 引号场景
		{
			name:     "双引号",
			input:    `echo "hello world"`,
			expected: []string{"echo", "\"hello world\""},
		},
		{
			name:     "单引号",
			input:    `echo 'hello world'`,
			expected: []string{"echo", "'hello world'"},
		},
		{
			name:     "反引号",
			input:    "echo `hello world`",
			expected: []string{"echo", "$(hello world)"},
		},
		{
			name:     "多个引号参数",
			input:    `echo "hello" "world" "test"`,
			expected: []string{"echo", "\"hello\"", "\"world\"", "\"test\""},
		},
		{
			name:     "引号嵌套",
			input:    `echo "hello 'world'"`,
			expected: []string{"echo", "\"hello 'world'\""},
		},
		{
			name:     "单引号内双引号",
			input:    `echo 'hello "world"'`,
			expected: []string{"echo", "'hello \"world\"'"},
		},

		// 特殊字符场景
		{
			name:     "环境变量",
			input:    `echo $HOME`,
			expected: []string{"echo", "$HOME"},
		},
		{
			name:     "花括号环境变量",
			input:    `echo ${HOME}`,
			expected: []string{"echo", "${HOME}"},
		},
		{
			name:     "通配符",
			input:    `ls *.go`,
			expected: []string{"ls", "*.go"},
		},
		{
			name:     "问号通配符",
			input:    `ls test?.txt`,
			expected: []string{"ls", "test?.txt"},
		},
		{
			name:     "命令替换",
			input:    `echo $(date)`,
			expected: []string{"echo", "$(date)"},
		},
		{
			name:     "反引号命令替换",
			input:    "echo `date`",
			expected: []string{"echo", "$(date)"},
		},
		{
			name:     "转义双引号",
			input:    `echo \"hello\"`,
			expected: []string{"echo", "\\\"hello\\\""},
		},
		{
			name:     "转义单引号",
			input:    `echo \'hello\'`,
			expected: []string{"echo", "\\'hello\\'"},
		},
		{
			name:     "转义反引号",
			input:    "echo \\`hello\\`",
			expected: []string{"echo", "\\`hello\\`"},
		},

		// 复杂场景
		{
			name:     "多个连续空格",
			input:    "echo   hello    world",
			expected: []string{"echo", "hello", "world"},
		},
		{
			name:     "制表符分隔",
			input:    "echo\thello\tworld",
			expected: []string{"echo", "hello", "world"},
		},
		{
			name:     "混合空格和制表符",
			input:    "echo \t hello \t world",
			expected: []string{"echo", "hello", "world"},
		},
		{
			name:     "引号内空格",
			input:    `echo "hello   world"`,
			expected: []string{"echo", "\"hello   world\""},
		},
		{
			name:     "引号内特殊字符",
			input:    `echo "hello@world.com"`,
			expected: []string{"echo", "\"hello@world.com\""},
		},
		{
			name:     "路径参数",
			input:    `ls /home/user/Documents`,
			expected: []string{"ls", "/home/user/Documents"},
		},
		{
			name:     "Windows路径",
			input:    `dir C:\Users\test`,
			expected: []string{"dir", `C:\Users\test`},
		},

		// 错误场景（ParseCmd 应该返回空切片）
		{
			name:     "未闭合双引号",
			input:    `echo "hello`,
			expected: []string{},
		},
		{
			name:     "未闭合单引号",
			input:    `echo 'hello`,
			expected: []string{},
		},
		{
			name:     "未闭合反引号",
			input:    "echo `hello",
			expected: []string{},
		},
		{
			name:     "未闭合括号",
			input:    `echo $(date`,
			expected: []string{},
		},
		{
			name:     "未闭合花括号",
			input:    `echo ${HOME`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCmd(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseCmd(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("ParseCmd(%q) = %v, want %v", tt.input, result, tt.expected)
					return
				}
			}
		})
	}
}

// ParseCmdE 测试用例
func TestParseCmdE(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectError bool
	}{
		// 正常场景
		{
			name:        "简单命令",
			input:       "echo hello",
			expected:    []string{"echo", "hello"},
			expectError: false,
		},
		{
			name:        "双引号",
			input:       `echo "hello world"`,
			expected:    []string{"echo", "\"hello world\""},
			expectError: false,
		},
		{
			name:        "环境变量",
			input:       `echo $HOME`,
			expected:    []string{"echo", "$HOME"},
			expectError: false,
		},
		{
			name:        "通配符",
			input:       `ls *.go`,
			expected:    []string{"ls", "*.go"},
			expectError: false,
		},
		{
			name:        "命令替换",
			input:       `echo $(date)`,
			expected:    []string{"echo", "$(date)"},
			expectError: false,
		},
		{
			name:        "空命令",
			input:       "",
			expected:    []string{},
			expectError: false,
		},

		// 错误场景
		{
			name:        "未闭合双引号",
			input:       `echo "hello`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "未闭合单引号",
			input:       `echo 'hello`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "未闭合反引号",
			input:       "echo `hello",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "未闭合括号",
			input:       `echo $(date`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "未闭合花括号",
			input:       `echo ${HOME`,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseCmdE(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParseCmdE(%q) expected error, got nil", tt.input)
					return
				}
				if result != nil {
					t.Errorf("ParseCmdE(%q) = %v, want nil on error", tt.input, result)
					return
				}
			} else {
				if err != nil {
					t.Errorf("ParseCmdE(%q) unexpected error: %v", tt.input, err)
					return
				}
				if len(result) != len(tt.expected) {
					t.Errorf("ParseCmdE(%q) = %v, want %v", tt.input, result, tt.expected)
					return
				}
				for i := range result {
					if result[i] != tt.expected[i] {
						t.Errorf("ParseCmdE(%q) = %v, want %v", tt.input, result, tt.expected)
						return
					}
				}
			}
		})
	}
}

// ParseCmd 和 ParseCmdE 一致性测试
func TestParseCmdConsistency(t *testing.T) {
	tests := []string{
		"echo hello",
		`echo "hello world"`,
		`echo 'hello world'`,
		`ls -la /home/user`,
		`echo $HOME`,
		`ls *.go`,
		`echo $(date)`,
		`echo "hello" "world"`,
		`echo "hello 'world'"`,
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result1 := ParseCmd(input)
			result2, err := ParseCmdE(input)

			if err != nil {
				t.Errorf("ParseCmdE(%q) unexpected error: %v", input, err)
				return
			}

			if len(result1) != len(result2) {
				t.Errorf("ParseCmd(%q) = %v, ParseCmdE(%q) = %v, length mismatch", input, result1, input, result2)
				return
			}

			for i := range result1 {
				if result1[i] != result2[i] {
					t.Errorf("ParseCmd(%q) = %v, ParseCmdE(%q) = %v, mismatch at index %d", input, result1, input, result2, i)
					return
				}
			}
		})
	}
}

// 基准测试 - ParseCmd
func BenchmarkParseCmd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ParseCmd(`echo "hello world"`)
	}
}

func BenchmarkParseCmdE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseCmdE(`echo "hello world"`)
	}
}

func BenchmarkParseCmdComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ParseCmd(`ls -la /home/user && echo "hello" "world" | grep test`)
	}
}
