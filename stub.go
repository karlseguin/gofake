package gofake

type Stub struct {
	values []interface{}
	count  int
}

func (s *Stub) Returning(values ...interface{}) *Stub {
	s.values = values
	return s
}

func (s *Stub) Once() *Stub {
	s.Times(1)
	return s
}

func (s *Stub) Times(count int) *Stub {
	s.count = count
	return s
}

func (s *Stub) Invoke(inputs []interface{}) []interface{} {
	s.count--
	return s.values
}

func (s *Stub) Exhausted() bool {
	return s.count == 0
}
