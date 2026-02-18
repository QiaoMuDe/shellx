package shellx

import (
	"strings"
	"unicode"
)

// shell特殊字符集合
//
// 包含常见的shell特殊字符，用于命令分隔、管道、重定向等操作
var specialChars = map[string]bool{
	";": true, // 命令分隔符
	"|": true, // 管道符
	">": true, // 重定向符
	"<": true, // 重定向符
	"&": true, // 后台运行/逻辑运算
	"!": true, // 逻辑非
	"`": true, // 命令替换
	"*": true, // 通配符
	"?": true, // 通配符
	"#": true, // 注释符
	"~": true, // 家目录
}

// splitState 命令拆分过程中的状态信息
//
// 封装了命令拆分过程中需要的所有状态变量，
// 便于状态管理和参数传递
type splitState struct {
	result         []string        // 拆分结果
	builder        strings.Builder // 当前命令片段构建器
	inQuotes       bool            // 是否在引号中
	quote          string          // 当前引号类型
	hasQuoteInWord bool            // 当前片段是否包含过引号(用于处理空引号情况)
	emptyQuote     bool            // 当前引号是否为空(用于区分空引号和非空引号）
}

// newSplitState 创建新的拆分状态
//
// 返回值:
//   - *splitState: 初始化后的拆分状态指针
func newSplitState() *splitState {
	return &splitState{
		result:         make([]string, 0, 8),
		builder:        strings.Builder{},
		inQuotes:       false,
		quote:          "",
		hasQuoteInWord: false,
		emptyQuote:     false,
	}
}

// isQuote 判断字符是否为引号（单引号、双引号、反引号）
//
// 支持的引号类型:
//   - 双引号: "
//   - 单引号: '
//   - 反引号: `
//
// 参数:
//   - ch: 要判断的字符
//
// 返回值:
//   - bool: 如果是引号返回 true，否则返回 false
func isQuote(ch string) bool {
	switch ch {
	case "\"", "'", "`":
		return true
	default:
		return false
	}
}

// isSpecialChar 判断字符是否为shell特殊字符
//
// 参数:
//   - ch: 要判断的字符
//
// 返回值:
//   - bool: 如果是特殊字符返回 true，否则返回 false
func isSpecialChar(ch string) bool {
	return specialChars[ch]
}

// handleSpecialChar 处理特殊字符（分号、管道符等）的逻辑
//
// 在非引号状态下，特殊字符作为独立token处理
// 在引号内，特殊字符作为普通字符处理
//
// 参数:
//   - state: 拆分状态
//   - ch: 当前字符
func (s *splitState) handleSpecialChar(ch string) {
	// 根据是否在引号内采用不同的处理策略
	if s.inQuotes {
		// 引号内：特殊字符作为普通字符处理
		s.builder.WriteString(ch)
	} else {
		// 引号外：特殊字符作为独立token处理

		// 1. 先处理当前累积的token
		if s.builder.Len() > 0 {
			s.result = append(s.result, s.builder.String())
			s.builder.Reset()
		}

		// 2. 将特殊字符作为独立token添加到结果中
		s.result = append(s.result, ch)
	}

	// 3. 重置空引号状态标记
	s.emptyQuote = false
}

// handleEscapeChar 处理转义字符的逻辑
//
// 将转义字符和下一个字符作为一个整体处理
// 对于某些特殊字符（如;），保持转义形式
//
// 参数:
//   - state: 拆分状态
//   - cmdStr: 输入字符串
//   - i: 当前位置
//
// 返回值:
//   - int: 新的位置索引
func (s *splitState) handleEscapeChar(cmdStr string, i int) int {
	// 先处理当前累积的token
	if s.builder.Len() > 0 {
		s.result = append(s.result, s.builder.String())
		s.builder.Reset()
	}

	// 获取被转义的字符
	escapedChar := cmdStr[i+1]

	// 对于某些特殊字符，保持转义形式
	switch escapedChar {
	case ';': // 分号在find -exec等命令中需要转义
		s.result = append(s.result, "\\;")
	default:
		// 其他情况，只保留被转义的字符，去掉反斜杠
		s.result = append(s.result, string(escapedChar))
	}

	s.emptyQuote = false
	return i + 1 // 返回跳过后的位置
}

