# 子包 `shx` API 参考

## 类型

### Shx

```go
type Shx struct { /* 未导出字段 */ }
```

- 配置方法不是并发安全的
- 每个 Shx 只能执行一次
- 不支持进程控制（无 PID/Kill/Signal）
- 同步执行，如需异步用 goroutine 包装

### ExitStatus

```go
type ExitStatus struct {
    Code uint8
}
func (e ExitStatus) Error() string
func IsExitStatus(err error) (uint8, bool)
```

## 构造函数

```go
func New(cmdStr string) *Shx                         // 从字符串创建
func NewArgs(cmd string, args ...string) *Shx        // 命令+参数
func NewCmds(cmds []string) *Shx                     // 从切片创建
func NewWithParser(cmdStr string, p *syntax.Parser) *Shx // 自定义解析器
```

- `New` 是最常用的方式
- `NewWithParser` 可注入自定义 `mvdan.cc/sh/v3` 解析器

## 配置方法

```go
func (s *Shx) WithDir(dir string) *Shx
func (s *Shx) WithEnv(key, value string) *Shx
func (s *Shx) WithEnvMap(envs map[string]string) *Shx
func (s *Shx) WithEnvs(envs []string) *Shx
func (s *Shx) WithStdin(r io.Reader) *Shx
func (s *Shx) WithStdout(w io.Writer) *Shx
func (s *Shx) WithStderr(w io.Writer) *Shx
func (s *Shx) WithTimeout(d time.Duration) *Shx
func (s *Shx) WithContext(ctx context.Context) *Shx
```

## 执行方法

```go
func (s *Shx) Exec() error
func (s *Shx) ExecOutput() ([]byte, error)
func (s *Shx) ExecContext(ctx context.Context) error
func (s *Shx) ExecContextOutput(ctx context.Context) ([]byte, error)
```

## 属性获取

```go
func (s *Shx) Raw() string
func (s *Shx) Dir() string
func (s *Shx) Env() expand.Environ
func (s *Shx) Timeout() time.Duration
func (s *Shx) Context() context.Context
func (s *Shx) IsExecuted() bool
```

## 便捷函数

```go
func Run(cmd string) error
func RunToTerminal(cmd string) error
func Out(cmd string) ([]byte, error)
func RunWith(cmd string, d time.Duration) error
func OutWith(cmd string, d time.Duration) ([]byte, error)
func RunWithIO(cmd string, stdin io.Reader, stdout, stderr io.Writer) error
func OutWithIO(cmd string, stdin io.Reader, stdout, stderr io.Writer) ([]byte, error)
func RunCtx(ctx context.Context, cmd string) error
func OutCtx(ctx context.Context, cmd string) ([]byte, error)
```

## 错误类型

```go
var ErrAlreadyExecuted = errors.New("command has already been executed")
var ErrNilContext      = errors.New("context cannot be nil")
var ErrNilReader       = errors.New("reader cannot be nil")
var ErrNilWriter       = errors.New("writer cannot be nil")
```

## 设计要点

- 纯 Go 实现，基于 `mvdan.cc/sh/v3` 解析和执行
- 执行流程：Parse → AST → interp.Runner → Run
- 跨平台一致性：Windows/Linux/macOS 行为一致
- 上下文优先级：`WithContext` > `WithTimeout` > `context.Background()`
- 不支持异步，mvdan/sh 本身是同步的
