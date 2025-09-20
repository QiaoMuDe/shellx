package shellx

import (
	"reflect"
	"testing"
)

func TestParseCmd(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		// 基本场景
		{
			name:     "简单命令",
			input:    "ls -la",
			expected: []string{"ls", "-la"},
		},
		{
			name:     "单个命令",
			input:    "pwd",
			expected: []string{"pwd"},
		},
		{
			name:     "多个参数",
			input:    "git commit -m message",
			expected: []string{"git", "commit", "-m", "message"},
		},

		// 引号场景
		{
			name:     "双引号参数",
			input:    `echo "hello world"`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "单引号参数",
			input:    `echo 'hello world'`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "混合引号",
			input:    `grep 'test' "file name"`,
			expected: []string{"grep", "test", "file name"},
		},
		{
			name:     "引号内包含空格",
			input:    `echo "hello   world   test"`,
			expected: []string{"echo", "hello   world   test"},
		},
		{
			name:     "引号内包含另一种引号",
			input:    `echo "He said 'hello'"`,
			expected: []string{"echo", "He said 'hello'"},
		},
		{
			name:     "引号内包含另一种引号2",
			input:    `echo 'He said "hello"'`,
			expected: []string{"echo", `He said "hello"`},
		},
		{
			name:     "反引号参数",
			input:    "echo `date`",
			expected: []string{"echo", "date"},
		},
		{
			name:     "反引号内包含空格",
			input:    "echo `ls -la`",
			expected: []string{"echo", "ls -la"},
		},
		{
			name:     "反引号内包含其他引号",
			input:    `echo ` + "`" + `grep 'test' "file"` + "`",
			expected: []string{"echo", `grep 'test' "file"`},
		},

		// 空白字符场景
		{
			name:     "多个空格",
			input:    "ls    -la    -h",
			expected: []string{"ls", "-la", "-h"},
		},
		{
			name:     "前后空格",
			input:    "  ls -la  ",
			expected: []string{"ls", "-la"},
		},
		{
			name:     "制表符分隔",
			input:    "ls\t-la\t-h",
			expected: []string{"ls", "-la", "-h"},
		},

		// 复杂场景
		{
			name:     "文件路径带空格",
			input:    `cp "source file.txt" "dest file.txt"`,
			expected: []string{"cp", "source file.txt", "dest file.txt"},
		},
		{
			name:     "长命令",
			input:    `find /path -name "*.go" -type f -exec grep "pattern" {} \;`,
			expected: []string{"find", "/path", "-name", "*.go", "-type", "f", "-exec", "grep", "pattern", "{}", `\;`},
		},
		{
			name:     "包含特殊字符",
			input:    `echo "Hello@#$%^&*()_+"`,
			expected: []string{"echo", "Hello@#$%^&*()_+"},
		},

		// 边界场景
		{
			name:     "空字符串",
			input:    "",
			expected: []string{},
		},
		{
			name:     "只有空格",
			input:    "   ",
			expected: []string{},
		},
		{
			name:     "只有制表符",
			input:    "\t\t\t",
			expected: []string{},
		},
		{
			name:     "空引号",
			input:    `echo ""`,
			expected: []string{"echo", ""},
		},
		{
			name:     "空单引号",
			input:    `echo ''`,
			expected: []string{"echo", ""},
		},
		{
			name:     "空反引号",
			input:    "echo ``",
			expected: []string{"echo", ""},
		},

		// 错误场景 - 未闭合引号
		{
			name:     "未闭合双引号",
			input:    `echo "hello world`,
			expected: []string{},
		},
		{
			name:     "未闭合单引号",
			input:    `echo 'hello world`,
			expected: []string{},
		},
		{
			name:     "未闭合反引号",
			input:    "echo `hello world",
			expected: []string{},
		},
		{
			name:     "引号类型不匹配",
			input:    `echo "hello'`,
			expected: []string{},
		},
		{
			name:     "多个未闭合引号",
			input:    `echo "hello 'world`,
			expected: []string{},
		},

		// 特殊引号场景
		{
			name:     "连续引号",
			input:    `echo ""hello""`,
			expected: []string{"echo", "hello"},
		},
		{
			name:     "引号紧挨着文字",
			input:    `echo hello"world"test`,
			expected: []string{"echo", "helloworldtest"},
		},
		{
			name:     "多段引号",
			input:    `echo "hello" "world"`,
			expected: []string{"echo", "hello", "world"},
		},
		{
			name:     "三种引号混合",
			input:    `echo "double" 'single' ` + "`backtick`",
			expected: []string{"echo", "double", "single", "backtick"},
		},

		// 实际使用场景
		{
			name:     "Git命令",
			input:    `git commit -m "Initial commit"`,
			expected: []string{"git", "commit", "-m", "Initial commit"},
		},
		{
			name:     "Docker命令",
			input:    `docker run -it --name "my container" ubuntu:latest`,
			expected: []string{"docker", "run", "-it", "--name", "my container", "ubuntu:latest"},
		},
		{
			name:     "SSH命令",
			input:    `ssh user@host "ls -la /home/user"`,
			expected: []string{"ssh", "user@host", "ls -la /home/user"},
		},
		{
			name:     "Curl命令",
			input:    `curl -H "Content-Type: application/json" -d '{"key":"value"}' http://api.example.com`,
			expected: []string{"curl", "-H", "Content-Type: application/json", "-d", `{"key":"value"}`, "http://api.example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCmd(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseCmd(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// 基准测试
func BenchmarkParseCmd(b *testing.B) {
	testCases := []string{
		"ls -la",
		`echo "hello world"`,
		`git commit -m "Initial commit"`,
		`find /path -name "*.go" -type f -exec grep "pattern" {} \;`,
		`docker run -it --name "my container" ubuntu:latest`,
	}

	for _, tc := range testCases {
		b.Run(tc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseCmd(tc)
			}
		})
	}
}

// 测试大量数据的性能
func BenchmarkParseCmdLarge(b *testing.B) {
	// 构造一个较长的命令
	longCmd := `find /very/long/path/to/search -name "*.go" -o -name "*.js" -o -name "*.py" -type f -exec grep -l "very long pattern to search for in files" {} \; | head -100`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseCmd(longCmd)
	}
}

// 模糊测试（如果Go版本支持）
func FuzzParseCmd(f *testing.F) {
	// 添加种子语料
	seeds := []string{
		"ls -la",
		`echo "hello"`,
		`echo 'world'`,
		"",
		"   ",
		`"unclosed`,
		`'unclosed`,
		`echo "hello 'world'"`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// 确保函数不会panic
		result := ParseCmd(input)

		// 基本不变性检查：确保函数不会panic
		_ = result

		// 如果输入只包含空白字符，结果应该为空
		hasNonSpace := false
		for _, r := range input {
			if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
				hasNonSpace = true
				break
			}
		}
		if !hasNonSpace && len(result) != 0 {
			t.Errorf("Expected empty result for whitespace-only input %q, got %v", input, result)
		}
	})
}
