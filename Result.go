// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了Result结构体，表示命令执行的结果，包含输出、错误、时间等完整信息。
package shellx

import (
	"os"
	"time"
)

// Result 表示命令执行的结果
type Result struct {
	// 基本执行信息
	cmd      *Command // 执行的命令引用
	exitCode int      // 退出码
	success  bool     // 是否执行成功

	// 输出信息
	stdout []byte // 标准输出内容
	stderr []byte // 标准错误内容
	output []byte // 合并输出

	// 时间信息
	startTime time.Time     // 开始执行时间
	endTime   time.Time     // 结束执行时间
	duration  time.Duration // 执行耗时

	// 进程信息
	pid          int              // 进程ID
	processState *os.ProcessState // 进程状态

	// 错误信息
	err error // 执行过程中的错误

	// 元数据
	metadata map[string]interface{} // 额外的元数据信息
}

// 提供公共访问方法
func (r *Result) Cmd() *Command                { return r.cmd }
func (r *Result) Code() int                    { return r.exitCode }
func (r *Result) Success() bool                { return r.success }
func (r *Result) StdOut() []byte               { return r.stdout }
func (r *Result) StdErr() []byte               { return r.stderr }
func (r *Result) Output() []byte               { return r.output }
func (r *Result) Start() time.Time             { return r.startTime }
func (r *Result) End() time.Time               { return r.endTime }
func (r *Result) Duration() time.Duration      { return r.duration }
func (r *Result) PID() int                     { return r.pid }
func (r *Result) State() *os.ProcessState      { return r.processState }
func (r *Result) Error() error                 { return r.err }
func (r *Result) Meta() map[string]interface{} { return r.metadata }
