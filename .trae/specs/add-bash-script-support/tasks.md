# Tasks

- [x] Task 1: 修改 types.go — 在 Shx 结构体中新增 `scriptFile string` 字段
- [x] Task 2: 修改 shx.go — 所有构造函数 parser 改为 `syntax.NewParser(syntax.Variant(syntax.LangBash))`，新增 `NewScript` 构造函数
- [x] Task 3: 修改 exec.go — `execWithContext` 增加文件解析分支，新增 `displayName()` 辅助方法
- [x] Task 4: 修改 funcs.go — 新增 `RunScript`、`OutScript` 便捷函数
- [x] Task 5: 新增 script_test.go — 添加脚本执行相关测试用例
- [x] Task 6: 运行测试验证 — `go test ./shx/...` 确保现有测试和新增测试全部通过

# Task Dependencies

- Task 1 是 Task 2、3 的基础
- Task 2、3 可并行实现（修改不同文件）
- Task 4 依赖 Task 2（NewScript 构造函数）
- Task 5 依赖 Task 2、3、4
- Task 6 依赖 Task 5
