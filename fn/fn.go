package fn

import (
	"log"
	"reflect"

	"github.com/frankill/gotools/array"
)

// FuncType 是一个空接口，接受任意函数签名。
type FuncType any

// 检查一个接口值是否是函数类型
func isFunction(fn interface{}) bool {
	t := reflect.TypeOf(fn)
	return t.Kind() == reflect.Func
}
func getFunctionParamCount(fn interface{}) (int, bool) {
	// 获取函数的 reflect.Type
	funcType := reflect.TypeOf(fn)

	// 检查类型是否为函数
	if funcType.Kind() != reflect.Func {
		return 0, false
	}

	// 返回函数参数的数量
	return funcType.NumIn(), true
}

// FuncWrapper 包装一个函数，并允许进行部分应用和调用。
type FuncWrapper struct {
	fun    FuncType
	params []reflect.Value
}

// NewFuncWrapper 创建一个新的 FuncWrapper 实例，初始化时只设置函数，参数为空。
func NewFunc(f FuncType) *FuncWrapper {
	if !isFunction(f) {
		log.Fatalln("Invalid function type.")
	}
	return &FuncWrapper{fun: f, params: make([]reflect.Value, 0)}
}

// Partial 添加新的参数到现有参数列表，返回一个指向当前 FuncWrapper 实例的指针。
func (fw *FuncWrapper) Partial(args ...any) *FuncWrapper {
	for _, arg := range args {
		fw.params = append(fw.params, reflect.ValueOf(arg))
	}
	return fw
}

// Call 调用函数，将所有已部分应用的参数与新传入的参数一起传递给函数，返回函数的结果。
func (fw *FuncWrapper) Call(args ...any) any {
	// 创建参数列表
	callArgs := append(fw.params, array.Map(func(x any) reflect.Value {
		return reflect.ValueOf(x)
	}, args)...)

	num, ok := getFunctionParamCount(fw.fun)
	if !ok {
		return nil
	}
	if num != len(callArgs) {
		return nil
	}
	// 获取函数的 reflect.Value
	funcValue := reflect.ValueOf(fw.fun)

	// 调用函数
	result := funcValue.Call(callArgs)
	// 返回第一个结果（如果函数有返回值）
	if len(result) > 0 {
		return result[0].Interface()
	}
	return nil
}

func (fw *FuncWrapper) Clone() *FuncWrapper {
	return &FuncWrapper{fun: fw.fun, params: fw.params}
}
