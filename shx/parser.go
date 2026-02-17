package shx

import (
	"bytes"
	"strings"
	"sync"

	"mvdan.cc/sh/v3/syntax"
)

// 全局 parser 和 printer 实例
var (
	globalParser *syntax.Parser  // 用于解析命令字符串
	printer      *syntax.Printer // 用于打印解析后的命令
	once         sync.Once       // 用于确保 parser 和 printer 只创建一次
)

// getParserAndPrinter 获取全局 parser 和 printer 实例（线程安全）
//
// 使用 sync.Once 确保 parser 和 printer 只创建一次
// Parser.Words 方法每次调用都会 reset，因此是线程安全的
func getParserAndPrinter() (*syntax.Parser, *syntax.Printer) {
	once.Do(func() {
		globalParser = syntax.NewParser()
		printer = syntax.NewPrinter()
	})
	return globalParser, printer
}

// parseCmdInternal 将命令字符串解析为命令切片（内部函数）
//
// 使用 mvdan.cc/sh/v3 的 Words 方法解析，支持完整的 shell 语法：
//   - 环境变量：${VAR}, $VAR
//   - 通配符：*.go, test?.txt
//   - 命令替换：$(cmd), `cmd`
//   - 转义字符：\", \', \`
//   - 引号嵌套：支持不同类型引号嵌套
//
// 参数:
//   - cmdStr: 要解析的命令字符串
//
// 返回值:
//   - []string: 解析后的命令切片
//   - error: 解析错误，成功时为 nil
func parseCmdInternal(cmdStr string) ([]string, error) {
	parser, printer := getParserAndPrinter()

	var builder strings.Builder // 局部字符串构建器，避免并发问题

	// 解析命令字符串，将每个单词添加到结果切片
	result := make([]string, 0, 8)

	err := parser.Words(bytes.NewReader([]byte(cmdStr)), func(word *syntax.Word) bool {
		builder.Reset()               // 重置字符串构建器
		printer.Print(&builder, word) // 打印单词到构建器
		str := builder.String()       // 获取构建器中的字符串

		// 如果字符串不为空，则添加到结果切片
		if str != "" {
			result = append(result, str)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
