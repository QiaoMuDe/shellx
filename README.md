<div align="center">

# ShellX

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org) [![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](https://gitee.com/MM-Q/shellx/blob/master/LICENSE) [![Gitee](https://img.shields.io/badge/Gitee-Repository-red?style=for-the-badge&logo=gitee)](https://gitee.com/MM-Q/shellx)

**一个功能完善、易于使用的 Go 语言 Shell 命令执行库**

[🏠 仓库地址](https://gitee.com/MM-Q/shellx) • [📖 子包 API 文档](shx/APIDOC.md)

</div>

---

## 项目简介

ShellX 是一个功能完善、易于使用的 Go 语言 Shell 命令执行库，包含两个子包：

- **主包 (shellx)** — 基于 Go 标准库 `os/exec` 封装，提供更友好的 API 和丰富的功能
- **子包 (shx)** — 基于 [mvdan.cc/sh/v3](https://mvdan.cc/sh/v3) 纯 Go 实现，默认使用 Bash 方言解析器，支持执行 `.sh` 脚本文件

两者 API 风格一致，按需选择即可。

---

## 安装

```bash
# 主包
go get gitee.com/MM-Q/shellx

# 子包 shx
go get gitee.com/MM-Q/shellx/shx
```

Go 版本要求：1.25.0 或更高

---

## 核心特性

### 主包 (shellx)

| 特性 | 说明 |
|------|------|
| **多种创建方式** | `NewCmd`(可变参数)、`NewCmds`(切片)、`NewCmdStr`(字符串解析) |
| **链式调用** | 流畅的方法链，支持工作目录、环境变量、超时、上下文等配置 |
| **多 Shell 支持** | sh、bash、cmd、powershell、pwsh 及系统默认 Shell |
| **进程控制** | 获取 PID、发送信号、Kill 进程 |
| **异步执行** | 非阻塞启动，可后续等待完成 |
| **丰富便捷函数** | `Exec`/`ExecStr`/`ExecOut`/`ExecOutStr` 及其带超时版本，开箱即用 |
| **精确超时控制** | 延迟构建 exec.Cmd，确保超时计时精确 |
| **并发安全** | 线程安全设计，支持多 goroutine 环境 |

### 子包 (shx)

| 特性 | 说明 |
|------|------|
| **纯 Go 实现** | 基于 mvdan.cc/sh/v3，不依赖系统 shell |
| **Bash 方言默认** | 默认使用 Bash 方言解析器，支持 `[[ ]]`、`function`、`select` 等 Bash 特有语法 |
| **脚本文件执行** | 原生支持执行 `.sh` 脚本文件 |
| **链式调用** | 支持流畅的方法链配置 |
| **超时控制** | 支持上下文超时和超时参数 |
| **输入输出重定向** | 灵活的标准输入输出配置 |
| **轻量级并发** | 使用 atomic.Bool 防止重复执行 |
| **退出码检测** | `IsExitStatus` 可提取命令退出码 |

---

## 使用示例

### 主包 (shellx)

```go
package main

import (
    "fmt"
    "log"
    "time"

    "gitee.com/MM-Q/shellx"
)

func main() {
    // ---- 便捷函数 ----
    shellx.Exec("echo", "hello")                // 简单执行
    shellx.ExecStr("ls -la")                     // 字符串方式
    output, _ := shellx.ExecOut("pwd")           // 获取输出
    fmt.Printf("Current: %s", output)

    shellx.ExecT(5*time.Second, "sleep", "10")   // 5秒超时

    // ---- 链式配置 ----
    err := shellx.NewCmd("echo", "Hello, World!").
        WithTimeout(10 * time.Second).
        WithWorkDir("/tmp").
        WithEnv("KEY", "value").
        Exec()
    if err != nil {
        log.Fatal(err)
    }

    // ---- 获取输出 ----
    out, err := shellx.NewCmd("ls", "-la").ExecOutput()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(out))

    // ---- 进程控制 ----
    cmd := shellx.NewCmd("sleep", "10")
    err = cmd.ExecAsync()
    if err != nil {
        log.Fatal(err)
    }
    pid := cmd.GetPID()
    fmt.Printf("PID: %d\n", pid)
    cmd.Wait()

    // ---- 选择 Shell 类型 ----
    shellx.NewCmdStr("echo hello").WithShell(shellx.ShellBash).Exec()
}
```

### 子包 (shx)

```go
package main

import (
    "fmt"
    "log"
    "time"

    "gitee.com/MM-Q/shellx/shx"
)

func main() {
    // ---- 便捷函数 ----
    shx.Run("echo hello")                        // 简单执行
    shx.RunToTerminal("ls -la")                  // 输出到终端

    output, _ := shx.Out("date")                 // 获取输出
    fmt.Println(string(output))

    shx.RunWith("sleep 10", 5*time.Second)       // 超时执行

    // ---- 链式配置 ----
    out, err := shx.New("echo hello").
        WithTimeout(5 * time.Second).
        WithDir("/tmp").
        WithEnv("FOO", "bar").
        ExecOutput()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(out))

    // ---- 使用上下文 ----
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err = shx.New("long-command").WithContext(ctx).Exec()

    // ---- 自定义输入输出 ----
    var stdout, stderr bytes.Buffer
    stdin := strings.NewReader("input")
    err = shx.New("cat").
        WithStdin(stdin).
        WithStdout(&stdout).
        WithStderr(&stderr).
        Exec()

    // ---- 执行脚本文件 ----
    shx.RunScript("deploy.sh")
    output, err = shx.OutScript("build.sh")

    out, err = shx.NewScript("test.sh").
        WithDir("/project").
        WithEnv("MODE", "ci").
        WithTimeout(30 * time.Second).
        ExecOutput()

    // ---- 检查退出码 ----
    err = shx.Run("exit 5")
    if code, ok := shx.IsExitStatus(err); ok {
        fmt.Printf("Exit code: %d\n", code)
    }
}
```

---

## Shell 类型说明

### 主包 — 支持 8 种 Shell 类型

| 类型 | 说明 |
|------|------|
| `ShellSh` | 标准 Unix shell |
| `ShellBash` | Bash shell |
| `ShellCmd` | Windows 命令提示符 |
| `ShellPowerShell` | Windows PowerShell |
| `ShellPwsh` | PowerShell Core（跨平台） |
| `ShellNone` | 直接执行，不使用 shell |
| `ShellDef1` | 自动选择：Windows → cmd，其他 → sh |
| `ShellDef2` | 自动选择：Windows → powershell，其他 → sh |

### 子包 — 使用 mvdan.cc/sh/v3 解析器

子包 shx 默认使用 **Bash 方言解析器**（`syntax.LangBash`），支持 `[[ ]]`、`function`、`select` 等 Bash 特有语法。

---

## 选择指南

| 需求 | 推荐包 |
|------|--------|
| 需要进程控制（获取 PID、Kill、Signal） | 主包 shellx |
| 需要 Windows cmd、PowerShell 命令支持 | 主包 shellx |
| 需要异步执行 | 主包 shellx |
| 纯 Go 实现，不依赖系统 shell | 子包 shx |
| 跨平台行为一致 | 子包 shx |
| 执行 `.sh` 脚本文件 | 子包 shx |
| 轻量级使用场景 | 子包 shx |

---

## 命令解析特性

两个包均支持：

- ✅ 单引号、双引号、反引号支持
- ✅ 引号嵌套处理
- ✅ 转义字符支持
- ✅ 多空格和制表符处理
- ✅ 未闭合引号检测
- ✅ Unicode 字符支持（中文、emoji 等多字节字符）

---

<div align="center">

**如果这个项目对您有帮助，请给它一个 ⭐ Star！**

[⬆顶部](#shellx)

</div>