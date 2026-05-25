# 主包 `shellx` API 参考

## 类型

### ShellType

```go
const (
    ShellSh         ShellType = iota // sh shell
    ShellBash                        // bash shell
    ShellPwsh                        // pwsh (PowerShell Core)
    ShellPowerShell                  // Windows PowerShell
    ShellCmd                         // Windows cmd
    ShellNone                        // 直接执行，不使用 shell
    ShellDef1                        // 自动（Windows: cmd, 其他: sh）
    ShellDef2                        // 自动（Windows: powershell, 其他: sh）
)
```

### Command

集配置、构建、执行于一体的命令对象。每个 Command 只能执行一次。

## 构造函数

```go
func NewCmd(name string, args ...string) *Command
func NewCmds(cmdArgs []string) *Command
func NewCmdStr(cmdStr string) *Command
```

- `NewCmd`: 可变参数，`NewCmd("ls", "-la")`
- `NewCmds`: 切片参数，`NewCmds([]string{"git", "status"})`
- `NewCmdStr`: 字符串自动分词，`NewCmdStr(`echo "hello world"`)`
- 默认 shell 为 `ShellDef1`（Windows: cmd, 其他: sh）
- 默认继承父进程环境变量

## 配置方法

```go
func (c *Command) WithWorkDir(dir string) *Command   // 工作目录（验证存在，不存在 panic）
func (c *Command) WithEnv(key, value string) *Command // 环境变量（验证 key 非空）
func (c *Command) WithEnvs(envs []string) *Command    // 批量环境变量（"key=value" 格式）
func (c *Command) WithTimeout(d time.Duration) *Command // 超时（d>0 才生效）
func (c *Command) WithContext(ctx context.Context) *Command // 上下文（覆盖超时）
func (c *Command) WithStdin(r io.Reader) *Command     // 标准输入
func (c *Command) WithStdout(w io.Writer) *Command    // 标准输出
func (c *Command) WithStderr(w io.Writer) *Command    // 标准错误
func (c *Command) WithShell(shell ShellType) *Command // Shell 类型
```

## 执行方法

```go
func (c *Command) Exec() error                             // 同步执行
func (c *Command) ExecOutput() ([]byte, error)             // 同步 + 合并输出
func (c *Command) ExecStdout() ([]byte, error)             // 同步 + 标准输出
func (c *Command) ExecAsync() error                        // 异步启动
func (c *Command) Wait() error                             // 等待完成
func (c *Command) WaitWithCode() (int, error)              // 等待 + 退出码
```

## 进程控制

```go
func (c *Command) Kill() error                    // 杀死进程
func (c *Command) Signal(sig os.Signal) error     // 发送信号
func (c *Command) IsRunning() bool                // 检查运行状态
func (c *Command) GetPID() int                    // 获取 PID
func (c *Command) Cmd() *exec.Cmd                 // 获取底层对象
```

## 属性获取

```go
func (c *Command) ShellType() ShellType
func (c *Command) Raw() string
func (c *Command) Name() string
func (c *Command) Args() []string         // 返回副本
func (c *Command) CmdStr() string
func (c *Command) WorkDir() string
func (c *Command) Env() []string          // 返回副本
func (c *Command) Timeout() time.Duration
func (c *Command) IsExecuted() bool
```

## 便捷函数

```go
func Exec(name string, args ...string) error
func ExecStr(cmdStr string) error
func ExecOut(name string, args ...string) ([]byte, error)
func ExecOutStr(cmdStr string) ([]byte, error)
func ExecT(d time.Duration, name string, args ...string) error
func ExecStrT(d time.Duration, cmdStr string) error
func ExecOutT(d time.Duration, name string, args ...string) ([]byte, error)
func ExecOutStrT(d time.Duration, cmdStr string) ([]byte, error)
func ExecCode(name string, args ...string) (int, error)
func ExecCodeStr(cmdStr string) (int, error)
```

## 工具函数

```go
func Split(cmdStr string) []string                    // 拆分命令字符串
func SplitE(cmdStr string) ([]string, error)          // 拆分 + 错误检测
func FindCmd(name string) (string, error)             // 查找命令（增强版）
func FindCommandPath(name string) string              // 查找命令（便捷版，找不到返回 ""）
```

- `FindCmd` 处理 ErrDot（Go 1.19+ 安全限制）+ 返回绝对路径
- `FindCommandPath` 找不到返回空字符串

## 错误类型

```go
var ErrAlreadyExecuted = errors.New("command has already been executed")
var ErrNotStarted      = errors.New("command has not been started")
var ErrNoProcess       = errors.New("no process to operate")

type UnclosedQuoteError struct {
    QuoteType rune // 未闭合的引号类型
}
```

错误消息模式：
- 超时: `"command execution timeout: <cmd> exceeded the <time> time limit"`
- 取消: `"command execution canceled: <cmd> was interrupted"`
- 未找到: `"command not found: <cmd> is not a valid command or executable"`
- 退出码: `"command execution failed: <cmd> exited with code <N>"`
