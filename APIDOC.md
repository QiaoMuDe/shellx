# Package shellx

```go
import "gitee.com/MM-Q/shellx"
```

Package shellx 定义了shell命令执行库的核心数据类型。 本文件定义了Builder结构体和相关构造函数，提供链式调用API来构建命令对象。

Builder是命令构建器的核心实现，支持：
- 三种命令创建方式：NewCmd、NewCmds、NewCmdStr
- 链式调用设置：工作目录、环境变量、超时、上下文、标准输入输出、Shell类型
- 并发安全的读写操作
- 灵活的命令配置和构建

Package shellx 定义了shell命令执行库的核心数据类型。 本文件定义了Command结构体，封装了exec.Cmd并提供了丰富的命令执行功能。

Command是命令执行对象的核心实现，支持：
- 同步执行：Exec、ExecOutput、ExecStdout、ExecResult
- 异步执行：ExecAsync、Wait
- 进程控制：Kill、Signal、IsRunning、GetPID
- 执行状态管理：IsExecuted（确保命令只执行一次）
- 完整的执行结果：Result对象包含输出、错误、时间、退出码等信息

Package shellx 提供了一个功能完善、易于使用的Go语言shell命令执行库。

本库基于Go标准库的os/exec包进行封装，提供了更加友好的API和丰富的功能， 支持同步和异步命令执行、输入输出重定向、超时控制、上下文管理、
多种shell类型支持等功能，并提供类型安全的API和友好的链式调用接口。

## 主要特性

- 支持三种命令创建方式：NewCmd(可变参数)、NewCmds(切片)、NewCmdStr(字符串解析)
- 链式调用API，支持流畅的方法链
- 完整的错误处理和类型安全
- 支持多种shell类型（sh、bash、cmd、powershell、pwsh等）
- 同步和异步执行支持
- 命令执行状态管理和进程控制
- 输入输出重定向和环境变量设置
- 超时控制和上下文取消
- 并发安全的设计
- 跨平台兼容（Windows、Linux、macOS）

## 核心组件

- Builder: 命令构建器，提供链式调用API
- Command: 命令执行对象，封装exec.Cmd并提供额外功能
- Result: 命令执行结果，包含输出、错误、时间等信息
- ShellType: Shell类型枚举，支持多种shell

## 基本用法

```go
import "gitee.com/MM-Q/shellx"

// 方式1：使用可变参数创建命令
cmd := shellx.NewCmd("ls", "-la").
	WithWorkDir("/tmp").
	WithTimeout(30 * time.Second).
	WithShell(shellx.ShellBash).
	Build()

// 方式2：使用字符串创建命令
cmd := shellx.NewCmdStr(`echo "hello world"`).
	WithEnv("MY_VAR", "value").
	Build()

// 同步执行
err := cmd.Exec()
if err != nil {
	log.Fatal(err)
}

// 获取输出
output, err := cmd.ExecOutput()
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(output))

// 获取完整结果
result := cmd.ExecResult()
fmt.Printf("Exit Code: %d\n", result.Code())
fmt.Printf("Success: %t\n", result.Success())
fmt.Printf("Duration: %v\n", result.Duration())
fmt.Printf("Output: %s\n", result.Output())

// 异步执行
err = cmd.ExecAsync()
if err != nil {
	log.Fatal(err)
}
// 等待完成
err = cmd.Wait()
```

## 高级用法

```go
// 设置标准输入输出
var stdout, stderr bytes.Buffer
stdin := strings.NewReader("input data")

cmd := shellx.NewCmd("cat").
	WithStdin(stdin).
	WithStdout(&stdout).
	WithStderr(&stderr).
	Build()

// 使用上下文控制
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

cmd := shellx.NewCmd("long-running-command").
	WithContext(ctx).
	Build()

// 进程控制
cmd.ExecAsync()
pid := cmd.GetPID()
isRunning := cmd.IsRunning()
cmd.Kill() // 或 cmd.Signal(syscall.SIGTERM)
```

## 命令解析

```go
// 支持复杂的命令字符串解析，包括引号处理
cmd := shellx.NewCmdStr(`git commit -m "Initial commit with 'quotes'"`)
// 解析结果：["git", "commit", "-m", "Initial commit with 'quotes'"]
```

## Shell类型

