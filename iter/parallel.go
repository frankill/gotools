package iter

import (
	"log"
)

// Pipeline 表示一系列处理步骤，可以对数据进行处理
type Pipeline[T any] struct {
	steps []func(chan T) chan T // 表示处理步骤的函数切片
	start func() chan T         // 创建数据流起点的函数
	end   func(ch chan T)       // 处理数据流终点的函数
}

// NewPipeline 创建一个新的 Pipeline 实例，支持任意数据类型
func NewPipeline[T any]() *Pipeline[T] {
	return &Pipeline[T]{}
}

// AddStep 向管道中添加一个处理步骤。
// 参数:
//
//	step: 处理步骤函数，该函数接受一个输入通道，返回一个输出通道。
//
// 返回:
//
//	void
func (p *Pipeline[T]) AddStep(f func(chan T) chan T) {
	p.steps = append(p.steps, f)
}

// SetStart 设置管道的起点函数
func (p *Pipeline[T]) SetStart(f func() chan T) {
	p.start = f
}

// SetEnd 设置管道的终点函数
func (p *Pipeline[T]) SetEnd(f func(chan T)) {
	p.end = f
}

// Compute 仅仅使用管道中的处理步骤，不包含起点和终点，直接对输入通道应用
func (p *Pipeline[T]) Compute(input chan T) chan T {
	ch := input
	for _, step := range p.steps {
		ch = step(ch) // 对通道应用每一个步骤
	}
	return ch // 返回最终的输出通道
}

// Run 启动管道处理流程。
// 参数:
//
//	input: 输入数据通道。
//
// 返回:
//
//	void
func (p *Pipeline[T]) Run() {
	if p.start == nil || p.end == nil {
		log.Panicln("Pipeline is missing start or end function")
	}

	input := p.start() // 获取起始通道
	output := input    // 初始输出通道是输入通道

	for _, step := range p.steps {
		output = step(output) // 对通道应用每一个步骤
	}

	p.end(output) // 处理最终的输出通道
}

var (
	parallerNum = 1
)

func SetParallel(parallel int) {

	bufferMutex.Lock()
	defer bufferMutex.Unlock()
	parallerNum = parallel
}
