package shx

import (
	"context"
	"strings"
	"testing"

	"mvdan.cc/sh/v3/interp"
)

func TestIsExitStatus(t *testing.T) {
	// 测试退出码错误
	err := ExitStatus{Code: 42}

	code, ok := IsExitStatus(err)
	if !ok {
		t.Fatal("expected true")
	}

	if code != 42 {
		t.Fatalf("expected code 42, got %d", code)
	}
}

func TestIsExitStatusNotExitStatus(t *testing.T) {
	// 测试非退出码错误
	err := &testError{msg: "not an exit status"}

	_, ok := IsExitStatus(err)
	if ok {
		t.Fatal("expected false")
	}
}

func TestExitStatusError(t *testing.T) {
	// 测试 ExitStatus.Error()
	err := ExitStatus{Code: 127}

	expected := "exit status 127"
	if err.Error() != expected {
		t.Fatalf("expected %s, got %s", expected, err.Error())
	}
}

func TestIsExitStatusWithInterpExitStatus(t *testing.T) {
	// 测试原生interp.ExitStatus
	// 注意：interp.ExitStatus是uint8类型，不是接口

	// 直接使用uint8类型作为interp.ExitStatus
	var err error = interp.ExitStatus(42)

	code, ok := IsExitStatus(err)
	if !ok {
		t.Fatal("expected true")
	}

	if code != 42 {
		t.Fatalf("expected code 42, got %d", code)
	}
}

func TestHandleErrorNil(t *testing.T) {
	// 测试 nil 错误
	err := handleError(nil, "test cmd", 0)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestHandleErrorContextCanceled(t *testing.T) {
	// 测试上下文取消错误
	err := handleError(context.Canceled, "test cmd", 0)

	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "command canceled") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleErrorDeadlineExceeded(t *testing.T) {
	// 测试超时错误（无超时时间）
	err := handleError(context.DeadlineExceeded, "test cmd", 0)

	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "command timed out") {
		t.Fatalf("unexpected error: %v", err)
	}

	// 不应该包含超时时间
	if strings.Contains(err.Error(), "after") {
		t.Fatalf("should not contain timeout duration: %v", err)
	}
}

func TestHandleErrorDeadlineExceededWithTimeout(t *testing.T) {
	// 测试超时错误（有超时时间）
	err := handleError(context.DeadlineExceeded, "test cmd", 5)

	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "command timed out after 5ns") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleErrorOther(t *testing.T) {
	// 测试其他错误
	testErr := &testError{msg: "test error"}
	err := handleError(testErr, "test cmd", 0)

	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "command failed") {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(err.Error(), "test cmd") {
		t.Fatalf("should contain command string: %v", err)
	}

	if !strings.Contains(err.Error(), "test error") {
		t.Fatalf("should contain original error: %v", err)
	}
}

// 测试用的错误类型
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
