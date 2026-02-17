package shellx

import (
	"testing"
)

// TestUnicodeSpaceHandling 测试Unicode空白字符处理
func TestUnicodeSpaceHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "全角空格",
			input:    "echo\u3000hello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "不换行空格",
			input:    "echo\u00A0hello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "各种宽度空格",
			input:    "echo\u2002hello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "混合空白",
			input:    "echo \t\u00A0\u3000hello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "Windows行结束符",
			input:    "echo \"line1\r\nline2\"",
			expected: []string{"echo", "line1\r\nline2"},
		},
		{
			name:     "零宽空格",
			input:    "echo\u200Bhello",
			expected: []string{"echo\u200Bhello"}, // 零宽空格不会被识别为分隔符
		},
		{
			name:     "换页符",
			input:    "echo\fhello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "垂直制表符",
			input:    "echo\vhello",
			expected: []string{"echo", "hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCmd(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseCmd(%q) = %v, expected %v", tt.input, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("ParseCmd(%q) = %v, expected %v", tt.input, result, tt.expected)
					break
				}
			}
		})
	}
}
