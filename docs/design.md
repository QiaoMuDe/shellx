# Shell命令执行库设计文档

## 项目概述

本项目旨在设计一个功能完善、易于使用的Go语言shell命令执行库，提供类型安全、灵活配置的命令执行能力。

## 1. 核心数据模型设计

### 1.1 命令结构体 (Cmd)

```go
// Cmd 表示一个待执行的shell命令
type Cmd struct {
    // 基本命令信息 - 支持两种方式
    name string   // 命令名称，如 "ls", "git", "docker" (当使用数组方式时)
    args []string // 命令参数列表 (当使用数组方式时)
    raw  string   // 原始命令字符串 (当使用字符串方式时)
    
    // 执行环境配置
    workDir string            // 工作目录
    env     map[string]string // 环境变量
    
    // 输入输出配置
    stdin  io.Reader // 标准输入
    stdout io.Writer // 标准输出重定向
    stderr io.Writer // 标准错误重定向
    
    // 执行控制
    timeout time.Duration   // 超时时间
    context context.Context // 上下文控制
    
    // 执行选项
    options *ExecuteOptions
}

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
    Shell     ShellType // 指定shell类型
    Capture   bool      // 是否捕获输出到Result

}

// 提供公共访问方法
func (c *Cmd) Name() string { return c.name }
func (c *Cmd) Args() []string { return c.args }
func (c *Cmd) Raw() string { return c.raw }
func (c *Cmd) Dir() string { return c.workDir }
func (c *Cmd) Env() map[string]string { return c.env }
func (c *Cmd) Input() io.Reader { return c.stdin }
func (c *Cmd) Output() io.Writer { return c.stdout }
func (c *Cmd) ErrOutput() io.Writer { return c.stderr }
func (c *Cmd) Timeout() time.Duration { return c.timeout }
func (c *Cmd) Ctx() context.Context { return c.context }
func (c *Cmd) Opts() *ExecuteOptions { return c.options }
```

**设计说明:**
- 支持两种命令创建方式：
  - **数组方式**: `name` + `args` 分离，便于参数处理和验证
  - **字符串方式**: `raw` 原始命令字符串，支持复杂shell语法
- 内部字段使用小写，遵循Go的封装原则
- `workDir` 支持指定执行目录，提高灵活性
- `env` 使用map存储环境变量，支持动态配置
- 支持标准输入输出重定向，满足复杂场景需求
- `context` 集成Go标准库的上下文管理
- `ExecuteOptions` 独立配置结构，便于扩展
- 类型安全的Shell类型定义
- 提供公共访问方法来获取内部字段

### 1.2 执行结果结构体 (Result)

```go
// Result 表示命令执行的结果
type Result struct {
    // 基本执行信息
    cmd        *Cmd          // 执行的命令引用
    exitCode   int           // 退出码
    success    bool          // 是否执行成功
    
    // 输出信息
    stdout     []byte        // 标准输出内容
    stderr     []byte        // 标准错误内容
    output     []byte        // 合并输出
    
    // 时间信息
    startTime  time.Time     // 开始执行时间
    endTime    time.Time     // 结束执行时间
    duration   time.Duration // 执行耗时
    
    // 进程信息
    pid          int                // 进程ID
    processState *os.ProcessState   // 进程状态
    
    // 错误信息
    err        error         // 执行过程中的错误
    
    // 元数据
    metadata   map[string]interface{} // 额外的元数据信息
}

// 提供公共访问方法
func (r *Result) Cmd() *Cmd { return r.cmd }
func (r *Result) Code() int { return r.exitCode }
func (r *Result) Success() bool { return r.success }
func (r *Result) StdOut() []byte { return r.stdout }
func (r *Result) StdErr() []byte { return r.stderr }
func (r *Result) Output() []byte { return r.output }
func (r *Result) Start() time.Time { return r.startTime }
func (r *Result) End() time.Time { return r.endTime }
func (r *Result) Duration() time.Duration { return r.duration }
func (r *Result) PID() int { return r.pid }
func (r *Result) State() *os.ProcessState { return r.processState }
func (r *Result) Error() error { return r.err }
func (r *Result) Meta() map[string]interface{} { return r.metadata }
```

**设计说明:**
- 包含完整的执行信息，便于调试和监控
- 分离标准输出和错误输出，同时提供合并输出
- 详细的时间信息，支持性能分析
- 进程信息便于进程管理
- 内部字段使用小写，遵循Go的封装原则
- 提供公共访问方法来获取内部字段
- `metadata` 字段提供扩展性

### 1.3 执行器接口 (Executor)

```go
// Executor 命令执行器接口
type Executor interface {
    // Exec 执行单个命令
    Exec(cmd *Cmd) (*Result, error)
    
    // ExecAsync 异步执行命令
    ExecAsync(cmd *Cmd) (<-chan *Result, error)
    
    // ExecPipe 执行命令管道 (可变参数方式)
    ExecPipe(commands ...*Cmd) (*Result, error)
    
    // ExecPipes 执行命令管道 (切片方式)
    ExecPipes(commands []*Cmd) (*Result, error)
    
    // Kill 终止正在执行的命令
    Kill(pid int) error
    
    // Running 检查命令是否正在运行
    Running(pid int) bool
}
```

