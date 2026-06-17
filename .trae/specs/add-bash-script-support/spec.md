# 添加 Bash 脚本执行支持 Spec

## Why

当前 `shx` 子包基于 `mvdan.cc/sh/v3`，但解析器默认使用 POSIX Shell 方言，且仅支持执行命令字符串，无法执行 `.sh` 脚本文件。需要将默认方言改为 Bash，并支持直接执行 bash 脚本文件。

## What Changes

1. **默认解析器改为 Bash 方言**: 所有现有构造函数的 `syntax.NewParser()` 改为 `syntax.NewParser(syntax.Variant(syntax.LangBash))`
2. **脚本文件执行支持**: `Shx` 结构体新增 `scriptFile` 字段，`NewScript` 构造函数，`RunScript`/`OutScript` 便捷函数
3. **执行引擎分支**: `execWithContext` 根据 `scriptFile` 是否为空选择从文件或字符串解析

## Impact

- Affected specs: shx 子包 API
- Affected code: `shx/types.go`, `shx/shx.go`, `shx/exec.go`, `shx/funcs.go`
- No breaking changes to existing API

## ADDED Requirements

### Requirement: Bash 方言全局默认

所有构造函数创建的解析器 SHALL 使用 `syntax.LangBash` 方言。

#### Scenario: 现有构造函数使用 Bash
- **WHEN** 调用 `New("echo hello")`、`NewArgs("ls", "-la")`、`NewCmds([]string{"echo", "hello"})`
- **THEN** 内部 `parser` 使用 `syntax.NewParser(syntax.Variant(syntax.LangBash))`

### Requirement: 脚本文件执行

系统 SHALL 提供执行 bash 脚本文件的能力。

#### Scenario: NewScript 构造函数
- **WHEN** 调用 `shx.NewScript("script.sh")`
- **THEN** 返回 `*Shx` 对象，`scriptFile` 字段为 `"script.sh"`，可通过链式调用配置并执行

#### Scenario: RunScript 便捷函数
- **WHEN** 调用 `shx.RunScript("script.sh")`
- **THEN** 执行 `script.sh` 脚本，返回 `error`

#### Scenario: OutScript 便捷函数
- **WHEN** 调用 `shx.OutScript("script.sh")`
- **THEN** 执行 `script.sh` 脚本，返回 `([]byte, error)`

#### Scenario: 脚本文件不存在
- **WHEN** 调用 `NewScript("nonexistent.sh").Exec()`
- **THEN** 返回文件打开失败的错误

#### Scenario: 空文件路径
- **WHEN** 调用 `NewScript("").Exec()`
- **THEN** 返回空命令错误

### Requirement: 执行引擎分支

`execWithContext` SHALL 根据 `scriptFile` 字段决定解析来源。

#### Scenario: scriptFile 非空
- **WHEN** `s.scriptFile != ""`
- **THEN** 使用 `os.Open(s.scriptFile)` 获取 `io.Reader`，传给 `parser.Parse()`

#### Scenario: scriptFile 为空
- **WHEN** `s.scriptFile == ""`
- **THEN** 保持现有逻辑，从 `s.raw` 字符串解析

### Requirement: 错误信息优化

脚本执行相关的错误信息 SHALL 使用脚本文件路径作为标识。

#### Scenario: 错误信息显示
- **WHEN** 脚本执行出错
- **THEN** 错误信息中包含脚本文件路径（而非空字符串）

### Requirement: 链式配置兼容

`NewScript` 创建的对象 SHALL 支持现有的所有 `WithXxx` 链式配置方法（`WithDir`、`WithEnv`、`WithTimeout` 等）。

#### Scenario: 链式配置脚本执行
- **WHEN** 调用 `shx.NewScript("build.sh").WithDir("/project").WithEnv("K", "V").WithTimeout(30*time.Second).Exec()`
- **THEN** 脚本在指定目录、指定环境变量、指定超时下执行
