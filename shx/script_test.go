package shx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// createTempScript 创建临时脚本文件
//
// 参数:
//   - t: 测试对象
//   - content: 脚本内容
//
// 返回:
//   - string: 脚本文件路径
//   - func(): 清理函数
func createTempScript(t *testing.T, content string) (string, func()) {
	t.Helper()

	dir, err := os.MkdirTemp("", "shx-script-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	scriptPath := filepath.Join(dir, "test.sh")
	if err := os.WriteFile(scriptPath, []byte(content), 0755); err != nil {
		_ = os.RemoveAll(dir)
		t.Fatalf("failed to write temp script: %v", err)
	}

	return scriptPath, func() {
		_ = os.RemoveAll(dir)
	}
}

func TestNewScript(t *testing.T) {
	scriptPath, cleanup := createTempScript(t, "echo hello script")
	defer cleanup()

	cmd := NewScript(scriptPath)

	if cmd.scriptFile != scriptPath {
		t.Fatalf("expected scriptFile %q, got %q", scriptPath, cmd.scriptFile)
	}

	if cmd.raw != scriptPath {
		t.Fatalf("expected raw %q, got %q", scriptPath, cmd.raw)
	}

	if cmd.Raw() != scriptPath {
		t.Fatalf("expected Raw() %q, got %q", scriptPath, cmd.Raw())
	}
}

func TestNewScriptPanicOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty script path")
		}
	}()

	NewScript("")
}

func TestRunScript(t *testing.T) {
	scriptPath, cleanup := createTempScript(t, "echo hello from script")
	defer cleanup()

	err := RunScript(scriptPath)
	if err != nil {
		t.Fatalf("RunScript failed: %v", err)
	}
}

func TestOutScript(t *testing.T) {
	scriptPath, cleanup := createTempScript(t, "echo hello output")
	defer cleanup()

	output, err := OutScript(scriptPath)
	if err != nil {
		t.Fatalf("OutScript failed: %v", err)
	}

	if !strings.Contains(string(output), "hello output") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestScriptFileNotFound(t *testing.T) {
	cmd := NewScript("nonexistent_script_12345.sh")

	err := cmd.Exec()
	if err == nil {
		t.Fatal("expected error for nonexistent script file")
	}

	if !strings.Contains(err.Error(), "open script file failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScriptWithTimeout(t *testing.T) {
	scriptPath, cleanup := createTempScript(t, "echo hello")
	defer cleanup()

	cmd := NewScript(scriptPath).WithTimeout(5 * time.Second)

	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec with timeout failed: %v", err)
	}
}

func TestScriptWithTimeoutActual(t *testing.T) {
	scriptPath, cleanup := createTempScript(t, "ping -n 10 127.0.0.1 > nul")
	defer cleanup()

	cmd := NewScript(scriptPath).WithTimeout(50 * time.Millisecond)

	err := cmd.Exec()
	if err == nil {
		t.Fatal("expected timeout error")
	}

	if !strings.Contains(err.Error(), "timed out") && !strings.Contains(err.Error(), "exit status") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScriptWithDir(t *testing.T) {
	// 创建一个脚本，pwd 输出当前目录
	scriptPath, cleanup := createTempScript(t, "pwd")
	defer cleanup()

	tmpDir, err := os.MkdirTemp("", "shx-script-dir-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cmd := NewScript(scriptPath).WithDir(tmpDir)

	var buf strings.Builder
	cmd.WithStdout(&buf).WithStderr(&buf)

	err = cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output != tmpDir {
		t.Fatalf("expected output %q, got %q", tmpDir, output)
	}
}

func TestScriptWithEnv(t *testing.T) {
	scriptPath, cleanup := createTempScript(t, "echo $TEST_VAR_12345")
	defer cleanup()

	cmd := NewScript(scriptPath).WithEnv("TEST_VAR_12345", "env_value")

	var buf strings.Builder
	cmd.WithStdout(&buf).WithStderr(&buf)

	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if !strings.Contains(strings.TrimSpace(buf.String()), "env_value") {
		t.Fatalf("expected env_value in output, got: %s", buf.String())
	}
}

func TestScriptChainedConfig(t *testing.T) {
	scriptPath, cleanup := createTempScript(t, "echo $MY_VAR from $(pwd)")
	defer cleanup()

	tmpDir, err := os.MkdirTemp("", "shx-script-chain-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cmd := NewScript(scriptPath).
		WithDir(tmpDir).
		WithEnv("MY_VAR", "hello").
		WithTimeout(10 * time.Second)

	output, err := cmd.ExecOutput()
	if err != nil {
		t.Fatalf("ExecOutput failed: %v", err)
	}

	if !strings.Contains(string(output), "hello") {
		t.Fatalf("expected 'hello' in output, got: %s", output)
	}
}

func TestScriptDuplicateExecution(t *testing.T) {
	scriptPath, cleanup := createTempScript(t, "echo first")
	defer cleanup()

	cmd := NewScript(scriptPath)

	// 第一次执行
	err := cmd.Exec()
	if err != nil {
		t.Fatalf("first Exec failed: %v", err)
	}

	// 第二次执行应该返回 ErrAlreadyExecuted
	err = cmd.Exec()
	if err == nil {
		t.Fatal("expected error for duplicate execution")
	}

	if !strings.Contains(err.Error(), "already been executed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScriptWithBinBash(t *testing.T) {
	// 测试 Bash 特有语法 ([[ ]])
	scriptPath, cleanup := createTempScript(t, `if [[ "hello" == "hello" ]]; then echo "bash matched"; fi`)
	defer cleanup()

	cmd := NewScript(scriptPath)

	var buf strings.Builder
	cmd.WithStdout(&buf).WithStderr(&buf)

	err := cmd.Exec()
	if err != nil {
		t.Fatalf("Exec with bash syntax failed: %v", err)
	}

	if !strings.Contains(buf.String(), "bash matched") {
		t.Fatalf("expected 'bash matched' in output, got: %s", buf.String())
	}
}
