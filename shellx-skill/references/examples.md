# 使用示例

## 1. 基础命令执行

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/shellx"
)

func main() {
    shellx.Exec("echo", "hello world")
    output, _ := shellx.ExecOut("ls", "-la")
    fmt.Println(string(output))
}
```

## 2. Shell 管道（shx）

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/shellx/shx"
)

func main() {
    output, _ := shx.Out("ls -la | grep .go | wc -l")
    fmt.Printf("Go files: %s", output)
}
```

## 3. 超时控制

```go
package main

import (
    "fmt"
    "time"
    "gitee.com/MM-Q/shellx"
)

func main() {
    err := shellx.ExecT(3*time.Second, "ping", "google.com")
    if err != nil {
        fmt.Println("Timeout or failed:", err)
    }
}
```

## 4. 上下文取消

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

err := shellx.NewCmd("sleep", "10").
    WithContext(ctx).
    WithStdout(os.Stdout).
    Exec()
```

## 5. 异步执行与进程控制

```go
cmd := shellx.NewCmd("sleep", "30")
cmd.ExecAsync()
fmt.Printf("PID: %d\n", cmd.GetPID())
time.Sleep(1 * time.Second)
cmd.Kill()
err := cmd.Wait()
```

## 6. 自定义 IO（shx）

```go
var stdout, stderr bytes.Buffer
stdin := strings.NewReader("hello\nworld\n")
shx.RunWithIO("cat", stdin, &stdout, &stderr)
```

## 7. Windows 路径

```go
// 反斜杠不会被误当作转义
shellx.NewCmdStr(`C:\Windows\System32\notepad.exe C:\temp\file.txt`)
```

## 8. Shell 语法（shx）

```go
output, _ := shx.Out(`
    for i in 1 2 3; do
        echo "number $i"
    done
`)
```

## 9. 环境变量

```go
output, _ := shellx.NewCmd("sh", "-c", "echo $VAR").
    WithEnv("VAR", "hello").
    ExecOutput()
```

## 10. 命令查找

```go
path, err := shellx.FindCmd("go")    // 增强版
path := shellx.FindCommandPath("go") // 便捷版
if path == "" { /* not found */ }
```

## 11. Shell 类型选择

```go
// PowerShell
shellx.NewCmdStr(`Get-ChildItem -Path "C:\"`).
    WithShell(shellx.ShellPowerShell).Exec()

// Bash
shellx.NewCmdStr(`find / -name "*.go"`).
    WithShell(shellx.ShellBash).Exec()
```

## 12. 错误处理

```go
err := cmd.Exec()
if err != nil {
    // 超时 / 取消 / 未找到 / 退出码
}
if code, ok := shx.IsExitStatus(err); ok {
    fmt.Printf("exit code: %d\n", code)
}
```

## 13. 执行脚本文件（shx）

```go
package main

import (
    "fmt"
    "time"
    "gitee.com/MM-Q/shellx/shx"
)

func main() {
    // 简单执行
    err := shx.RunScript("deploy.sh")
    if err != nil {
        fmt.Println("Script failed:", err)
    }

    // 带超时
    output, err := shx.OutScriptWith("build.sh", 60*time.Second)
    if err != nil {
        fmt.Println("Build timed out or failed:", err)
    }
    fmt.Println(string(output))

    // 链式配置
    output, err = shx.NewScript("test.sh").
        WithDir("/project").
        WithEnv("MODE", "ci").
        WithTimeout(30 * time.Second).
        ExecOutput()
}
```