// checkMultiCharOperator 检查并处理多字符操作符
//
// 参数:
//   - state: 拆分状态
//   - cmdStr: 输入字符串
//   - i: 当前位置
//
// 返回值:
//   - bool: 如果是多字符操作符返回 true，否则返回 false
//   - int: 新的位置索引（如果是多字符操作符）
func checkMultiCharOperator(state *splitState, cmdStr string, i int) (bool, int) {
	// 检查多字符操作符
	if i+1 >= len(cmdStr) {
		return false, i
	}

	twoChar := cmdStr[i : i+2]
	switch twoChar {
	case "&&", "||", ">>", "<<":
		// 先处理当前累积的token
		if state.builder.Len() > 0 {
			state.result = append(state.result, state.builder.String())
			state.builder.Reset()
		}
		// 将多字符操作符作为独立token
		state.result = append(state.result, twoChar)
		// 跳过多字符操作符的长度
		state.emptyQuote = false
		return true, i + 1 // 跳过下一个字符

	default:
		return false, i
	}
}

// splitInternal 将命令字符串拆分为命令切片（内部函数）
//
// 实现原理：
//  1. 去除首尾空白
//  2. 遍历每个字符
//  3. 处理引号状态切换
//  4. 在非引号状态下遇到空格时分割
//  5. 检测未闭合的引号
//
// 参数:
//   - cmdStr: 要拆分的命令字符串
//
// 返回值:
//   - []string: 拆分后的命令切片
//   - error: 拆分错误，成功时为 nil
func splitInternal(cmdStr string) ([]string, error) {
	// 去除首尾空白
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return []string{}, nil
	}

	// 初始化拆分状态
	state := newSplitState()

	// 遍历每个字符
	for i := 0; i < len(cmdStr); i++ {
		// 处理转义字符
		if cmdStr[i] == '\\' && i+1 < len(cmdStr) && !state.inQuotes {
			i = state.handleEscapeChar(cmdStr, i)
			continue
		}

		// 获取当前字符
		currentChar := cmdStr[i : i+1]

		// 检查多字符操作符
		if isMultiOp, newPos := checkMultiCharOperator(state, cmdStr, i); isMultiOp {
			i = newPos
			continue
		}

		switch {
		case isQuote(currentChar):
			// 处理引号字符
			state.handleQuoteChar(currentChar)

		case isSpecialChar(currentChar):
			// 处理特殊字符（分号、管道符等）
			state.handleSpecialChar(currentChar)

		case unicode.IsSpace(rune(currentChar[0])) && !state.inQuotes:
			// 统一处理所有空白字符分隔符（包括换行符）
			state.handleSeparator()

		default:
			// 处理普通字符
			state.builder.WriteString(currentChar)
			state.emptyQuote = false
		}
	}

	// 添加最后一个命令片段
	if state.builder.Len() > 0 || state.hasQuoteInWord {
		state.result = append(state.result, state.builder.String())
	}

	// 检查引号是否闭合，如果未闭合返回带有未闭合引号类型的错误
	if state.inQuotes {
		return state.result, &UnclosedQuoteError{QuoteType: state.quote}
	}

	return state.result, nil
}

// handleQuoteChar 处理引号字符的逻辑
//
// 参数:
//   - ch: 当前字符
func (s *splitState) handleQuoteChar(ch string) {
	switch {
	case !s.inQuotes:
		// 开始引号
		s.inQuotes = true   // 进入引号状态
		s.quote = ch        // 记录当前引号类型
		s.emptyQuote = true // 假设为空引号，直到遇到非引号字符

	case ch == s.quote:
		// 结束引号
		s.inQuotes = false
		// 处理空引号逻辑：非空引号或独立空引号才设置 hasQuoteInWord
		if !s.emptyQuote || s.builder.Len() == 0 {
			s.hasQuoteInWord = true
		}

	default:
		// 引号内的其他引号字符
		s.builder.WriteString(ch) // 直接添加引号字符
		s.emptyQuote = false      // 有内容，不是空引号
	}
}

// handleSeparator 处理分隔符（空格或制表符）的逻辑
//
// 判断逻辑:
//   - builder.Len() > 0：片段中有内容（非空）
//   - 只在有内容时添加token，避免空字符串token
func (s *splitState) handleSeparator() {
	// 判断是否需要添加当前命令片段
	if s.builder.Len() > 0 {
		s.result = append(s.result, s.builder.String()) // 添加当前命令片段
		s.builder.Reset()                               // 重置构建器，准备下一个片段
	}
	s.hasQuoteInWord = false // 重置标记，为下一个片段做准备
}
