package shellx

import "os/exec"

// Command 命令对象
type Command struct {
	cmd *exec.Cmd
}
