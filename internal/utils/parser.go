// Package internal 提供shell命令执行库的内部工具函数。
// 本文件实现了命令字符串解析功能，支持引号处理和参数分割。
package internal

import (
	"strings"
)

// ParseCmd 将命令字符串解析为命令切片，支持引号处理(单引号、双引号、反引号)，出错时返回空切片
//
// 实现原理：
//  1. 去除首尾空白
//  2. 遍历每个字符
//  3. 处理引号状态切换
//  4. 在非引号状态下遇到空格时分割
//  5. 检查引号是否闭合
//
// 参数:
//   - cmdStr: 要解析的命令字符串
//
// 返回值:
//   - []string: 解析后的命令切片
func ParseCmd(cmdStr string) []string {
	// 去除首尾空白
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return []string{}
	}

	var (
		result    []string // 解析结果
		current   []rune   // 当前命令片段
		inQuotes  bool     // 是否在引号中
		quote     rune     // 当前引号类型
		hadQuotes bool     // 当前片段是否包含过引号
	)

	// 遍历每个字符
	for _, r := range cmdStr {
		if r == '"' || r == '\'' || r == '`' {
			if !inQuotes {
				inQuotes = true // 开始引号
				quote = r
				hadQuotes = true // 标记当前片段包含引号

			} else if r == quote {
				inQuotes = false // 引号闭合

			} else {
				current = append(current, r) // 引号内的字符直接添加
			}

		} else if (r == ' ' || r == '\t') && !inQuotes {
			if len(current) > 0 || hadQuotes {
				result = append(result, string(current)) // 非引号状态下遇到空格或制表符，添加当前命令片段
				current = current[:0]
				hadQuotes = false
			}

		} else {
			current = append(current, r)
		}
	}

	// 添加最后一个命令片段
	if len(current) > 0 || hadQuotes {
		result = append(result, string(current))
	}

	// 检查引号是否闭合
	if inQuotes {
		return []string{}
	}

	return result
}
