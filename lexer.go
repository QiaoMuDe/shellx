package shellx

import (
	"strings"
	"unicode"
)

// splitState 命令拆分过程中的状态信息
//
// 封装了命令拆分过程中需要的所有状态变量，
// 便于状态管理和参数传递
type splitState struct {
	result         []string        // 拆分结果
	builder        strings.Builder // 当前命令片段构建器
	inQuotes       bool            // 是否在引号中
	quote          rune            // 当前引号类型
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
		quote:          0,
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
//   - r: 要判断的字符
//
// 返回值:
//   - bool: 如果是引号返回 true，否则返回 false
func isQuote(r rune) bool {
	return strings.ContainsRune("\"'`", r)
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
	for _, r := range cmdStr {
		switch {
		case isQuote(r):
			// 处理引号字符
			state.handleQuoteChar(r)

		case unicode.IsSpace(r) && !state.inQuotes:
			// 特殊处理换行符：直接添加到当前片段，不分割
			if r == '\n' || r == '\r' {
				state.builder.WriteRune(r) // 直接添加换行符
				state.emptyQuote = false   // 换行符后重置为空引号状态
				continue
			}
			// 处理其他空白字符分隔符
			state.handleSeparator()

		default:
			// 处理普通字符
			state.builder.WriteRune(r)
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
//   - r: 当前字符
func (s *splitState) handleQuoteChar(r rune) {
	switch {
	case !s.inQuotes:
		// 开始引号
		s.inQuotes = true   // 进入引号状态
		s.quote = r         // 记录当前引号类型
		s.emptyQuote = true // 假设为空引号，直到遇到非引号字符

	case r == s.quote:
		// 结束引号
		s.inQuotes = false
		// 处理空引号逻辑：非空引号或独立空引号才设置 hasQuoteInWord
		if !s.emptyQuote || s.builder.Len() == 0 {
			s.hasQuoteInWord = true
		}

	default:
		// 引号内的其他引号字符
		s.builder.WriteRune(r) // 直接添加引号字符
		s.emptyQuote = false   // 有内容，不是空引号
	}
}

// handleSeparator 处理分隔符（空格或制表符）的逻辑
//
// 判断逻辑:
//   - 1. builder.Len() > 0：片段中有内容（非空）
//   - 2. hasQuoteInWord：片段中包含过引号（空引号或非空引号）
//   - 这样可以确保独立空引号被保留，连续引号中的空引号被忽略
func (s *splitState) handleSeparator() {
	// 判断是否需要添加当前命令片段
	if s.builder.Len() > 0 || s.hasQuoteInWord {
		s.result = append(s.result, s.builder.String()) // 添加当前命令片段
		s.builder.Reset()                               // 重置构建器，准备下一个片段
		s.hasQuoteInWord = false                        // 重置标记，为下一个片段做准备
	}
}
