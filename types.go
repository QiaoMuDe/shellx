package shellx

import (
	"runtime"
	"time"
)

// ShellType 定义shell类型
type ShellType int

const (
	ShellSh         ShellType = iota // sh shell
	ShellBash                        // bash shell
	ShellPwsh                        // pwsh (PowerShell Core)
	ShellPowerShell                  // powershell (Windows PowerShell)
	ShellCmd                         // cmd (Windows Command Prompt)
	ShellNone                        // 无shell, 直接原生的执行命令
	ShellDefault                     // 默认shell, 根据操作系统自动选择(Windows系统默认为cmd, 其他系统默认为sh)
)

// String 返回shell类型的字符串表示
func (s ShellType) String() string {
	switch s {
	case ShellSh:
		return "sh"

	case ShellBash:
		return "bash"

	case ShellPwsh:
		return "pwsh"

	case ShellPowerShell:
		return "powershell"

	case ShellCmd:
		return "cmd"

	case ShellNone:
		return "none"

	case ShellDefault:
		if runtime.GOOS == "windows" {
			return "cmd"
		}
		return "sh"

	default:
		return "unknown"
	}
}

// shellFlags 用于返回shell类型的标志
func (s ShellType) shellFlags() string {
	switch s {
	case ShellSh:
		return "-c"

	case ShellBash:
		return "-c"

	case ShellPwsh:
		return "-Command"

	case ShellPowerShell:
		return "-Command"

	case ShellCmd:
		return "/c"

	case ShellNone:
		return ""

	case ShellDefault:
		if runtime.GOOS == "windows" {
			return "/c"
		}
		return "-c"

	default:
		return ""
	}
}

// Result 表示命令执行的结果
type Result struct {
	// 基本执行信息
	exitCode int  // 退出码：0=成功，非0=失败
	success  bool // 是否执行成功

	// 输出信息
	output []byte // 命令输出内容(合并标准输出和标准错误后的内容)

	// 时间信息
	startTime time.Time     // 开始执行时间
	endTime   time.Time     // 结束执行时间
	duration  time.Duration // 执行耗时

	// 错误信息
	err error // 执行过程中的错误
}

// 提供公共访问方法
func (r *Result) Code() int               { return r.exitCode }
func (r *Result) Success() bool           { return r.success }
func (r *Result) Output() []byte          { return r.output }
func (r *Result) Start() time.Time        { return r.startTime }
func (r *Result) End() time.Time          { return r.endTime }
func (r *Result) Duration() time.Duration { return r.duration }
func (r *Result) Error() error            { return r.err }
