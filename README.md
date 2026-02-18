<div align="center">

# ShellX 🚀


[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](https://gitee.com/MM-Q/shellx/blob/master/LICENSE)
[![Gitee](https://img.shields.io/badge/Gitee-Repository-red?style=for-the-badge&logo=gitee)](https://gitee.com/MM-Q/shellx)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](https://gitee.com/MM-Q/shellx)

**一个功能完善、易于使用的 Go 语言 Shell 命令执行库**

[🏠 仓库地址](https://gitee.com/MM-Q/shellx) • [📖 API文档](APIDOC.md) • [🚀 快速开始](#安装指南) • [💡 示例](#使用示例)

</div>

---

## 📋 项目简介

ShellX 是一个功能完善、易于使用的 Go 语言 Shell 命令执行库。本项目包含两个子包：

- **主包 (shellx)**：基于 Go 标准库 `os/exec` 包封装的高级命令执行库，提供了更加友好的 API 和丰富的功能
- **子包 (shx)**：基于 [mvdan.cc/sh/v3](https://mvdan.cc/sh/v3) 的纯 Go shell 命令执行功能，具有更好的跨平台一致性

无论您是需要执行简单的系统命令，还是构建复杂的命令行工具，ShellX 都能为您提供强大而灵活的解决方案。

## ✨ 核心特性

### 主包 (shellx) 特性

| 特性 | 描述 |
|------|------|
| 🎯 **一体化设计** | Command集配置、构建、执行于一体，无需Build()方法，简化API使用 |
| 🔧 **多种创建方式** | 支持 `NewCmd`(可变参数)、`NewCmds`(切片)、`NewCmdStr`(字符串解析) 三种命令创建方式 |
| ⚡ **丰富便捷函数** | 提供 `Exec`、`ExecStr`、`ExecOut`、`ExecOutStr` 及其带超时版本，开箱即用 |
| ⛓️ **链式调用 API** | 流畅的方法链，支持工作目录、环境变量、超时等配置 |
| ⏱️ **精确超时控制** | 延迟构建exec.Cmd，确保超时计时精确，避免配置时间损耗 |
| 🛡️ **类型安全** | 完整的错误处理和类型安全保证 |
| 🐚 **多 Shell 支持** | 支持 sh、bash、cmd、powershell、pwsh 等多种 shell 类型 |
| ⚡ **同步/异步执行** | 灵活的执行模式，支持阻塞和非阻塞操作 |
| 🎛️ **进程控制** | 完整的进程生命周期管理，支持信号发送、进程终止等 |
| 📊 **执行状态管理** | 智能的执行状态跟踪，防止重复执行 |
| 🔄 **输入输出重定向** | 灵活的标准输入输出配置 |
| 🔒 **并发安全** | 线程安全的设计，支持多 goroutine 环境 |
| 🌍 **跨平台兼容** | 支持 Windows、Linux、macOS 等主流操作系统 |
| 🧠 **智能解析** | 强大的命令字符串解析，支持复杂引号处理 |

### 子包 (shx) 特性

| 特性 | 描述 |
|------|------|
| 🟢 **纯 Go 实现** | 基于 mvdan.cc/sh/v3，不依赖系统 shell |
| 🌍 **跨平台一致** | Windows/Linux/macOS 行为完全一致 |
| 🔒 **轻量级并发** | 使用 atomic.Bool 防止重复执行 |
| ⛓️ **链式调用** | 支持流畅的方法链配置 |
| ⏱️ **超时控制** | 支持上下文超时和超时参数 |
| 🔄 **输入输出重定向** | 灵活的标准输入输出配置 |

## 📦 安装指南

### 使用 Go Modules（推荐）

```bash
# 安装主包
go get gitee.com/MM-Q/shellx

# 安装子包 shx (基于 mvdan.cc/sh/v3)
go get gitee.com/MM-Q/shellx/shx
```

### 版本要求

- Go 1.25.0 或更高版本
- 支持 Go Modules

### 包说明

ShellX 项目包含两个包，您可以根据需求选择使用：

| 包 | 导入路径 | 特点 |
|----|----------|------|
| 主包 | `gitee.com/MM-Q/shellx` | 基于 os/exec，功能丰富，支持进程控制 |
| 子包 | `gitee.com/MM-Q/shellx/shx` | 纯 Go 实现，跨平台一致性好 |

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
    // 使用可变参数创建命令
    err := shellx.NewCmd("echo", "Hello, World!").
        WithTimeout(10 * time.Second).
        Exec()
    if err != nil {
        log.Fatal(err)
    }
    
    // 获取输出
    output, err := shellx.NewCmd("ls", "-la").ExecOutput()
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
    "time"
    
    "gitee.com/MM-Q/shellx"
)

func main() {
    // 基础执行函数
    err := shellx.Exec("echo", "Hello, World!")        // 执行命令，输出到控制台
    err = shellx.ExecStr("ls -la")                      // 字符串方式执行
    
    // 获取输出的函数
    output, err := shellx.ExecOut("pwd")                // 执行并返回输出
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Current directory: %s", output)
    
    output, err = shellx.ExecOutStr("git status --porcelain") // 字符串方式执行并返回输出
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Git status: %s", output)
    
    // 带超时的执行函数
    err = shellx.ExecT(5*time.Second, "sleep", "10")                    // 5秒超时
    err = shellx.ExecStrT(3*time.Second, "ping google.com")             // 字符串方式，3秒超时
    output, err = shellx.ExecOutT(2*time.Second, "curl", "example.com") // 返回输出，2秒超时
    output, err = shellx.ExecOutStrT(1*time.Second, "date")             // 字符串方式，返回输出，1秒超时
}
```

### 字符串解析

```go
// 使用字符串创建命令（支持复杂引号处理）
cmd := shellx.NewCmdStr(`git commit -m "feat: add new feature with 'quotes'"`).
    WithWorkDir("/path/to/repo").
    WithEnv("GIT_AUTHOR_NAME", "John Doe")

// 执行命令
err := cmd.Exec()
if err != nil {
    log.Fatal(err)
}

// 如果需要获取退出码，可以使用 WaitWithCode
exitCode, err := cmd.WaitWithCode()
if err != nil {
    log.Printf("Command failed: %v", err)
}
fmt.Printf("Exit Code: %d\n", exitCode)
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
        WithShell(shellx.ShellBash)
    
    // 使用上下文控制
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    cmdWithContext := shellx.NewCmd("long-running-command").
        WithContext(ctx)
    
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
    "Bash":        shellx.ShellBash,
    "PowerShell":  shellx.ShellPwsh,
    "CMD":         shellx.ShellCmd,
    "直接执行":    shellx.ShellNone,
    "系统默认1":   shellx.ShellDef1,
    "系统默认2":   shellx.ShellDef2,
}

for name, shellType := range examples {
    cmd := shellx.NewCmdStr("echo 'Hello from " + name + "'").
        WithShell(shellType)
    
    output, err := cmd.ExecOutput()
    if err != nil {
        fmt.Printf("%s failed: %v\n", name, err)
        continue
    }
    fmt.Printf("%s: %s", name, output)
}
```

### 子包 shx 使用示例

shx 子包提供基于 mvdan.cc/sh/v3 的纯 Go shell 执行功能，具有更好的跨平台一致性。

#### 基础用法

```go
package main

import (
    "fmt"
    "log"
    
    "gitee.com/MM-Q/shellx/shx"
)

func main() {
    // 简单执行
    err := shx.Run("echo hello world")
    if err != nil {
        log.Fatal(err)
    }
    
    // 获取输出
    output, err := shx.Out("ls -la")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(output))
}
```

#### 链式配置

```go
// 使用链式配置
output, err := shx.New("echo hello").
    WithTimeout(5 * time.Second).
    WithDir("/tmp").
    ExecOutput()
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(output))
```

#### 使用上下文

```go
// 使用上下文控制超时
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := shx.New("long-running-command").
    WithContext(ctx).
    Exec()
if err != nil {
    log.Printf("Command failed: %v", err)
}
```

#### 自定义输入输出

```go
// 自定义标准输入输出
var stdout, stderr bytes.Buffer
stdin := strings.NewReader("hello")

err := shx.New("cat").
    WithStdin(stdin).
    WithStdout(&stdout).
    WithStderr(&stderr).
    Exec()
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(stdout.Bytes()))
```

#### 检查退出状态

```go
// 检查命令退出状态
err := shx.Run("exit 5")
if exitCode, ok := shx.IsExitStatus(err); ok {
    fmt.Printf("Command exited with code: %d\n", exitCode)
}
```

## 🎯 支持的功能

### 主包 Shell 类型支持

- **sh** - 标准 Unix shell
- **bash** - Bash shell  
- **cmd** - Windows 命令提示符
- **powershell** - Windows PowerShell
- **pwsh** - PowerShell Core (跨平台)
- **none** - 直接执行，不使用 shell
- **def1** - 根据操作系统自动选择(Windows系统默认为cmd, 其他系统默认为sh)
- **def2** - 根据操作系统自动选择(Windows系统默认为powershell, 其他系统默认为sh)

### 子包 shx Shell 类型支持

shx 子包使用 mvdan.cc/sh/v3 作为解析器，支持以下 shell 类型：
- **bash** - Bash shell (默认)
- **sh** - POSIX shell 

### 命令解析特性

- ✅ 单引号、双引号、反引号支持
- ✅ 引号嵌套处理
- ✅ 转义字符支持
- ✅ 多空格和制表符处理
- ✅ 未闭合引号检测

### 执行模式

- 🔄 **同步执行**：阻塞等待命令完成
- ⚡ **异步执行**（仅主包）：非阻塞启动，可后续等待
- 📊 **结果获取**：完整的执行结果信息
- 🎯 **输出捕获**：标准输出和错误输出
- ⏱️ **超时控制**：支持上下文超时和超时参数

### 选择指南

| 需求 | 推荐使用 |
|------|----------|
| 需要进程控制（获取PID、Kill、Signal） | 主包 shellx |
| 需要 Windows cmd 命令支持 | 主包 shellx |
| 需要 PowerShell 命令支持 | 主包 shellx |
| 需要异步执行 | 主包 shellx |
| 纯 Go 实现，不依赖系统 shell | 子包 shx |
| 跨平台行为一致 | 子包 shx |
| 轻量级使用场景 | 子包 shx |

详细的 API 文档请参考：
- [📖 主包 API 文档](APIDOC.md)
- [📖 子包 shx API 文档](shx/APIDOC.md)

## ⚙️ 配置选项

### 主包 (shellx) 配置

#### 环境配置

```go
cmd := shellx.NewCmd("command").
    WithWorkDir("/custom/path").           // 设置工作目录
    WithEnv("KEY", "value").              // 添加环境变量
    WithTimeout(30 * time.Second).        // 设置超时时间
    WithContext(ctx)                      // 设置上下文
