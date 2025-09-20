package shellx

// ShellType 定义shell类型
type ShellType int

const (
	ShellSh         ShellType = iota // sh shell
	ShellBash                        // bash shell
	ShellPwsh                        // pwsh (PowerShell Core)
	ShellPowerShell                  // powershell (Windows PowerShell)
	ShellCmd                         // cmd (Windows Command Prompt)
	ShellNone                        // 无shell, 直接原生的执行命令
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

	default:
		return ""
	}
}