```go
// 支持多种shell类型
shellx.ShellSh         // sh shell
shellx.ShellBash       // bash shell
shellx.ShellCmd        // Windows cmd
shellx.ShellPowerShell // Windows PowerShell
shellx.ShellPwsh       // PowerShell Core
shellx.ShellNone       // 直接执行，不使用shell
shellx.ShellDefault    // 根据操作系统自动选择
```

## 注意事项

- 每个Command对象只能执行一次，重复执行会返回错误
- Builder是并发安全的，可以在多个goroutine中安全使用
- 命令执行会继承父进程的环境变量，可通过WithEnv添加额外变量
- 超时设置仅在支持的Go版本中有效
- 异步执行需要调用Wait()等待完成或使用Kill()终止

Package shellx 定义了shell命令执行库的核心数据类型。
本文件定义了ShellType枚举和Result结构体，提供了shell类型管理和执行结果封装。

## 主要类型

- ShellType: Shell类型枚举，支持sh、bash、cmd、powershell等多种shell
- Result: 命令执行结果结构体，包含退出码、输出、时间、错误等完整信息

## ShellType支持的shell类型

- ShellSh: Unix/Linux sh shell
- ShellBash: Bash shell
- ShellCmd: Windows Command Prompt
- ShellPowerShell: Windows PowerShell
- ShellPwsh: PowerShell Core (跨平台)
- ShellNone: 直接执行命令，不使用shell
- ShellDefault: 根据操作系统自动选择默认shell

Package shellx 定义了shell命令执行库的核心数据类型。 本文件定义了工具函数，提供命令字符串处理和解析功能。

## 主要功能

- getCmdStr: 从Builder对象获取完整的命令字符串
- ParseCmd: 智能解析命令字符串，支持复杂的引号处理
- FindCmd: 查找系统中的命令路径

## ParseCmd函数特性

- 支持单引号、双引号、反引号三种引号类型
- 正确处理引号内的空格和特殊字符
- 支持引号嵌套（不同类型的引号可以嵌套）
- 自动检测未闭合的引号并返回空结果
- 处理多个连续空格和制表符
- 支持复杂的命令行参数解析

## 解析示例

- `ls -la` → ["ls", "-la"]
- `echo "hello world"` → ["echo", "hello world"]
- `git commit -m "fix: update 'config' file"` → ["git", "commit", "-m", "fix: update 'config' file"]
- `find . -name "*.go" -exec grep "pattern" {} \;` → ["find", ".", "-name", "*.go", "-exec", "grep", "pattern", "{}", "\\;"]

## FUNCTIONS

### func FindCmd

```go
func FindCmd(name string) (string, error)
```

FindCmd 查找命令

**参数:**
- name: 命令名称

**返回:**
- string: 命令路径
- error: 错误信息

### func Exec

```go
func Exec(name string, args ...string) error
```

Exec 执行命令(阻塞)

**参数:**
- name: 命令名
- args: 命令参数

**返回:**
- error: 错误信息

### func ExecOutput

```go
func ExecOutput(name string, args ...string) ([]byte, error)
```

ExecOutput 执行命令并返回合并后的输出(阻塞)

**参数:**
- name: 命令名
- args: 命令参数

**返回:**
- []byte: 输出
- error: 错误信息

### func ExecOutputStr

```go
func ExecOutputStr(cmdStr string) ([]byte, error)
```

ExecOutputStr 执行命令并返回合并后的输出(阻塞)

**参数:**
- cmdStr: 命令字符串

**返回:**
- []byte: 输出
- error: 错误信息

### func ExecStr

```go
func ExecStr(cmdStr string) error
```

ExecStr 执行命令(阻塞)

**参数:**
- cmdStr: 命令字符串

**返回:**
- error: 错误信息

### func FindCmd

```go
func FindCmd(name string) (string, error)
```

FindCmd 查找命令

**参数:**
- name: 命令名称

**返回:**
- string: 命令路径
- error: 错误信息

### func ParseCmd

```go
func ParseCmd(cmdStr string) []string
```

ParseCmd 将命令字符串解析为命令切片，支持引号处理(单引号、双引号、反引号)，出错时返回空切片

**实现原理：**
1. 去除首尾空白
2. 遍历每个字符
3. 处理引号状态切换
4. 在非引号状态下遇到空格时分割
5. 检查引号是否闭合

**参数:**
- cmdStr: 要解析的命令字符串

