# ShellX 🚀

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](https://gitee.com/MM-Q/shellx/blob/master/LICENSE)
[![Gitee](https://img.shields.io/badge/Gitee-Repository-red?style=for-the-badge&logo=gitee)](https://gitee.com/MM-Q/shellx)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](https://gitee.com/MM-Q/shellx)

**一个功能完善、易于使用的 Go 语言 Shell 命令执行库**

[🏠 仓库地址](https://gitee.com/MM-Q/shellx) • [📖 文档](#api-文档概述) • [🚀 快速开始](#安装指南) • [💡 示例](#使用示例)

</div>

---

## 📋 项目简介

ShellX 是一个基于 Go 标准库 `os/exec` 包封装的高级命令执行库，提供了更加友好的 API 和丰富的功能。它支持同步和异步命令执行、输入输出重定向、超时控制、上下文管理、多种 shell 类型支持等功能，并提供类型安全的 API 和友好的链式调用接口。

无论您是需要执行简单的系统命令，还是构建复杂的命令行工具，ShellX 都能为您提供强大而灵活的解决方案。

## ✨ 核心特性

| 特性 | 描述 |
|------|------|
| 🔧 **多种创建方式** | 支持 `NewCmd`(可变参数)、`NewCmds`(切片)、`NewCmdStr`(字符串解析) 三种命令创建方式 |
| ⚡ **便捷函数** | 提供 `Exec`、`ExecStr`、`ExecOutput`、`ExecOutputStr` 等便捷函数，无需手动构建对象 |
| ⛓️ **链式调用 API** | 流畅的方法链，支持工作目录、环境变量、超时等配置 |
| 🛡️ **类型安全** | 完整的错误处理和类型安全保证 |
| 🐚 **多 Shell 支持** | 支持 sh、bash、cmd、powershell、pwsh 等多种 shell 类型 |
| ⚡ **同步/异步执行** | 灵活的执行模式，支持阻塞和非阻塞操作 |
| 🎛️ **进程控制** | 完整的进程生命周期管理，支持信号发送、进程终止等 |
| 📊 **执行状态管理** | 智能的执行状态跟踪，防止重复执行 |
| 🔄 **输入输出重定向** | 灵活的标准输入输出配置 |
| ⏱️ **超时控制** | 支持命令执行超时和上下文取消 |
| 🔒 **并发安全** | 线程安全的设计，支持多 goroutine 环境 |
| 🌍 **跨平台兼容** | 支持 Windows、Linux、macOS 等主流操作系统 |
| 🧠 **智能解析** | 强大的命令字符串解析，支持复杂引号处理 |

## 📦 安装指南

### 使用 Go Modules（推荐）

```bash
go get gitee.com/MM-Q/shellx
```

### 版本要求

- Go 1.25.0 或更高版本
- 支持 Go Modules

## 🚀 使用示例

### 基础用法

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "gitee.com/MM-Q/shellx"
)

func main() {
    // 方式1：使用可变参数创建命令
    cmd := shellx.NewCmd("echo", "Hello, World!").
        WithTimeout(10 * time.Second).
        Build()
    
    // 同步执行
    err := cmd.Exec()
    if err != nil {
        log.Fatal(err)
    }
    
    // 获取输出
    output, err := shellx.NewCmd("ls", "-la").Build().ExecOutput()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(output))
}
```

### 便捷函数用法

```go
package main

import (
    "fmt"
    "log"
    
    "gitee.com/MM-Q/shellx"
)

func main() {
    // 直接执行命令（可变参数方式）
    err := shellx.Exec("echo", "Hello, World!")
    if err != nil {
        log.Fatal(err)
    }
    
    // 直接执行命令（字符串方式）
    err = shellx.ExecStr("ls -la")
    if err != nil {
        log.Fatal(err)
    }
    
    // 执行命令并获取输出（可变参数方式）
    output, err := shellx.ExecOutput("pwd")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Current directory: %s", output)
    
    // 执行命令并获取输出（字符串方式）
    output, err = shellx.ExecOutputStr("git status --porcelain")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Git status: %s", output)
}
```

### 字符串解析

```go
// 方式2：使用字符串创建命令（支持复杂引号处理）
cmd := shellx.NewCmdStr(`git commit -m "feat: add new feature with 'quotes'"`).
    WithWorkDir("/path/to/repo").
    WithEnv("GIT_AUTHOR_NAME", "John Doe").
    Build()

result, err := cmd.ExecResult()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Exit Code: %d\n", result.Code())
fmt.Printf("Success: %t\n", result.Success())
fmt.Printf("Duration: %v\n", result.Duration())
```

### 高级用法

```go
package main

import (
    "bytes"
    "context"
    "strings"
    "time"
    
    "gitee.com/MM-Q/shellx"
)

func advancedExample() {
    // 设置标准输入输出
    var stdout, stderr bytes.Buffer
    stdin := strings.NewReader("input data\n")
    
    cmd := shellx.NewCmd("cat").
        WithStdin(stdin).
        WithStdout(&stdout).
        WithStderr(&stderr).
        WithWorkDir("/tmp").
        WithEnv("MY_VAR", "custom_value").
        WithShell(shellx.ShellBash).
        Build()
    
    // 使用上下文控制
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    cmdWithContext := shellx.NewCmd("long-running-command").
        WithContext(ctx).
        Build()
    
    // 异步执行
    err := cmdWithContext.ExecAsync()
    if err != nil {
        log.Fatal(err)
    }
    
    // 进程控制
    pid := cmdWithContext.GetPID()
    fmt.Printf("Process ID: %d\n", pid)
    
    if cmdWithContext.IsRunning() {
        fmt.Println("Command is still running...")
        // 可以选择等待或终止
        // cmdWithContext.Kill()
        // 或发送信号
        // cmdWithContext.Signal(syscall.SIGTERM)
    }
    
    // 等待完成
    err = cmdWithContext.Wait()
    if err != nil {
        log.Printf("Command failed: %v", err)
    }
}
```

### 不同 Shell 类型示例

```go
// 使用不同的 Shell 类型
examples := map[string]shellx.ShellType{
    "Bash":       shellx.ShellBash,
    "PowerShell": shellx.ShellPwsh,
    "CMD":        shellx.ShellCmd,
    "直接执行":      shellx.ShellNone,
    "系统默认":      shellx.ShellDefault,
}

for name, shellType := range examples {
    cmd := shellx.NewCmdStr("echo 'Hello from " + name + "'").
        WithShell(shellType).
        Build()
    
    output, err := cmd.ExecOutput()
    if err != nil {
        fmt.Printf("%s failed: %v\n", name, err)
        continue
    }
    fmt.Printf("%s: %s", name, output)
}
```

## 📚 API 文档概述

### 核心类型

| 类型 | 描述 |
|------|------|
| `Builder` | 命令构建器，提供链式调用 API |
| `Command` | 命令执行对象，封装 exec.Cmd 并提供额外功能 |
| `Result` | 命令执行结果，包含输出、错误、时间等信息 |
| `ShellType` | Shell 类型枚举，支持多种 shell |

### 主要方法

#### 创建命令

```go
// 可变参数方式
func NewCmd(name string, args ...string) *Builder

// 切片方式
func NewCmds(cmdArgs []string) *Builder

// 字符串解析方式
func NewCmdStr(cmdStr string) *Builder
```

#### 链式配置

```go
func (b *Builder) WithWorkDir(dir string) *Builder
func (b *Builder) WithEnv(key, value string) *Builder
func (b *Builder) WithTimeout(timeout time.Duration) *Builder
func (b *Builder) WithContext(ctx context.Context) *Builder
func (b *Builder) WithStdin(stdin io.Reader) *Builder
func (b *Builder) WithStdout(stdout io.Writer) *Builder
func (b *Builder) WithStderr(stderr io.Writer) *Builder
func (b *Builder) WithShell(shell ShellType) *Builder
func (b *Builder) Build() *Command
```

#### 便捷函数

```go
// 直接执行命令
func Exec(name string, args ...string) error
func ExecStr(cmdStr string) error

// 执行命令并获取输出
func ExecOutput(name string, args ...string) ([]byte, error)
func ExecOutputStr(cmdStr string) ([]byte, error)
```

#### 命令执行

```go
// 同步执行
func (c *Command) Exec() error
func (c *Command) ExecOutput() ([]byte, error)
func (c *Command) ExecStdout() ([]byte, error)
func (c *Command) ExecResult() (*Result, error)

// 异步执行
func (c *Command) ExecAsync() error
func (c *Command) Wait() error

// 进程控制
func (c *Command) Kill() error
func (c *Command) Signal(sig os.Signal) error
func (c *Command) IsRunning() bool
func (c *Command) GetPID() int
func (c *Command) IsExecuted() bool
```

## 🎯 支持的功能

### Shell 类型支持

| Shell 类型 | 常量 | 平台支持 | 描述 |
|------------|------|----------|------|
| **sh** | `ShellSh` | Unix/Linux/macOS | 标准 Unix shell |
| **bash** | `ShellBash` | Unix/Linux/macOS | Bash shell |
| **cmd** | `ShellCmd` | Windows | Windows 命令提示符 |
| **powershell** | `ShellPowerShell` | Windows | Windows PowerShell |
| **pwsh** | `ShellPwsh` | 跨平台 | PowerShell Core |
| **none** | `ShellNone` | 跨平台 | 直接执行，不使用 shell |
| **default** | `ShellDefault` | 跨平台 | 根据操作系统自动选择 |

### 命令解析特性

- ✅ 单引号 `'` 支持
- ✅ 双引号 `"` 支持  
- ✅ 反引号 `` ` `` 支持
- ✅ 引号嵌套处理
- ✅ 转义字符支持
- ✅ 多空格处理
- ✅ 制表符处理
- ✅ 未闭合引号检测

### 执行模式

- 🔄 **同步执行**：阻塞等待命令完成
- ⚡ **异步执行**：非阻塞启动，可后续等待
- 📊 **结果获取**：完整的执行结果信息
- 🎯 **输出捕获**：标准输出和错误输出

## ⚙️ 配置选项

### 环境配置

```go
cmd := shellx.NewCmd("command").
    WithWorkDir("/custom/path").           // 设置工作目录
    WithEnv("KEY", "value").              // 添加环境变量
    WithTimeout(30 * time.Second).        // 设置超时时间
    WithContext(ctx).                     // 设置上下文
    Build()
```

### 输入输出配置

```go
var stdout, stderr bytes.Buffer
stdin := strings.NewReader("input")

cmd := shellx.NewCmd("command").
    WithStdin(stdin).                     // 设置标准输入
    WithStdout(&stdout).                  // 设置标准输出
    WithStderr(&stderr).                  // 设置标准错误
    Build()
```

### Shell 配置

```go
cmd := shellx.NewCmd("command").
    WithShell(shellx.ShellBash).          // 指定 shell 类型
    Build()
```

## 📁 项目结构

```
shellx/
├── README.md                 # 项目文档
├── LICENSE                   # MIT 许可证
├── go.mod                    # Go 模块文件
├── shellx.go                 # 主包文档
├── builder.go                # 命令构建器实现
├── command.go                # 命令执行对象实现
├── types.go                  # 类型定义
├── utils.go                  # 工具函数
├── builder_test.go           # Builder 测试
├── command_test.go           # Command 测试
├── utils_test.go             # 工具函数测试
├── builder_bench_test.go     # Builder 基准测试
├── command_bench_test.go     # Command 基准测试
└── TEST_README.md            # 测试说明文档
```

## 🧪 测试说明

### 运行测试

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestBuilder
go test -v -run TestCommand

# 运行基准测试
go test -bench=. -benchmem

# 查看测试覆盖率
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 测试覆盖

- ✅ **功能测试**：完整的功能正确性测试
- ✅ **边界测试**：边界条件和错误处理测试
- ✅ **并发测试**：并发安全性测试
- ✅ **性能测试**：基准测试和性能监控
- ✅ **模糊测试**：命令解析的模糊测试

### 测试结果

最新测试结果显示所有测试通过，性能表现良好：

- 功能测试：15/15 通过 ✅
- 命令测试：14/14 通过 ✅
- 基准测试：性能优异 🚀

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

```
MIT License

Copyright (c) 2025 M乔木

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. **Fork** 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 **Pull Request**

### 贡献类型

- 🐛 **Bug 修复**
- ✨ **新功能**
- 📚 **文档改进**
- 🎨 **代码优化**
- 🧪 **测试增强**
- 🔧 **工具改进**

### 开发规范

- 遵循 Go 语言编码规范
- 添加适当的测试用例
- 更新相关文档
- 确保所有测试通过

## 📞 联系方式

- **作者**：M乔木
- **仓库**：[https://gitee.com/MM-Q/shellx](https://gitee.com/MM-Q/shellx)
- **问题反馈**：[Issues](https://gitee.com/MM-Q/shellx/issues)
- **功能请求**：[Feature Requests](https://gitee.com/MM-Q/shellx/issues)

## 🔗 相关链接

- 📖 [Go 官方文档](https://golang.org/doc/)
- 🔧 [os/exec 包文档](https://pkg.go.dev/os/exec)
- 🏠 [项目主页](https://gitee.com/MM-Q/shellx)
- 📋 [更新日志](https://gitee.com/MM-Q/shellx/releases)

---

<div align="center">

**如果这个项目对您有帮助，请给它一个 ⭐ Star！**

[⬆️ 回到顶部](#shellx-)

</div>