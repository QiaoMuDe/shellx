# 支持 bash shell 脚本执行 — 分析与实施方案（简化版）

## 一、现状分析

### 1.1 shx 子包当前实现

当前 `shx` 子包基于 `mvdan.cc/sh/v3` 库，**仅支持执行命令字符串**，不支持执行 `.sh` 脚本文件。

**核心执行流程**（[exec.go](file:///d:/资源池/下水道/Dev/本地项目/shellx/shx/exec.go)）：

```
命令字符串 (s.raw)
    → bytes.NewReader([]byte(s.raw))          [io.Reader]
    → s.parser.Parse(reader, "")              [解析为 AST *syntax.File]
    → runner.Run(ctx, file)                   [执行 AST]
```

**关键依赖**：

* `mvdan.cc/sh/v3/syntax.Parser.Parse(io.Reader, string)` — 接受任何 `io.Reader`

* `mvdan.cc/sh/v3/interp.Runner.Run(ctx, *syntax.File)`

### 1.2 mvdan.cc/sh/v3 天然支持从文件读取

`syntax.Parser.Parse(io.Reader, string)` 天然支持从任何 `io.Reader` 读取，包括 `*os.File`。**无需修改库的使用方式**，只需在 shx 中增加从文件读取的分支。

***

## 二、设计方案（简化版）

### 2.1 核心原则

1. **默认写死 Bash** — 解析器固定使用 `syntax.NewParser(syntax.Variant(syntax.LangBash))`
2. **复用现有字段** — 新增最少字段，尽量复用已有机制
3. **简洁 API** — 只加 `NewScript` 构造函数 + 少量便捷函数

### 2.2 Shx 结构体改动

在 [types.go](file:///d:/资源池/下水道/Dev/本地项目/shellx/shx/types.go) 中仅新增一个字段：

```go
type Shx struct {
    // ... 现有字段不变 ...
    scriptFile string  // 脚本文件路径（空=执行命令字符串，非空=执行脚本文件）
}
```

复用 `raw` 字段存储脚本文件路径的显示名称，关键在于 `scriptFile` 字段的存在决定解析方式。

### 2.3 新增构造函数

在 [shx.go](file:///d:/资源池/下水道/Dev/本地项目/shellx/shx/shx.go) 中新增：

```go
// NewScript 从 bash 脚本文件创建命令
//
// 参数:
//   - filePath: 脚本文件路径
//
// 返回:
//   - *Shx: 命令对象
//
// 示例:
//
//	cmd := shx.NewScript("deploy.sh")
//	cmd := shx.NewScript("/path/to/script.sh")
func NewScript(filePath string) *Shx
```

同时将所有现有构造函数的解析器改为 `syntax.NewParser(syntax.Variant(syntax.LangBash))`。

### 2.4 执行逻辑改动

在 [exec.go](file:///d:/资源池/下水道/Dev/本地项目/shellx/shx/exec.go) 的 `execWithContext` 中增加分支：

```go
func (s *Shx) execWithContext(ctx context.Context) error {
    if strings.TrimSpace(s.raw) == "" && s.scriptFile == "" {
        return fmt.Errorf("command cannot be empty")
    }

    // 分支：从脚本文件解析 vs 从命令字符串解析
    var file *syntax.File
    var err error

    if s.scriptFile != "" {
        f, err := os.Open(s.scriptFile)
        if err != nil {
            return fmt.Errorf("open script file failed: %w", err)
        }
        defer f.Close()
        file, err = s.parser.Parse(f, s.scriptFile)
    } else {
        file, err = s.parser.Parse(bytes.NewReader([]byte(s.raw)), "")
    }

    if err != nil {
        return fmt.Errorf("parse error: %w", err)
    }

    // 构建 runner 并执行（不变）
    runner, err := s.buildRunner()
    if err != nil {
        return err
    }
    err = runner.Run(ctx, file)
    return handleError(err, s.displayName(), s.timeout)
}

// displayName 返回用于错误信息显示的标识
func (s *Shx) displayName() string {
    if s.scriptFile != "" {
        return s.scriptFile
    }
    return s.raw
}
```

### 2.5 新增便捷函数

在 [funcs.go](file:///d:/资源池/下水道/Dev/本地项目/shellx/shx/funcs.go) 中新增：

```go
// RunScript 执行 bash 脚本文件
//
// 参数:
//   - filePath: 脚本文件路径
//
// 返回:
//   - error: 执行错误
//
// 示例:
//
//	err := shx.RunScript("deploy.sh")
func RunScript(filePath string) error {
    return NewScript(filePath).Exec()
}

// OutScript 执行 bash 脚本文件并获取输出
//
// 参数:
//   - filePath: 脚本文件路径
//
// 返回:
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例:
//
//	output, err := shx.OutScript("deploy.sh")
func OutScript(filePath string) ([]byte, error) {
    return NewScript(filePath).ExecOutput()
}
```

***

## 三、改动文件清单

| 文件                                                                | 改动内容                                                                                     |
| ----------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| [shx/types.go](file:///d:/资源池/下水道/Dev/本地项目/shellx/shx/types.go)   | 新增 `scriptFile string` 字段                                                                |
| [shx/shx.go](file:///d:/资源池/下水道/Dev/本地项目/shellx/shx/shx.go)       | 新增 `NewScript` 构造函数；现有构造函数 parser 改为 `syntax.NewParser(syntax.Variant(syntax.LangBash))` |
| [shx/exec.go](file:///d:/资源池/下水道/Dev/本地项目/shellx/shx/exec.go)     | `execWithContext` 增加文件解析分支；新增 `displayName()` 辅助方法                                       |
| [shx/funcs.go](file:///d:/资源池/下水道/Dev/本地项目/shellx/shx/funcs.go)   | 新增 `RunScript`、`OutScript` 便捷函数                                                          |
| [shx/errors.go](file:///d:/资源池/下水道/Dev/本地项目/shellx/shx/errors.go) | 可选：`handleError` 中对文件路径错误做更好的提示                                                          |

**不修改的文件**：

* `option.go` — 无需新增配置方法

* `APIDOC.md` — 如需文档更新后面再补充

***

## 四、使用示例

```go
// 1. 执行脚本文件
err := shx.RunScript("deploy.sh")

// 2. 执行脚本并获取输出
output, err := shx.OutScript("deploy.sh")

// 3. 链式配置（复用现有 WithXxx 方法）
err := shx.NewScript("build.sh").
    WithDir("/project").
    WithEnv("MODE", "production").
    WithTimeout(60 * time.Second).
    Exec()

// 4. 获取输出 + 链式配置
output, err := shx.NewScript("test.sh").
    WithTimeout(30 * time.Second).
    ExecOutput()
```

***

## 五、验证步骤

1. 运行现有测试套件确保无回归：`go test ./shx/...`
2. 新增 `script_test.go`，覆盖：

   * 正常脚本文件执行

   * 不存在的脚本文件（返回错误）

   * 空文件路径（返回错误）

   * 脚本 + 超时场景

   * 脚本 + 链式配置

***

## 六、实施顺序

1. **Step 1**: `types.go` — 加 `scriptFile` 字段
2. **Step 2**: `shx.go` — 加 `NewScript`，改所有构造函数的 parser 为 Bash 方言
3. **Step 3**: `exec.go` — 加文件分支 + `displayName()`
4. **Step 4**: `funcs.go` — 加 `RunScript`、`OutScript`
5. **Step 5**: 新增 `script_test.go` 并运行测试验证

