---
name: shellx-skill
description: 为 Go 语言的 ShellX 库（gitee.com/MM-Q/shellx）生成 Shell 命令执行代码。当用户需要在 Go 中执行系统命令、调用子进程、获取命令输出/退出码、设置命令超时/取消、跨平台兼容执行 Shell 脚本、执行 .sh 脚本文件、使用管道/重定向等 Shell 特性时使用此技能。覆盖主包 shellx（基于 os/exec，支持进程控制和异步执行）和子包 shx（基于 mvdan.cc/sh/v3 的纯 Go 实现，跨平台一致，支持 Bash 方言和脚本文件执行）。对于用户涉及"执行命令""运行程序""调用子进程""获取输出""设置超时"等表述时应主动使用此技能。
---

# ShellX Skill

为 [ShellX](https://gitee.com/MM-Q/shellx) 库生成 Go 代码。

## 快速选择

| 需求 | 推荐包 | 原因 |
|------|--------|------|
| 进程 PID/Kill/Signal | **shellx**（主包） | 子包不支持进程控制 |
| 异步执行 / 非阻塞 | **shellx**（主包） | 子包仅同步 |
| Windows cmd 命令 | **shellx**（主包） | 支持 ShellCmd 类型 |
| 跨平台一致性 | **shx**（子包） | pure Go，不依赖系统 shell |
| 管道/重定向 | **shx**（子包） | mvdan.cc/sh 完整支持 |
| 执行 `.sh` 脚本文件 | **shx**（子包） | 通过 NewScript/RunScript 原生支持 |
| 零外部依赖要求 | **shellx**（主包） | 仅依赖 Go 标准库 |

## 主包 `shellx`

基于 Go 标准库 `os/exec`，提供完整的命令生命周期管理。

### 创建命令

```go
import "gitee.com/MM-Q/shellx"

// 三种创建方式
cmd := shellx.NewCmd("ls", "-la", "/tmp")                  // 可变参数
cmd := shellx.NewCmds([]string{"git", "status"})           // 切片
cmd := shellx.NewCmdStr(`git commit -m "feat: x"`)        // 字符串（自动分词）
```

### 链式配置

```go
cmd := shellx.NewCmd("myapp").
    WithWorkDir("/app").                   // 工作目录（验证存在性）
    WithTimeout(30 * time.Second).         // 超时控制
    WithEnv("KEY", "value").              // 环境变量
    WithEnvs([]string{"A=1", "B=2"}).     // 批量环境变量
    WithShell(shellx.ShellBash).          // Shell 类型
    WithStdin(os.Stdin).                  // 标准输入
    WithStdout(os.Stdout).                // 标准输出
    WithStderr(os.Stderr).                // 标准错误
    WithContext(ctx)                      // 上下文（覆盖超时）
```

### 执行方式

```go
// 同步执行
err := cmd.Exec()
output, err := cmd.ExecOutput()      // 合并输出
stdout, err := cmd.ExecStdout()      // 仅标准输出

// 异步执行
cmd.ExecAsync()
pid := cmd.GetPID()
cmd.IsRunning()
cmd.Kill()
cmd.Signal(syscall.SIGTERM)
err = cmd.Wait()
code, err := cmd.WaitWithCode()
```

### 便捷函数

```go
// 基础
shellx.Exec("ls", "-la")
shellx.ExecStr("echo hello")
shellx.ExecOut("ls", "-la")
shellx.ExecOutStr("date")

// 带超时
shellx.ExecT(5*time.Second, "sleep", "10")
shellx.ExecStrT(3*time.Second, "ping google.com")
shellx.ExecOutT(2*time.Second, "curl", "example.com")

// 退出码
code, _ := shellx.ExecCode("git", "status")
code, _ := shellx.ExecCodeStr("docker ps")
```

### 命令查找

```go
path, err := shellx.FindCmd("go")              // 增强版
if err != nil { /* 未找到 */ }

path := shellx.FindCommandPath("go")           // 便捷版（找不到返回 ""）
```

### 命令字符串拆分

```go
// Selective escaping: `\` only escapes before quote/special/space/`\`
// Windows paths like `.\path\file` are handled correctly
shellx.Split(`git commit -m "fix: bug"`)
// ["git", "commit", "-m", "fix: bug"]

shellx.SplitE(`echo "unclosed quote`)
// 返回 UnclosedQuoteError
```

### 上下文与超时优先级

```
WithContext(ctx)  →  最高优先级，完全覆盖 WithTimeout
WithTimeout(d)    →  次优先级，ctx 为 nil 时生效
默认              →  context.Background()（无超时，无取消）
```

### 并发安全说明

- 配置方法（`WithXxx`）**不是**并发安全的，不要多个 goroutine 同时配置
- 执行方法（`Exec`, `ExecAsync` 等）是并发安全的（`atomic.Bool` 保护）
- 每个 `Command` 只能执行一次，重复执行返回 `ErrAlreadyExecuted`
- 属性获取方法（`Name()`, `Args()`等）不是并发安全的

## 子包 `shx`

基于 `mvdan.cc/sh/v3`，纯 Go 实现，跨平台行为一致。默认使用 Bash 方言解析器（`syntax.LangBash`），支持 `[[ ]]`、`function`、`select` 等 Bash 特有语法。也支持直接执行 `.sh` 脚本文件。

### 创建与配置

```go
import "gitee.com/MM-Q/shellx/shx"

// 命令字符串
cmd := shx.New("echo $VAR").
    WithDir("/tmp").
    WithTimeout(5 * time.Second).
    WithEnv("VAR", "value").
    WithEnvMap(map[string]string{"A": "1"}).
    WithStdin(reader).
    WithStdout(&buf).
    WithStderr(&buf).
    WithContext(ctx)

// 脚本文件
cmd := shx.NewScript("deploy.sh").
    WithDir("/app").
    WithTimeout(30 * time.Second).
    WithEnv("MODE", "production")
```

### 执行

```go
err := cmd.Exec()
output, err := cmd.ExecOutput()
err := cmd.ExecContext(ctx)
output, err := cmd.ExecContextOutput(ctx)
```

### 便捷函数

```go
// 命令字符串
shx.Run("echo hello")
shx.RunToTerminal("ls -la")
shx.Out("date")
shx.RunWith("sleep 10", 5*time.Second)
shx.OutWith("sleep 5", 10*time.Second)
shx.RunWithIO("cat", stdin, stdout, stderr)
shx.OutWithIO("cat", stdin, stdout, stderr)
shx.RunCtx(ctx, "echo")
shx.OutCtx(ctx, "ls")

// 脚本文件
shx.RunScript("deploy.sh")
shx.RunScriptToTerminal("build.sh")
shx.OutScript("deploy.sh")
shx.RunScriptWith("long_task.sh", 30*time.Second)
shx.OutScriptWith("build.sh", 60*time.Second)
shx.RunScriptWithIO("script.sh", stdin, stdout, stderr)
shx.OutScriptWithIO("script.sh", stdin, stdout, stderr)
shx.RunCtxScript(ctx, "long_task.sh")
shx.OutCtxScript(ctx, "build.sh")
```

### 退出码判断

```go
if code, ok := shx.IsExitStatus(err); ok {
    fmt.Printf("exit code: %d\n", code)
}
```

## 参考文档

- [主包 API 参考](references/main_package.md) — shellx 完整 API
- [子包 API 参考](references/subpackage.md) — shx 完整 API
- [使用示例](references/examples.md) — 12 个场景的完整代码
