package gofake

import (
	"reflect"
	"runtime"
	"strings"
)

type Invoker interface {
	Invoke(inputs []interface{}) []interface{}
	Exhausted() bool
}

type Fake struct {
	registry map[string][]Invoker
}

func New() Fake {
	return Fake{
		registry: make(map[string][]Invoker),
	}
}

func (f Fake) Stub(function interface{}) *Stub {
	name := getFunctionName(reflect.ValueOf(function).Pointer())
	array, exists := f.registry[name]
	if exists == false {
		array = make([]Invoker, 0, 1)
	}
	stub := &Stub{count: -1}
	f.registry[name] = append(array, stub)
	return stub
}

func (f Fake) Called(inputs ...interface{}) *Return {
	pc, _, _, _ := runtime.Caller(1)
	name := getFunctionName(pc)
	if invokers, exists := f.registry[name]; exists && len(invokers) > 0 {
		invoker := invokers[0]
		values := invoker.Invoke(inputs)
		if invoker.Exhausted() {
			f.registry[name] = invokers[1:]
		}
		return &Return{values}
	}
	return &Return{make([]interface{}, 0)}
}

func getFunctionName(pc uintptr) string {
	name := runtime.FuncForPC(pc).Name()
	dot := strings.LastIndex(name, ".")
	function := name[dot+1:]
	if index := strings.LastIndex(function, "·"); index != -1 {
		return function[:index]
	}
	return function
}
