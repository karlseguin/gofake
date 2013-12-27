package gofake

import (
	"errors"
	"testing"
	"github.com/karlseguin/gspec"
)

func TestReturnsTheStubbedValueMultipleTimes(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Stub(fake.GetEmail).Returning("invalid")
	spec.Expect(fake.GetEmail("leto")).ToEqual("invalid")
	spec.Expect(fake.GetEmail("paul")).ToEqual("invalid")
}

func TestReturnsTheDefaultValueMultipleTimes(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	spec.Expect(fake.GetEmail("leto")).ToEqual("leto@caladan.gov")
	spec.Expect(fake.GetEmail("paul")).ToEqual("leto@caladan.gov")
}

func TestLimitsTheStubToASingleInvokation(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Stub(fake.GetEmail).Returning("first").Once()
	spec.Expect(fake.GetEmail("leto")).ToEqual("first")
	spec.Expect(fake.GetEmail("paul")).ToEqual("leto@caladan.gov")
}

func TestLimitsTheStubToATheSpecifiedNumberOfInvokation(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Stub(fake.GetEmail).Returning("first").Times(2)
	spec.Expect(fake.GetEmail("leto")).ToEqual("first")
	spec.Expect(fake.GetEmail("jessica")).ToEqual("first")
	spec.Expect(fake.GetEmail("paul")).ToEqual("leto@caladan.gov")
}

//silly, but let's make sure it doesn't panic or anything
func TestStubbingMethodWithNoReturnIsANoop(t *testing.T) {
	fake := newFake()
	fake.Exec()
}

func TestStubsAllReturnValues(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Stub(fake.Count).Returning(32, nil)
	c, err := fake.Count()
	spec.Expect(c).ToEqual(32)
	spec.Expect(err).ToBeNil()
}

func TestStubsASingleReturnValueAndDefaultsTheRest(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Stub(fake.Count).Returning(22)
	c, err := fake.Count()
	spec.Expect(c).ToEqual(22)
	spec.Expect(err.Error()).ToEqual("invalid")
}

func TestDefaultsNilValues(t *testing.T) {
	spec := gspec.New(t)
	fake := newFake()
	fake.Stub(fake.Count).Returning(nil, "some error")
	c, err := fake.Count()
	spec.Expect(c).ToEqual(10)
	spec.Expect(err.Error()).ToEqual("some error")
}

type Repository interface {
	GetEmail(id string) string
	Exec()
	Count() (int, error)
}

type FakeRepository struct {
	Fake
}

func newFake() *FakeRepository {
	return &FakeRepository{New()}
}

func (f *FakeRepository) GetEmail(id string) string {
	r := f.Called(id)
	return r.String(0, "leto@caladan.gov")
}

func (f *FakeRepository) Exec() {
	f.Called()
}

func (f *FakeRepository) Count() (int, error) {
	r := f.Called()
	return r.Int(0, 10), r.Error(1, errors.New("invalid"))
}
