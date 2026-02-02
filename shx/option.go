package shx

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"mvdan.cc/sh/v3/expand"
)

// WithDir 设置工作目录
//
// 参数：
//   - dir: 工作目录路径
//
// 返回：
//   - *Shx: 命令对象（支持链式调用）
//
// 注意：
//   - 如果命令已经执行过，会 panic
//   - 如果目录不存在或不是目录，会 panic
func (s *Shx) WithDir(dir string) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if dir == "" {
		return s
	}

	// 验证目录是否存在
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("directory %s does not exist", dir))
		}
		panic(fmt.Sprintf("stat %s failed: %v", dir, err))
	}
	if !info.IsDir() {
		panic(fmt.Sprintf("%s is not a directory", dir))
	}

	// 转换为绝对路径
	absDir, err := filepath.Abs(dir)
	if err != nil {
		s.dir = dir
	} else {
		s.dir = absDir
	}

	return s
}

// WithEnv 设置环境变量
//
// 参数：
//   - key: 环境变量名
//   - value: 环境变量值
//
// 返回：
//   - *Shx: 命令对象（支持链式调用）
//
// 注意：
//   - 如果命令已经执行过，会 panic
//   - 如果 key 为空，则忽略
func (s *Shx) WithEnv(key, value string) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if key == "" {
		return s
	}

	// 获取当前环境变量列表
	var envList []string
	s.env.Each(func(name string, vr expand.Variable) bool {
		if name != key {
			envList = append(envList, fmt.Sprintf("%s=%s", name, vr.String()))
		}
		return true
	})

	// 添加新的环境变量
	envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	s.env = expand.ListEnviron(envList...)

	return s
}

// WithEnvs 批量设置环境变量
//
// 参数：
//   - envs: 环境变量映射（key-value）
//
// 返回：
//   - *Shx: 命令对象（支持链式调用）
//
// 注意：
//   - 如果命令已经执行过，会 panic
func (s *Shx) WithEnvs(envs map[string]string) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if len(envs) == 0 {
		return s
	}

	// 获取当前环境变量列表
	envList := make([]string, 0)
	s.env.Each(func(name string, vr expand.Variable) bool {
		envList = append(envList, fmt.Sprintf("%s=%s", name, vr.String()))
		return true
	})

	// 添加新的环境变量（会覆盖已有的）
	for key, value := range envs {
		if key != "" {
			envList = append(envList, fmt.Sprintf("%s=%s", key, value))
		}
	}

	s.env = expand.ListEnviron(envList...)
	return s
}

// WithStdin 设置标准输入
//
// 参数：
//   - r: 输入读取器
//
// 返回：
//   - *Shx: 命令对象（支持链式调用）
//
// 注意：
//   - 如果命令已经执行过，会 panic
//   - 如果 r 为 nil，会 panic
func (s *Shx) WithStdin(r io.Reader) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if r == nil {
		panic("stdin cannot be nil")
	}

	s.stdin = &expandEnvReader{reader: r}
	return s
}

// WithStdout 设置标准输出
//
// 参数：
//   - w: 输出写入器
//
// 返回：
//   - *Shx: 命令对象（支持链式调用）
//
// 注意：
//   - 如果命令已经执行过，会 panic
//   - 如果 w 为 nil，会 panic
func (s *Shx) WithStdout(w io.Writer) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if w == nil {
		panic("stdout cannot be nil")
	}

	s.stdout = &expandEnvWriter{writer: w}
	return s
}

// WithStderr 设置标准错误
//
// 参数：
//   - w: 错误输出写入器
//
// 返回：
//   - *Shx: 命令对象（支持链式调用）
//
// 注意：
//   - 如果命令已经执行过，会 panic
//   - 如果 w 为 nil，会 panic
func (s *Shx) WithStderr(w io.Writer) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if w == nil {
		panic("stderr cannot be nil")
	}

	s.stderr = &expandEnvWriter{writer: w}
	return s
}

// WithTimeout 设置超时时间
//
// 参数：
//   - d: 超时时间
//
// 返回：
//   - *Shx: 命令对象（支持链式调用）
//
// 注意：
//   - 如果命令已经执行过，会 panic
//   - 如果 d <= 0，则忽略（不设置超时）
func (s *Shx) WithTimeout(d time.Duration) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if d > 0 {
		s.timeout = d
	}

	return s
}

// WithContext 设置上下文
//
// 参数：
//   - ctx: 上下文
//
// 返回：
//   - *Shx: 命令对象（支持链式调用）
//
// 注意：
//   - 如果命令已经执行过，会 panic
//   - 如果 ctx 为 nil，会 panic
//   - 设置的上下文优先级高于 WithTimeout
func (s *Shx) WithContext(ctx context.Context) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if ctx == nil {
		panic("context cannot be nil")
	}

	s.ctx = ctx
	return s
}
