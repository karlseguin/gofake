# GoFake
A stub helper for writing effective golang tests.

GoFake isn't magical. It doesn't eliminate as much boilerplate code as other libraries. And, critically, it still requires that you program against interfaces (which is just the static-typed way of doing thing, deal with it)

## Usage
The first thing to do is to look at the interface we want to fake:

```go
    type Repository interface {
      GetEmail(id string) string
    }
```

Next, let's jump ahead to how we'll use our fake implementation:

```go
    repo := tests.NewRepository()
    repo.Stub(repo.GetEmail).Returning("")
    email := Show("paul", repo)
    ...
```
The above creates our test repository (which will fulfill the `Repository` interface) and stubs the call to `GetEmail`, telling it to always return an empty string. Now, let's implement our fake:

```go
    import (
      "github.com/karlseguin/gofake"
    )

    type FakeRepository struct {
      gofake.Fake
    }

    func NewRepository() *FakeRepository {
      return &FakeRepository{gofake.New()}
    }

    func (f *FakeRepository) GetEmail(id string) string {
      returns := f.Called(id)
      return returns.String(0, "default@mail.com")
    }
```

By embedding `gofake.Fake` our implementation inherits a number of behaviors. The first is the `Stub` method, which we saw above and which tests will make use of. The other is the `Called` method which is used internally by our fake. 

## Stubs
Only a few methods are available on a stub:

* `Returning(values ...interface{})` - the values to return
* `Once()` - only stub this for 1 call (defaults to unlimited)
* `Times(n int)` - only stub this for `n` calls

## Mocks
Mocks are created with the `Expect` method, rather than the `Stub` method. They expose the same methods as stubs, with the addition of `With` and `Assert`:

```go
    func TestPassesTheIdToTheRepoToGetTheEmail(t *testing.T) {
      repo.Expect(repo.GetEmail).With("leto@caladan.gov").Returning("")
      Show("paul", repo)
      repo.Assert(t)
    }
```

It's ok to mix stubs and mocks as `Assert` on a stub always passes. Also, you can use `gofake.Any` as a paremeter to `With` in order to match any value.

* `With(inptus ...interface{})` - the input values we expect to be called
* `Returning(values ...interface{})` - the values to return
* `Once()` - only stub this for 1 call (defaults to unlimited)
* `Never()` - this method should never be called
* `Times(n int)` - only stub this for `n` calls

## Returns
The `*Return` type which `Called` returns has a `Values` array that contains the returned value as `interface{}`. A number of helper functions exist to type these (`String`, `Int`, `Uint64`, `Bool`, ....). There's a helper for all built-in types. Each helper also takes the default value to return.

Unlike the other helpers, the `Error` helper will convert a string to an error, so that it's possible to do:

```go
    fake.Stub(fake.Load).Returning(nil, "some error")
```

Alternatively, you can always pass an actual `error`.

For your own custom type, you'll want to do something like:
```go
    func (f *FakeRepository) LoadUser(id) *User {
      r := f.Called(id)
      if len(r.Values) <= 0 {
        return SOME_DEFAULT_USER_MAYBE_NIL
      }
      return r.Values[0].(*User)
    }
```

Which is the same as:

```go
    func (f *FakeRepository) LoadUser(id) *User {
      r := f.Called(id)
      return r.At(0, nil).(*User)
    }
```

## Stubs vs Mocks
Since these terms tend to be loosely defined, here's my definition: stubs are relaxed, mocks are strict. A stub doesn't care how or even if it's called, or with what arguments. Stubs provide the minimum amount of coupling to a dependent component and are well suited for focused unit tests. A mock expects to be called with specific parameters and a specific number of times. Mocks are used specifically to test the interaction between dependencies. Mocks are useful for basic sanity checks, but they're no substitute for integration tests (conversely, integration tests can remove the need for mocks altogether, though there are practical reasons (speed) to favor having some mock-using unit tests).