**返回值:**
- []string: 解析后的命令切片

## TYPES

### type Builder

```go
type Builder struct {
	// Has unexported fields.
}
```

Builder 命令构建器，提供链式调用

#### func NewCmd

```go
func NewCmd(name string, args ...string) *Builder
```

NewCmd 创建新的命令构建器 (数组方式 - 可变参数)

**参数：**
- name: 命令名
- args: 命令参数列表

**返回：**
- *Builder: 命令构建器对象

#### func NewCmdStr

```go
func NewCmdStr(cmdStr string) *Builder
```

NewCmdStr 创建新的命令构建器 (字符串方式)

**参数：**
- cmdStr: 命令字符串

**返回：**
- *Builder: 命令构建器对象

#### func NewCmds

```go
func NewCmds(cmdArgs []string) *Builder
```

NewCmds 创建新的命令构建器 (数组方式 - 切片参数，第一个元素为命令名)

**参数：**
- cmdArgs: 命令参数列表，第一个元素为命令名，后续元素为参数

**返回：**
- *Builder: 命令构建器对象

#### func (*Builder) Args

```go
func (b *Builder) Args() []string
```

Args 获取命令参数列表

**返回:**
- []string: 命令参数列表

#### func (*Builder) Build

```go
func (b *Builder) Build() *Command
```

Build 构建并返回命令对象

**返回:**
- *Command: 构建的命令对象

#### func (*Builder) Env

```go
func (b *Builder) Env() []string
```

Env 获取命令环境变量列表

**返回:**
- []string: 命令环境变量列表

#### func (*Builder) Name

```go
func (b *Builder) Name() string
```

Name 获取命令名称

**返回:**
- string: 命令名称

#### func (*Builder) Raw

```go
func (b *Builder) Raw() string
```

Raw 获取原始命令字符串

**返回:**
- string: 原始命令字符串

#### func (*Builder) ShellType

```go
func (b *Builder) ShellType() ShellType
```

ShellType 获取shell类型

**返回:**
- ShellType: shell类型

#### func (*Builder) Timeout

```go
func (b *Builder) Timeout() time.Duration
```

Timeout 获取命令执行超时时间

**返回:**
- time.Duration: 命令执行超时时间

#### func (*Builder) WithContext

```go
func (b *Builder) WithContext(ctx context.Context) *Builder
```

WithContext 设置命令的上下文

**参数：**
- ctx: context.Context类型，用于取消命令执行和超时控制

**返回：**
- *Builder: 命令构建器对象

#### func (*Builder) WithEnv

```go
func (b *Builder) WithEnv(key, value string) *Builder
```

WithEnv 设置命令的环境变量

**参数：**
- key: 环境变量的键
- value: 环境变量的值

**返回：**
- *Builder: 命令构建器对象

#### func (*Builder) WithShell

```go
func (b *Builder) WithShell(shell ShellType) *Builder
```

WithShell 设置命令的shell类型

**参数：**
- shell: ShellType类型，表示要使用的shell类型

**返回：**
- *Builder: 命令构建器对象

#### func (*Builder) WithStderr

```go
func (b *Builder) WithStderr(stderr io.Writer) *Builder
```

WithStderr 设置命令的标准错误输出

**参数：**
- stderr: io.Writer类型，用于接收命令的标准错误输出

**返回：**
- *Builder: 命令构建器对象

#### func (*Builder) WithStdin

```go
func (b *Builder) WithStdin(stdin io.Reader) *Builder
```

WithStdin 设置命令的标准输入

**参数：**
- stdin: io.Reader类型，用于提供命令的标准输入

**返回：**
- *Builder: 命令构建器对象

#### func (*Builder) WithStdout

```go
func (b *Builder) WithStdout(stdout io.Writer) *Builder
```

WithStdout 设置命令的标准输出

**参数：**
- stdout: io.Writer类型，用于接收命令的标准输出

**返回：**
- *Builder: 命令构建器对象

#### func (*Builder) WithTimeout

```go
func (b *Builder) WithTimeout(timeout time.Duration) *Builder
```

WithTimeout 设置命令的超时时间

**参数：**
- timeout: time.Duration类型，命令执行的超时时间

**返回：**
- *Builder: 命令构建器对象

#### func (*Builder) WithWorkDir

```go
func (b *Builder) WithWorkDir(dir string) *Builder
```

