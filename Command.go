package shellx

import "os/exec"

// Command 命令对象
type Command struct {
	cmd *exec.Cmd
}

func (c *Command) Cmd() *exec.Cmd {
	return c.cmd
}

func (c *Command) Exec() error {
	return nil
}

func (c *Command) ExecOutput() ([]byte, error) {
	return nil, nil
}

func (c *Command) ExecAsync() error {
	return nil
}

func (c *Command) Wait() error {
	return nil
}
