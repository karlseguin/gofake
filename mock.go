package gofake

import (
	"reflect"
)

type Mock struct {
	expectedInputs []interface{}
	actualInputs   [][]interface{}
	outputs        []interface{}
	originalCount  int
	count          int
	never          bool
}

func (m *Mock) With(inputs ...interface{}) *Mock {
	m.expectedInputs = inputs
	return m
}

func (m *Mock) Returning(outputs ...interface{}) *Mock {
	m.outputs = outputs
	return m
}

func (m *Mock) Once() *Mock {
	return m.Times(1)
}

func (m *Mock) Never() *Mock {
	m.never = true
	m.originalCount = 0
	m.count = 0
	return m
}

func (m *Mock) Times(count int) *Mock {
	m.originalCount = count
	m.count = count
	return m
}

func (m *Mock) Invoke(inputs []interface{}) []interface{} {
	if m.actualInputs == nil {
		m.actualInputs = make([][]interface{}, 0, 1)
	}
	m.actualInputs = append(m.actualInputs, inputs)
	m.count--
	return m.outputs
}

func (m *Mock) Exhausted() bool {
	return m.count == 0
}

func (m *Mock) Assert(name string, t AssertionLogger) {
	if m.never == true && m.count < 0 {
		t.Errorf("expected %s to not be called, it was", name)
		return
	}

	if m.count > 0 {
		t.Errorf("expected %s to be called %d %s, was called %d %s", name, m.originalCount, pluralize(m.originalCount, "times", "time"), m.originalCount-m.count, pluralize(m.originalCount-m.count, "times", "time"))
		return
	}

	for callIndex, actualInputs := range m.actualInputs {
		for index, actual := range actualInputs {
			expected := m.expectedInputs[index]
			if expected == Any {
				continue
			}
			if expected == nil && (actual == nil || reflect.ValueOf(actual).IsNil()) {
				continue
			}
			if expected == actual {
				continue
			}
			t.Errorf("expected %s to be called with %v, got %v (call #%d)", name, m.expectedInputs, actualInputs, callIndex)
			return
		}
	}
}

func pluralize(count int, plural string, single string) string {
	if count == 1 {
		return single
	}
	return plural
}