WithWorkDir 设置命令的工作目录

**参数：**
- dir: 命令的工作目录

**返回：**
- *Builder: 命令构建器对象

#### func (*Builder) WorkDir

```go
func (b *Builder) WorkDir() string
```

WorkDir 获取命令执行的工作目录

**返回:**
- string: 命令执行目录

### type Command

```go
type Command struct {
	// Has unexported fields.
}
```

Command 命令对象

#### func (*Command) Cmd

```go
func (c *Command) Cmd() *exec.Cmd
```

Cmd 获取底层的 exec.Cmd 对象

**返回:**
- *exec.Cmd: 底层的 exec.Cmd 对象

#### func (*Command) Exec

```go
func (c *Command) Exec() error
```

Exec 执行命令(阻塞)

**返回:**
- error: 错误信息

#### func (*Command) ExecAsync

```go
func (c *Command) ExecAsync() error
```

ExecAsync 异步执行命令(非阻塞)

**返回:**
- error: 错误信息

#### func (*Command) ExecOutput

```go
func (c *Command) ExecOutput() ([]byte, error)
```

ExecOutput 执行命令并返回合并后的输出(阻塞)

**返回:**
- []byte: 命令输出
- error: 错误信息

#### func (*Command) ExecResult

```go
func (c *Command) ExecResult() (*Result, error)
```

ExecResult 执行命令并返回完整的执行结果(阻塞)

**使用示例:**

```go
result, err := cmd.ExecResult()
if err != nil {
    // 处理错误情况
    log.Printf("Command failed: %v", err)
    return
}
// 处理成功情况
fmt.Println(string(result.Output()))
```

**返回:**
- *Result: 执行结果对象，包含输出、时间、退出码等信息
- error: 执行过程中的错误信息

#### func (*Command) ExecStdout

```go
func (c *Command) ExecStdout() ([]byte, error)
```

ExecStdout 执行命令并返回标准输出(阻塞)

**返回:**
- []byte: 标准输出
- error: 错误信息

#### func (*Command) GetPID

```go
func (c *Command) GetPID() int
```

GetPID 获取进程ID

**返回:**
- int: 进程ID，如果进程不存在返回0

#### func (*Command) IsExecuted

```go
func (c *Command) IsExecuted() bool
```

IsExecuted 检查命令是否已经执行过

**返回:**
- bool: 是否已执行

#### func (*Command) IsRunning

```go
func (c *Command) IsRunning() bool
```

IsRunning 检查进程是否还在运行

**返回:**
- bool: 是否在运行

#### func (*Command) Kill

```go
func (c *Command) Kill() error
```

Kill 杀死当前命令的进程

**返回:**
- error: 错误信息

#### func (*Command) Signal

```go
func (c *Command) Signal(sig os.Signal) error
```

Signal 向当前进程发送信号

**参数:**
- sig: 信号类型

**返回:**
- error: 错误信息

#### func (*Command) Wait

```go
func (c *Command) Wait() error
```

Wait 等待命令执行完成(仅在异步执行时有效)

**返回:**
- error: 错误信息

### type Result

```go
type Result struct {
	// Has unexported fields.
}
```

Result 表示命令执行的结果

#### func (*Result) Code

```go
func (r *Result) Code() int
```

提供公共访问方法

#### func (*Result) Duration

```go
func (r *Result) Duration() time.Duration
```

#### func (*Result) End

```go
func (r *Result) End() time.Time
```

#### func (*Result) Output

```go
func (r *Result) Output() []byte
```

#### func (*Result) Start

```go
func (r *Result) Start() time.Time
```

#### func (*Result) Success

```go
func (r *Result) Success() bool
```

### type ShellType

```go
type ShellType int
```

ShellType 定义shell类型

```go
const (
	ShellSh         ShellType = iota // sh shell
	ShellBash                        // bash shell
	ShellPwsh                        // pwsh (PowerShell Core)
	ShellPowerShell                  // powershell (Windows PowerShell)
	ShellCmd                         // cmd (Windows Command Prompt)
	ShellNone                        // 无shell, 直接原生的执行命令
	ShellDefault                     // 默认shell, 根据操作系统自动选择(Windows系统默认为cmd, 其他系统默认为sh)
)
```

#### func (ShellType) String

```go
func (s ShellType) String() string
```

String 返回shell类型的字符串表示

