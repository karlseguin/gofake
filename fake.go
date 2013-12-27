// A stub helper for writing effective golang tests.
package gofake

import (
	"reflect"
	"runtime"
	"strings"
)

var Any = new(interface{})

type AssertionLogger interface {
	Errorf(format string, args ...interface{})
}

type Invoker interface {
	Invoke(inputs []interface{}) []interface{}
	Exhausted() bool
	Assert(name string, t AssertionLogger)
}

type Fake struct {
	registry map[string][]Invoker
	runtime  map[string][]Invoker
}

func New() Fake {
	return Fake{
		registry: make(map[string][]Invoker),
		runtime:  make(map[string][]Invoker),
	}
}

func (f Fake) Stub(function interface{}) *Stub {
	stub := &Stub{count: -1}
	f.register(function, stub)
	return stub
}

func (f Fake) Expect(function interface{}) *Mock {
	mock := &Mock{}
	f.register(function, mock)
	return mock.Once()
}

func (f Fake) Assert(t AssertionLogger) {
	for name, invokers := range f.registry {
		for _, invoker := range invokers {
			invoker.Assert(name, t)
		}
	}
}

func (f Fake) register(function interface{}, invoker Invoker) {
	name := getFunctionName(reflect.ValueOf(function).Pointer())
	array, exists := f.registry[name]
	if exists == false {
		array = make([]Invoker, 0, 1)
	}
	f.registry[name] = append(array, invoker)
	f.runtime[name] = append(array, invoker)
}

func (f Fake) Called(inputs ...interface{}) *Return {
	pc, _, _, _ := runtime.Caller(1)
	name := getFunctionName(pc)
	if invokers, exists := f.runtime[name]; exists && len(invokers) > 0 {
		invoker := invokers[0]
		values := invoker.Invoke(inputs)
		if invoker.Exhausted() {
			f.runtime[name] = invokers[1:]
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
