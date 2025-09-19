package shellx

// ShellType 定义shell类型
type ShellType int

const (
	ShellSh         ShellType = iota // sh shell
	ShellBash                        // bash shell
	ShellPwsh                        // pwsh (PowerShell Core)
	ShellPowerShell                  // powershell (Windows PowerShell)
	ShellCmd                         // cmd (Windows Command Prompt)
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
	default:
		return "unknown"
	}
}

// ExecuteOptions 执行选项配置
type ExecuteOptions struct {
	Shell   ShellType // 指定shell类型
	Capture bool      // 是否捕获输出到Result
}
