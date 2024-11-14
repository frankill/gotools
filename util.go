package gotools

import (
	"cmp"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Logical interface {
	bool
}

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}

type String interface {
	string
}

type Comparable interface {
	comparable
}

type Ordered interface {
	cmp.Ordered
}

func Identity[T any](x T) T {
	return x
}

var (
	ASCInt  = ASCGeneric[int]
	DESCInt = DESCGeneric[int]
)

func ASCGeneric[T Ordered](x, y T) bool {
	return x < y
}

func DESCGeneric[T Ordered](x, y T) bool {
	return x > y
}

func SysStop() context.Context {
	// 创建一个可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 创建一个通道用于接收系统信号
	signalChan := make(chan os.Signal, 1)
	// 监听指定的信号
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 在独立的goroutine中等待信号
	go func() {
		// 阻塞直到接收到信号
		<-signalChan
		// 接收到信号时取消上下文
		cancel()
	}()

	return ctx
}

func Clear(i int) {
	time.Sleep(time.Second * time.Duration(i))
}
