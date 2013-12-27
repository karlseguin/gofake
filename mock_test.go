package gofake

import (
	"fmt"
	"github.com/karlseguin/gspec"
	"testing"
)

func TestMockReturnsTheValueOnceByDefaultTimes(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Expect(fake.GetEmail).Returning("invalid")
	spec.Expect(fake.GetEmail("leto")).ToEqual("invalid")
	spec.Expect(fake.GetEmail("paul")).ToEqual("leto@caladan.gov")
}

func TestMockIsLimitedToASingleInvocation(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Expect(fake.GetEmail).Returning("first").Once()
	spec.Expect(fake.GetEmail("leto")).ToEqual("first")
	spec.Expect(fake.GetEmail("paul")).ToEqual("leto@caladan.gov")
}

func TestMockIsLimitedToTheSpecifiedNumberOfInvocations(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Expect(fake.GetEmail).Returning("first").Times(2)
	spec.Expect(fake.GetEmail("leto")).ToEqual("first")
	spec.Expect(fake.GetEmail("jessica")).ToEqual("first")
	spec.Expect(fake.GetEmail("paul")).ToEqual("leto@caladan.gov")
}

func TestMockReturnsAllSpecifiedValues(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Expect(fake.Count).Returning(32, nil)
	c, err := fake.Count()
	spec.Expect(c).ToEqual(32)
	spec.Expect(err).ToBeNil()
}

func TestMockReturnsASingleSpecifiedValue(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Expect(fake.Count).Returning(22)
	c, err := fake.Count()
	spec.Expect(c).ToEqual(22)
	spec.Expect(err.Error()).ToEqual("invalid")
}

func TestMockReturnsADefaultOnNil(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Expect(fake.Count).Returning(nil, "some error")
	c, err := fake.Count()
	spec.Expect(c).ToEqual(10)
	spec.Expect(err.Error()).ToEqual("some error")
}

func TestMockIsNotValidWhenCalledButShouldNotBe(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	logger := newLogger()
	fake.Expect(fake.GetEmail).Never()
	fake.GetEmail("x")
	fake.Assert(logger)
	spec.Expect(logger.Count()).ToEqual(1)
	spec.Expect(logger.At(0)).ToEqual(`expected GetEmail to not be called, it was`)
}

func TestMockIsNotValidWhenNotCalledOnce(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	logger := newLogger()
	fake.Expect(fake.GetEmail)
	fake.Assert(logger)
	spec.Expect(logger.Count()).ToEqual(1)
	spec.Expect(logger.At(0)).ToEqual(`expected GetEmail to be called 1 time, was called 0 times`)
}

func TestMockIsNotValidWhenNotCalledTheSpecifiedNumberOftimes(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	logger := newLogger()
	fake.Expect(fake.GetEmail).Times(3)
	fake.GetEmail("x")
	fake.Assert(logger)
	spec.Expect(logger.Count()).ToEqual(1)
	spec.Expect(logger.At(0)).ToEqual(`expected GetEmail to be called 3 times, was called 1 time`)
}

func TestMockIsNotValidWhenCalledWithTheWrongInputs(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	logger := newLogger()
	fake.Expect(fake.GetEmail).With("leto")
	fake.GetEmail("paul")
	fake.Assert(logger)
	spec.Expect(logger.Count()).ToEqual(1)
	spec.Expect(logger.At(0)).ToEqual(`expected GetEmail to be called with [leto], got [paul] (call #0)`)
}

func TestMockIsNotValidWhenCalledWithTheWrongInputsForMultipleInvocations(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	logger := newLogger()
	fake.Expect(fake.GetEmail).With("leto").Times(2)
	fake.GetEmail("leto")
	fake.GetEmail("jessica")
	fake.Assert(logger)
	spec.Expect(logger.Count()).ToEqual(1)
	spec.Expect(logger.At(0)).ToEqual(`expected GetEmail to be called with [leto], got [jessica] (call #1)`)
}

func TestMockAssertsFailureOnMultipleArguments(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	logger := newLogger()
	fake.Expect(fake.Remove).With(22, false)
	fake.Remove(11, true)
	fake.Assert(logger)
	spec.Expect(logger.Count()).ToEqual(1)
	spec.Expect(logger.At(0)).ToEqual(`expected Remove to be called with [22 false], got [11 true] (call #0)`)
}

func TestMockAssertSuccessWithMatchingParameters(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	logger := newLogger()
	fake.Expect(fake.Remove).With(22, false)
	fake.Remove(22, false)
	fake.Assert(logger)
	spec.Expect(logger.Count()).ToEqual(0)
}

func TestMockAssertSuccessWithMatchingParametersUsingAny(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	logger := newLogger()
	fake.Expect(fake.Remove).With(Any, false)
	fake.Remove(1111, false)
	fake.Assert(logger)
	spec.Expect(logger.Count()).ToEqual(0)
}

func TestMockAssertSuccessWithNils(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	logger := newLogger()
	fake.Expect(fake.LogError).With(nil)
	fake.LogError(nil)
	fake.Assert(logger)
	spec.Expect(logger.Count()).ToEqual(0)
}

type RecordingLogger struct {
	records []string
}

func newLogger() *RecordingLogger {
	return &RecordingLogger{
		records: make([]string, 0, 1),
	}
}

func (l *RecordingLogger) Errorf(format string, args ...interface{}) {
	l.records = append(l.records, fmt.Sprintf(format, args...))
}

func (l *RecordingLogger) Count() int {
	return len(l.records)
}

func (l *RecordingLogger) At(index int) string {
	return l.records[0]
}
