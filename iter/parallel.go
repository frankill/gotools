package iter

import (
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

// Parallel 允许并行运行多个函数，并动态调整最大并发数
type Parallel struct {
	arr []func()       // 存储要执行的函数的切片
	wg  sync.WaitGroup // 用于同步 goroutine 的 WaitGroup
	mu  sync.Mutex     // 保护并发度调整的互斥锁
	num int            // 当前并发度
	sem chan struct{}  // 信号量，用于限制同时运行的协程数量
}

// NewParallel 创建一个新的 Parallel 实例，指定初始最大并发数
func NewParallel(initialConcurrency int) *Parallel {
	return &Parallel{
		num: initialConcurrency,
		sem: make(chan struct{}, initialConcurrency), // 创建带有初始并发数的信号量通道
	}
}

// SetConcurrency 动态调整最大并发数
func (p *Parallel) SetConcurrency(newConcurrency int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if newConcurrency == p.num {
		return // 并发数没有变化
	}

	// 创建新的信号量通道
	newSem := make(chan struct{}, newConcurrency)

	// 复制旧信号量通道中的信号到新的信号量通道
	// 避免修改正在运行的协程
	for i := 0; i < p.num; i++ {
		select {
		case <-p.sem:
			// 将旧的信号量迁移到新的信号量通道
			newSem <- struct{}{}
		default:
		}
	}

	// 将信号量通道替换为新的信号量通道
	p.sem = newSem
	p.num = newConcurrency
}

// Add 添加一个新的函数到并行执行列表
func (p *Parallel) Add(f func()) {
	p.arr = append(p.arr, f)
}

// Compute 并行执行所有添加的函数，并等待它们完成
func (p *Parallel) Compute() {
	for _, f := range p.arr {
		p.sem <- struct{}{} // 获取信号量，控制最大并发数
		p.wg.Add(1)         // 增加 WaitGroup 计数器
		go func(f func()) {
			defer p.wg.Done()          // 函数完成时减少 WaitGroup 计数器
			defer func() { <-p.sem }() // 释放信号量
			f()
		}(f)
	}

	p.wg.Wait() // 等待所有函数完成
}

// 获取系统负载的示例函数
func getSystemLoad() float64 {
	// 模拟系统负载变化（0.0 到 1.0 之间的值）
	return rand.Float64()
}

// 自动调整并发度的函数
func AutoAdjustConcurrency(p *Parallel, maxConcurrency int, minConcurrency int) {
	for {
		load := getSystemLoad()
		var newConcurrency int

		if load > 0.8 {
			// 系统负载高，减少并发度
			newConcurrency = minConcurrency
		} else if load < 0.2 {
			// 系统负载低，增加并发度
			newConcurrency = maxConcurrency
		} else {
			// 维持当前并发度
			newConcurrency = p.num
		}

		p.SetConcurrency(newConcurrency)

		time.Sleep(5 * time.Second) // 每5秒调整一次
	}
}

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
