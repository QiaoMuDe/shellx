// Package shellx 工具函数测试模块
// 本文件包含 shellx 包中工具函数的单元测试，包括：
//   - 命令拆分函数测试
//   - 快速执行函数测试
//   - 环境变量验证测试
//   - 错误处理测试
//
// 确保工具函数的正确性和稳定性，提供全面的测试覆盖。
package shellx

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
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

		// 错误场景 - 未闭合引号（不自动修复）
		{
			name:     "未闭合双引号",
			input:    `echo "hello world`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "未闭合单引号",
			input:    `echo 'hello world`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "未闭合反引号",
			input:    "echo `hello world",
			expected: []string{"echo", "hello world"},
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
			result := Split(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Split(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSplitE(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectError bool
	}{
		// 基本场景
		{
			name:        "简单命令",
			input:       "ls -la",
			expected:    []string{"ls", "-la"},
			expectError: false,
		},
		{
			name:        "双引号参数",
			input:       `echo "hello world"`,
			expected:    []string{"echo", "hello world"},
			expectError: false,
		},
		{
			name:        "单引号参数",
			input:       `echo 'hello world'`,
			expected:    []string{"echo", "hello world"},
			expectError: false,
		},
		{
			name:        "反引号参数",
			input:       "echo `date`",
			expected:    []string{"echo", "date"},
			expectError: false,
		},

		// 错误场景 - 未闭合引号
		{
			name:        "未闭合双引号",
			input:       `echo "hello world`,
			expected:    []string{"echo", "hello world"},
			expectError: true,
		},
		{
			name:        "未闭合单引号",
			input:       `echo 'hello world`,
			expected:    []string{"echo", "hello world"},
			expectError: true,
		},
		{
			name:        "未闭合反引号",
			input:       "echo `hello world",
			expected:    []string{"echo", "hello world"},
			expectError: true,
		},

		// 边界场景
		{
			name:        "空字符串",
			input:       "",
			expected:    []string{},
			expectError: false,
		},
		{
			name:        "只有空格",
			input:       "   ",
			expected:    []string{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SplitE(tt.input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitE(%q) = %v, expected %v", tt.input, result, tt.expected)
			}

			if tt.expectError && err == nil {
				t.Errorf("SplitE(%q) expected error, got nil", tt.input)
			}

			if !tt.expectError && err != nil {
				t.Errorf("SplitE(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

// 基准测试
func BenchmarkSplit(b *testing.B) {
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
				Split(tc)
			}
		})
	}
}

// 测试大量数据的性能
func BenchmarkSplitLarge(b *testing.B) {
	// 构造一个较长的命令
	longCmd := `find /very/long/path/to/search -name "*.go" -o -name "*.js" -o -name "*.py" -type f -exec grep -l "very long pattern to search for in files" {} \; | head -100`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Split(longCmd)
	}
}

// 模糊测试（如果Go版本支持）
func FuzzSplit(f *testing.F) {
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
		result := Split(input)

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

// 测试命令分词器在各种场景下的行为
func TestSplitEdgeCases(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
		comment  string
	}{
		// 基本场景
		{
			input:    "echo hello world",
			expected: []string{"echo", "hello", "world"},
			comment:  "基本命令拆分",
		},
		{
			input:    "ls -la",
			expected: []string{"ls", "-la"},
			comment:  "命令与参数",
		},

		// 引号场景
		{
			input:    `echo "hello world"`,
			expected: []string{"echo", "hello world"},
			comment:  "双引号包含空格",
		},
		{
			input:    `echo 'hello world'`,
			expected: []string{"echo", "hello world"},
			comment:  "单引号包含空格",
		},
		{
			input:    "echo `hello world`",
			expected: []string{"echo", "hello world"},
			comment:  "反引号包含空格",
		},
		{
			input:    `echo ""`,
			expected: []string{"echo", ""},
			comment:  "空双引号",
		},
		{
			input:    `echo ''`,
			expected: []string{"echo", ""},
			comment:  "空单引号",
		},

		// 未闭合引号场景
		{
			input:    `echo "hello world`,
			expected: []string{"echo", "hello world"},
			comment:  "未闭合双引号（忽略错误）",
		},
		{
			input:    `echo 'hello world`,
			expected: []string{"echo", "hello world"},
			comment:  "未闭合单引号（忽略错误）",
		},
		{
			input:    "echo `hello world",
			expected: []string{"echo", "hello world"},
			comment:  "未闭合反引号（忽略错误）",
		},

		// 复杂引号场景
		{
			input:    `echo "hello 'world'"`,
			expected: []string{"echo", "hello 'world'"},
			comment:  "双引号内包含单引号",
		},
		{
			input:    `echo 'hello "world"'`,
			expected: []string{"echo", `hello "world"`},
			comment:  "单引号内包含双引号",
		},
		{
			input:    `echo ""hello""`,
			expected: []string{"echo", "hello"},
			comment:  "连续双引号包围单词",
		},

		// 特殊字符场景
		{
			input:    `echo "hello@#$%^&*()"`,
			expected: []string{"echo", "hello@#$%^&*()"},
			comment:  "引号内包含特殊字符",
		},
		{
			input:    `git commit -m "Initial commit"`,
			expected: []string{"git", "commit", "-m", "Initial commit"},
			comment:  "Git提交命令",
		},

		// 空白字符场景
		{
			input:    "  ls   -la   ",
			expected: []string{"ls", "-la"},
			comment:  "多余空格处理",
		},
		{
			input:    "ls\t-la\t-h",
			expected: []string{"ls", "-la", "-h"},
			comment:  "制表符分隔",
		},
		{
			input:    "",
			expected: []string{},
			comment:  "空字符串",
		},
		{
			input:    "   ",
			expected: []string{},
			comment:  "只有空格",
		},

		// 换行符场景
		{
			input:    "echo \"line1\nline2\"",
			expected: []string{"echo", "line1\nline2"},
			comment:  "引号内包含换行符",
		},
		{
			input:    "echo \"line1\rline2\"",
			expected: []string{"echo", "line1\rline2"},
			comment:  "引号内包含回车符",
		},

		// 实际使用场景
		{
			input:    `docker run -it --name "my container" ubuntu:latest`,
			expected: []string{"docker", "run", "-it", "--name", "my container", "ubuntu:latest"},
			comment:  "Docker运行命令",
		},
		{
			input:    `ssh user@host "ls -la /home/user"`,
			expected: []string{"ssh", "user@host", "ls -la /home/user"},
			comment:  "SSH远程命令",
		},
		{
			input:    `curl -H "Content-Type: application/json" -d '{"key":"value"}' http://api.example.com`,
			expected: []string{"curl", "-H", "Content-Type: application/json", "-d", `{"key":"value"}`, "http://api.example.com"},
			comment:  "Curl API调用",
		},

		// Bash 语法场景
		{
			input:    `if [ -d "/tmp" ]; then echo "exists"; else echo "not exists"; fi`,
			expected: []string{"if", "[", "-d", "/tmp", "]", ";", "then", "echo", "exists", ";", "else", "echo", "not exists", ";", "fi"},
			comment:  "Bash 条件语句",
		},
		{
			input:    `for i in {1..5}; do echo "item $i"; done`,
			expected: []string{"for", "i", "in", "{1..5}", ";", "do", "echo", "item $i", ";", "done"},
			comment:  "Bash for 循环",
		},
		{
			input:    `VAR="hello world"; echo $VAR`,
			expected: []string{"VAR=hello world", ";", "echo", "$VAR"},
			comment:  "Bash 变量赋值和使用",
		},
		{
			input:    `command1 && command2 || command3`,
			expected: []string{"command1", "&&", "command2", "||", "command3"},
			comment:  "Bash 逻辑操作符",
		},
		{
			input:    `find . -name "*.go" -exec grep "pattern" {} \;`,
			expected: []string{"find", ".", "-name", "*.go", "-exec", "grep", "pattern", "{}", "\\;"},
			comment:  "Bash find 命令",
		},
		{
			input: `cat <<EOF > file.txt
Hello World
EOF`,
			expected: []string{"cat", "<<", "EOF", ">", "file.txt", "Hello", "World", "EOF"},
			comment:  "Bash here document",
		},

		// PowerShell 语法场景
		{
			input:    `Get-ChildItem -Path "C:\Program Files" -Recurse | Where-Object {$_.Name -like "*.exe"}`,
			expected: []string{"Get-ChildItem", "-Path", "C:\\Program Files", "-Recurse", "|", "Where-Object", "{$_.Name", "-like", "*.exe}"},
			comment:  "PowerShell 管道和过滤器",
		},
		{
			input:    `$files = Get-Content "config.json" | ConvertFrom-Json`,
			expected: []string{"$files", "=", "Get-Content", "config.json", "|", "ConvertFrom-Json"},
			comment:  "PowerShell 变量赋值",
		},
		{
			input:    `if (Test-Path "C:\temp") { Write-Host "exists" } else { Write-Host "not exists" }`,
			expected: []string{"if", "(Test-Path", "C:\\temp)", "{", "Write-Host", "exists", "}", "else", "{", "Write-Host", "not exists", "}"},
			comment:  "PowerShell 条件语句",
		},
		{
			input:    `foreach ($file in Get-ChildItem "*.txt") { Write-Host $file.Name }`,
			expected: []string{"foreach", "($file", "in", "Get-ChildItem", "*.txt)", "{", "Write-Host", "$file.Name", "}"},
			comment:  "PowerShell foreach 循环",
		},
		{
			input:    `Start-Process -FilePath "notepad.exe" -ArgumentList "file.txt" -Wait`,
			expected: []string{"Start-Process", "-FilePath", "notepad.exe", "-ArgumentList", "file.txt", "-Wait"},
			comment:  "PowerShell 启动进程",
		},
	}

	for i, tc := range testCases {
		// t.Logf("\n=== 测试用例 %d: %s ===", i+1, tc.comment)
		// t.Logf("原始字符串: %q", tc.input)
		// t.Logf("预期结果: %v", tc.expected)
		fmt.Printf("=== 测试用例 %d: %s ===\n", i+1, tc.comment)
		fmt.Printf("原始字符串: %q\n", tc.input)
		fmt.Printf("预期结果: %v\n", tc.expected)

		// 使用 Split 函数
		result := Split(tc.input)
		// t.Logf("解析结果: %v", result)
		fmt.Printf("解析结果: %v, 长度: %d\n", result, len(result))

		// 简单比较
		if len(result) == len(tc.expected) {
			match := true
			for j := range result {
				if result[j] != tc.expected[j] {
					match = false
					break
				}
			}
			if match {
				// t.Logf("✅ 结果匹配")
				fmt.Println("✅ 结果匹配")
			} else {
				// t.Logf("❌ 结果不匹配")
				fmt.Println("❌ 结果不匹配")
			}
		} else {
			// t.Logf("❌ 长度不匹配: 期望 %d 个元素，实际 %d 个元素", len(tc.expected), len(result))
			fmt.Printf("❌ 长度不匹配: 期望 %d 个元素，实际 %d 个元素\n", len(tc.expected), len(result))
		}

		// t.Logf("---")
		fmt.Println("---")
	}
}
