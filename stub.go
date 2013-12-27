package gofake

type Stub struct {
	outputs []interface{}
	count   int
}

func (s *Stub) Returning(outputs ...interface{}) *Stub {
	s.outputs = outputs
	return s
}

func (s *Stub) Once() *Stub {
	return s.Times(1)
}

func (s *Stub) Times(count int) *Stub {
	s.count = count
	return s
}

func (s *Stub) Invoke(inputs []interface{}) []interface{} {
	s.count--
	return s.outputs
}

func (s *Stub) Exhausted() bool {
	return s.count == 0
}

func (s *Stub) Assert(name string, t AssertionLogger) {

}