```

#### 输入输出配置

```go
var stdout, stderr bytes.Buffer
stdin := strings.NewReader("input")

cmd := shellx.NewCmd("command").
    WithStdin(stdin).                     // 设置标准输入
    WithStdout(&stdout).                  // 设置标准输出
    WithStderr(&stderr)                   // 设置标准错误
```

#### Shell 配置

```go
cmd := shellx.NewCmd("command").
    WithShell(shellx.ShellBash)           // 指定 shell 类型
```

### 子包 (shx) 配置

#### 环境配置

```go
cmd := shx.New("command").
    WithDir("/custom/path").              // 设置工作目录
    WithEnv("KEY", "value").              // 添加环境变量
    WithTimeout(30 * time.Second).        // 设置超时时间
    WithContext(ctx)                      // 设置上下文
```

#### 输入输出配置

```go
var stdout, stderr bytes.Buffer
stdin := strings.NewReader("input")

cmd := shx.New("command").
    WithStdin(stdin).                     // 设置标准输入
    WithStdout(&stdout).                  // 设置标准输出
    WithStderr(&stderr)                   // 设置标准错误
```

#### 批量环境变量

```go
// 方式1: 使用 map (WithEnvMap)
cmd := shx.New("command").
    WithEnvMap(map[string]string{
        "KEY1": "value1",
        "KEY2": "value2",
    })

// 方式2: 使用切片 (WithEnvs)
cmd := shx.New("command").
    WithEnvs([]string{
        "KEY1=value1",
        "KEY2=value2",
    })

// 方式3: 配合 os.Environ() 使用
envs := os.Environ()
envs = append(envs, "CUSTOM_VAR=my_value")
cmd := shx.New("command").WithEnvs(envs)
```

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
- 📦 [mvdan.cc/sh/v3](https://mvdan.cc/sh/v3) - shx 子包依赖

---

<div align="center">

**如果这个项目对您有帮助，请给它一个 ⭐ Star！**

[⬆️ 回到顶部](#shellx-)

</div>