**设计说明:**
- 接口设计遵循单一职责原则
- 支持同步和异步执行模式
- 提供管道执行能力
- 包含进程管理功能

### 1.4 命令构建器 (Builder)

```go
// Builder 命令构建器，提供链式调用
type Builder struct {
    cmd *Cmd
}

// NewCmd 创建新的命令构建器 (数组方式 - 可变参数)
func NewCmd(name string, args ...string) *Builder

// NewCmds 创建新的命令构建器 (数组方式 - 切片参数，第一个元素为命令名)
func NewCmds(cmdArgs []string) *Builder

// NewCmdString 创建新的命令构建器 (字符串方式)
func NewCmdString(cmdStr string) *Builder

// 链式方法
func (b *Builder) WithArgs(args ...string) *Builder
func (b *Builder) WithWorkDir(dir string) *Builder
func (b *Builder) WithEnv(key, value string) *Builder
func (b *Builder) WithTimeout(timeout time.Duration) *Builder
func (b *Builder) WithContext(ctx context.Context) *Builder
func (b *Builder) WithStdin(stdin io.Reader) *Builder
func (b *Builder) WithStdout(stdout io.Writer) *Builder
func (b *Builder) WithStderr(stderr io.Writer) *Builder
func (b *Builder) WithOptions(opts *ExecuteOptions) *Builder
func (b *Builder) Build() *Cmd
```

**设计说明:**
- Builder模式提供友好的API体验
- 支持三种创建方式：
  - `NewCmd()`: 数组方式，使用可变参数，适合简单命令
  - `NewCmds()`: 数组方式，使用切片参数，适合动态构建的命令
  - `NewCmdString()`: 字符串方式，适合复杂shell命令
- 链式调用减少代码冗余
- 每个配置项都有对应的方法
- `Build()` 方法返回最终的Cmd对象

### 1.5 错误类型定义

```go
// ExecutionError 执行错误类型
type ExecutionError struct {
    Cmd       *Cmd
    ExitCode  int
    Stderr    string
    Err       error
    Timestamp time.Time
}

func (e *ExecutionError) Error() string {
    cmdStr := e.Cmd.Name()
    if e.Cmd.Raw() != "" {
        cmdStr = e.Cmd.Raw()
    }
    return fmt.Sprintf("command '%s' failed with exit code %d: %s", 
        cmdStr, e.ExitCode, e.Stderr)
}

// TimeoutError 超时错误
type TimeoutError struct {
    Cmd     *Cmd
    Timeout time.Duration
}

func (e *TimeoutError) Error() string {
    cmdStr := e.Cmd.Name()
    if e.Cmd.Raw() != "" {
        cmdStr = e.Cmd.Raw()
    }
    return fmt.Sprintf("command '%s' timed out after %v", 
        cmdStr, e.Timeout)
}

// ValidationError 验证错误
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error in field '%s': %s", 
        e.Field, e.Message)
}
```

**设计说明:**
- 定义具体的错误类型，便于错误处理
- 每个错误类型包含相关的上下文信息
- 实现Error接口，提供有意义的错误消息
- 错误消息支持两种命令格式的显示

## 2. 使用示例

### 2.1 基本使用

```go
// 方式1: 数组方式创建命令 (可变参数)
cmd := NewCmd("ls", "-la").
    WithWorkDir("/tmp").
    WithTimeout(30 * time.Second).
    Build()

// 方式2: 数组方式创建命令 (切片参数)
cmdArgs := []string{"ls", "-la", "-h"}
cmd2 := NewCmds(cmdArgs).
    WithWorkDir("/tmp").
    WithTimeout(30 * time.Second).
    Build()

// 方式3: 字符串方式创建命令
cmd3 := NewCmdString("ls -la | grep test").
    WithWorkDir("/tmp").
    WithTimeout(30 * time.Second).
    Build()

executor := NewExecutor()
result, err := executor.Exec(cmd)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Exit Code: %d\n", result.Code())
fmt.Printf("Output: %s\n", result.StdOut())
```

### 2.2 环境变量配置

```go
// 数组方式
cmd := NewCmd("env").
    WithEnv("MY_VAR", "my_value").
    WithEnv("PATH", "/usr/local/bin:/usr/bin").
    Build()

// 字符串方式 - 适合复杂的shell命令
cmd2 := NewCmdString("export MY_VAR=my_value && echo $MY_VAR").
    WithEnv("PATH", "/usr/local/bin:/usr/bin").
    Build()

result, err := executor.Exec(cmd)
```

### 2.3 Shell类型配置

```go
// 指定使用bash执行命令
opts := &ExecuteOptions{
    Shell:   ShellBash,
    Capture: true,
}

cmd := NewCmd("echo", "Hello World").
    WithOptions(opts).
    Build()

// 在Windows上使用PowerShell
optsWin := &ExecuteOptions{
    Shell:   ShellPowerShell,
    Capture: true,
}

cmd2 := NewCmdString("Get-Process | Where-Object {$_.Name -eq 'notepad'}").
    WithOptions(optsWin).
    Build()

result, err := executor.Execute(cmd)
```

### 2.4 输入输出重定向

```go
var stdout, stderr bytes.Buffer

// 数组方式
cmd := NewCmd("grep", "error").
    WithStdin(strings.NewReader("some input\nerror line\nother line")).
    WithStdout(&stdout).
    WithStderr(&stderr).
    Build()

// 字符串方式 - 支持shell重定向语法
cmd2 := NewCmdString("grep error < input.txt > output.txt 2> error.txt").
    Build()

result, err := executor.Execute(cmd)
```

### 2.5 管道命令

```go
// 方式1: 可变参数方式 (推荐)
result, err := executor.ExecPipe(
    NewCmd("ps", "aux").Build(),
    NewCmd("grep", "go").Build(),
    NewCmd("wc", "-l").Build(),
)

// 方式2: 切片方式
pipeline := []*Cmd{
    NewCmd("ps", "aux").Build(),
    NewCmd("grep", "go").Build(),
    NewCmd("wc", "-l").Build(),
}
result, err := executor.ExecPipes(pipeline)

// 方式3: 使用字符串方式的单个管道命令
cmd := NewCmdString("ps aux | grep go | wc -l").Build()
result, err := executor.Exec(cmd)
```

### 2.6 异步执行

```go
resultChan, err := executor.ExecAsync(cmd)
if err != nil {
    log.Fatal(err)
}

go func() {
    result := <-resultChan
    fmt.Printf("Command finished with exit code: %d\n", result.Code())
}()
```

### 2.6 上下文控制

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// 数组方式
cmd := NewCmd("long-running-command").
    WithContext(ctx).
    Build()

// 字符串方式 - 支持复杂命令
cmd2 := NewCmdString("sleep 60 && echo done").
    WithContext(ctx).
    Build()

result, err := executor.Execute(cmd)
```

## 3. 设计特点

### 3.1 类型安全
- 使用强类型定义，避免运行时错误
- 明确的接口定义，便于测试和mock

### 3.2 易用性
- Builder模式提供友好的API
- 链式调用减少代码冗余
- 合理的默认值设置
- 支持两种命令创建方式，满足不同场景需求

### 3.3 灵活性
- 支持多种执行选项和环境配置
- 可扩展的元数据字段
- 插件化的执行器设计
- 通过重定向控制输出行为，更加灵活


### 3.4 错误处理
- 详细的错误类型和信息
- 结构化的错误数据
- 便于调试的错误消息

### 3.5 异步支持
- 支持同步和异步执行
- 基于channel的异步通信
- 上下文控制支持取消操作

### 3.6 管道支持
- 支持命令管道操作
- 自动处理管道间的数据传递
- 错误传播机制
- 支持shell原生管道语法

### 3.7 跨平台兼容
- 考虑不同操作系统的差异
- 支持不同的shell类型
- 统一的API接口

### 3.8 封装性
- 内部字段使用小写，遵循Go的封装原则
- 提供公共访问方法
- 保护内部状态不被意外修改

## 4. 扩展性考虑

### 4.1 插件系统
- 支持自定义执行器实现
- 中间件模式支持功能扩展
- 钩子函数支持自定义逻辑

### 4.2 监控和日志
- 内置执行统计功能
- 可配置的日志输出
- 性能监控指标

### 4.3 安全性
- 命令参数验证
- 环境变量过滤
- 权限检查机制

### 4.4 命令解析
- 支持复杂shell语法解析
- 参数转义和引号处理
- 环境变量展开

## 5. 下一步计划

1. **实现核心数据结构**: 完成Cmd、Result等基础类型
2. **实现Builder模式**: 提供友好的命令构建API
3. **实现命令解析器**: 支持字符串到数组的转换
4. **实现基础执行器**: 支持基本的命令执行功能
5. **添加错误处理**: 完善错误类型和处理机制
6. **实现异步执行**: 支持非阻塞的命令执行
7. **添加管道支持**: 实现命令管道功能
8. **跨平台适配**: 处理不同操作系统的差异
9. **编写测试用例**: 确保代码质量和稳定性
10. **性能优化**: 优化执行效率和资源使用
11. **文档完善**: 提供详细的使用文档和示例

## 6. 技术栈

- **语言**: Go 1.19+
- **依赖**: 尽量使用标准库，减少外部依赖
- **测试**: 使用Go标准测试框架
- **文档**: 使用godoc格式注释

## 7. 命令解析策略

### 7.1 数组方式
- 直接使用提供的name和args
- 适合简单、明确的命令
- 参数已经分离，无需解析

### 7.2 字符串方式
- 需要解析原始命令字符串
- 支持shell语法：管道、重定向、环境变量等
- 考虑引号、转义字符的处理
- 可能需要调用系统shell进行解析

这个设计为shell命令执行库提供了一个坚实的基础，既保证了功能的完整性，又确保了代码的可维护性和扩展性。通过支持两种命令创建方式，可以满足从简单脚本到复杂shell命令的各种使用场景。