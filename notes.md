##### 100 common go mistakes review

###### main question to keep in mind:
- data-race bugs? and how to avoid them with golang?
- reducing allocation while parallelizing execution? how and where?
- the impact of data alignment in performance?
- variable shadowing and nested code abuse?
- using fallback mechanism in case of an error or a not desired condiiton happened.
- 

###### valuable notes:
- Feature-wise, Go has no type inheritance, no exceptions, no macros, no partial
functions, no support for lazy variable evaluation or immutability, no operator over-
loading, no pattern matching, and on and on. Why are these features missing from
the language? The official [Go FAQ](https://go.dev/doc/faq) gives us some valuable insight.
- **Variable shadowing**: in go, a variable declared in a block can be redeclared in an inner block, so we could have two variable with the same name, one in the outer block and one in the inner block and they are not the same.
```Go
// variable shadowing:
var client *http.Client
if tracing {
    client, err := CreateHttpClientWithTracing()
    if err != nil {
        return err
    }
    log.Println(client)
} else {
    client, err := CreateDefaultHttpClient()
    if err != nil {
        return err
    }
    log.Println(client)
}
// some logic with client
// and the client in the outer block (here) is still nil, cause we used ":=" in inner block and both the clients are not the same variable. --> variable shadowing
```
- When an if block returns, we should omit the else block in all cases. so called we should keep the happy path to be able to examine the execution flow on the column wise.
```Go
if condition {
    // some logic 
    return ...
}
if bluh-bluh {
    //
    return ...
}
//
```
- An <u>**init function**</u> is a function used to initialize the state of an application, When a package is initialized,
all the constant and variable declarations in the package are evaluated. **Then**, the init
functions are executed.
- *Global variables* have some severe drawbacks:
    - 1. Any functions can alter global variables within the package.
    - 2. Unit tests can be more complicated because a function that depends on a
         global variable won’t be isolated anymore.
- We should be cautious with **init functions**. They can be helpful in some situations,
however, such as defining [<u>static</u> configuration](https://cs.opensource.google/go/x/website/+/e0d934b4:blog/blog.go;l=32). Otherwise,
and in most cases, we should handle initializations through ad hoc functions.
- In programming, data encapsulation refers to hiding the values or state of an object.
*Getters and setters* are means to enable encapsulation by providing exported methods
on top of unexported object fields.
- An <u>**interface**</u> provides a way to specify the behavior of an object. We use interfaces to
create common abstractions that multiple objects can implement.
- cases we should consider using interfaces:
    - 1. common behavior: use interfaces when multiple types implement a
        common behavior. In such a case, we can factor out the <u>behavior</u> inside an interface.
        ```Go
        // all we need for sorting is these three behavior. whether it's mergesort or quicksort
        type Sort interface {
            Len() int // the number of element to sort
            Less(i, j int) bool // checking the "less than" to sort
            Swap(i, j int) // swap two element          
        }
        ```
    - 2. decoupling: decoupling our code from an implementation. If
        we rely on an abstraction instead of a concrete implementation, the implementation
        itself can be replaced with another without even having to change our code.
    - 3. restricting behavior
- <b>```the main caveat of interfaces, as a way to create abstraction; is that abstraction should be discovered, not created.```</b> so we shouldn't desing with interfaces and wait for a concrete need. Said differently, we should create an interface when we need it, not when we foresee that we could need it. so before introducing an interface type we should ask this question: `Why not call the implementation directly?` or is there any other object or usecase for this abstraction?
- it's always a best practice to declare the interface on the client or consumer side and place the actual implementation on the producer side. so every client could declare it's own interface and import just the behavior or functionality that it needs, not the full funcitonality. it is also good in the sense that there could be no dependency between the package that implements the actual functionality and the package that is declaring interface, cause in Go interfaces are implemented implicitly
- in the concept of returning a struct (actual implementation) or interface there is a rule of thumb:
    - 1. returning a struct
    - 2. accepting interface if possible
- by using `any` we lose some of the core aspect and benefits of golang as a statically typed language.
- **generics**: In a nutshell, this allows writing code with types
that can be specified later and instantiated when needed.
- using generics and the function parameter using Type Parameter is some how defining a **constraint interface** for a type and **instantiate** it to the parameter.
- What’s the difference between a constraint using `~int` or one using `int`? Using `int`
restricts it to that type, whereas `~int` restricts all the types whose underlying type is
an `int`
- [`comparable` interface](https://pkg.go.dev/builtin#comparable) as a type constraint interface: comparable is an interface that is implemented by all comparable types (booleans, numbers, strings, pointers, channels, arrays of comparable types, structs whose fields are all comparable types). The comparable interface may only be used as a **type parameter constraint**, not as the type of a variable.
- so the last word about generics and how to use them and don't using them, is that generics provide a kind of abstraction for the data types and functional behvaior as it can restrict and constraint the data types and the functional behvaior of an object, so `just solve what needs to be solved, and whenever you feel the need to write boilerplate code for a general type and writing biolerplate and general function for an object type, consider if using generics make your code clearer or not.`
- when using **embedded fields** to promote behavior or data (fields) to another struct be conscious that this promotion stays public to the struct and if the struct is public as well, clients and outside world could grant access to that field or behvaior. although using type embedding could be useful to prevent writing redundant boilerplate codes that are just passing functionality to the inner fields.
- by using `Builder` pattern we can segregate the logic of some optional parameters in a different struct and then by using methods on that struct, handle optional configuration and provide patterns on that specific field. --> downside is error handling and pointer value in Builder fields as the pointer values to distinguish the difference between nil value and the actual 0 value of that specific type. so for the actual nil value we should provide `nil` in the configBuilder struct's field. in error handling case, it's much harder to chain methods if the middle ones could throw an exception and return error. so the logic of error handling could defer to client to manage error cases.
- using functional options pattern is the idiomatic way of dealing with optional parameter in the context of the providing optional configuration.
with this exact type, we define a struct called options then a type that is a funciton to update the optional fields of the options struct
each field in the options needs a closure to deal with configuration. a **closure** in go is an anonymous function that references a field outside of its body.
so the client just have to provide the desired closure in using the API.
- using linters and formatter, documentating exported functions and methods and also packages and constant can make our code easier for reader and maintainer.

##### chapter 3.
- keep in mind that floating point values are approximate, so by this reason when we are comparing two floating point, in may return a result that may not be true. so `testify` package contains a `InDelta` function to assert that two values are in ceratin amount of distance by each other.
- in Go, a slice is backed by an array. that means that the slice's data is continuously added to the underlying array. `so internally a slice, holds a pointer to the backing array and a capacity and length variable`. The length is the number of elements the slice contains, whereas the capacity is the number of elements in the backing array.
-  after the backing array is full, go cope with this by doubling the capacity of the backing array. generally, In Go, a slice grows by doubling its size until it contains 1,024 elements, after which it grows by 25%.
-  when slicing a slice, hence, they are both backed by a single array, so changes to one, can impact the other, but appending have a different effect, which by appending to one of the slices, the appended element is only visible to the slice that it was appended to.
-  in the case of nil and empty slices, it's always a good practice to check whether the function that is working with the slice (either its from standard library or an external library) distinguish the difference between nil and empty slices. such as `encoding/json` which `Marshal()` return `null` for nil slice marshaling and returns `[]` for empty slice marshaling.
-  when designing Interfaces, it's a good practice to not to distinguish the difference between `nil` and empty slices. it means that we always should return `nil` for the both empty and nil slices, for the case where the caller's condition is `x != nil` or `len(x) != 0` which in both cases the `nil` value returns false and prevent subtle programming errors. and as a caller, it's always good to keep in mind to check the length of the slice, not whether it's nil or not.
- To use `copy` effectively, it’s essential to understand that the number of elements
copied to the destination slice corresponds to the minimum between. so, If we want to perform a complete copy, the destination slice must have a length greater than or equal to the source slice’s length.
- To copy one slice to another using the copy built-in function, remember that the number of copied elements corresponds to the minimum between the two slice’s lengths.
-  when slicing a slice(!), and the perform `append` on the latter one, we should be careful that if the capacity of the latter slice is greater than the element we are appending, the backing array of both may change, so the append operation may change the former(first) slice. in this kind of scenarios, we can use **full slice expression** which sets a limit on the sliced slice(!). this way, we are limiting on appending and if the limit we have set, don't contian capacity, appending doesn't affect the original slice and backing array.
-  Using `copy` or the *full slice expression* is a way to prevent `append` from creating
conflicts if two different functions use slices backed by the same array.
-  when working with slice is important to notice this rule, **a slice is a <u>pointer</u> to the backing array** so for more boldness --> `a slice is a pointer to the underlying backing array so when working with slices: if the element is a pointer or a struct with pointer fields, the elements won’t be reclaimed by the GC.`
- while benchmarking, consider using `runtime.MemStats` and it's method `Alloc`,  and `runtime.ReadMemStats()`. also `runtime.GC()` for manually calling garbage collector and `runtime.KeepAlive` to keep a reference of a variable to prevent garbage collector to collect it.
- in memory-leak point of view on slice and slicing subject there was two problems:
    - 1. The first was about slicing an existing slice or array to preserve the capacity. If we handle large slices and reslice them to keep only a fraction, a lot of memory will remain allocated but unused.
    - 2. The second problem is that when we use the slicing operation with *pointers* or structs
    with pointer fields, we need to know that the GC won't reclaim these elements.
###### maps:
- A `map` provides an <u>unordered</u> collection of key-value pairs in which all the keys are distinct. In Go, a map is based on the hash table data structure. Internally, a hash table is
an array of buckets, and each bucket is a pointer to an array of key-value pairs.
- just like with slices, initializing a size for the map, up front, remove the potential of reallocation of memory and rebalancing all the current buckets.
- it worth remembering that maps can only decrease the bucket but never shrinks. so we should be careful to not to allocate much more memory for map and bucket if we are gonna delete it and never use the bucket --> the solution may be to copy the actual map to another map with lower bucket size or use pointer as the values (which allocates 8 bit on 64-bit machine).
###### comparison:
- note that `==` and `!=` operators dont work with slice and maps. they only works on operand that are comparable (comparable is an interface on `Boolean`, `Numberics`, `Strings`, `Channels`, `Interface`, `Pointers`, `Structs`, `Array`)
- We can also use the ?, >=, <, and > operators with numeric types to compare values and with strings to compare their lexical order.
#### chapter 4: Control structure
- In general, `range` produces two values for each data structure except a receiving channel
- `Important Notice`: `In Go, everything we assign is a copy`:
    - 1. If we assign the result of a function returning a struct, it performs a copy of that
    struct.
    - 2. If we assign the result of a function returning a pointer, it performs a copy of the
    memory address (an address is 64 bits long on a 64-bit architecture). 
- It’s crucial to keep this in mind to avoid common mistakes, including those related to
range loops. Indeed, when a range loop iterates over a data structure, it performs a
copy of each element to the value variable (the second item).
- **expression evaluated only once in `range`:
```go
s := []int{1, 2, 3}
for range s {
    s = append(s, 10)
}
// keep in mind that range will copy the expression in a temp variable and uses that
// to evaluate the expression, so the size of the slice remain 3 for the temp variable
// used by range. it's not the case with traditional for loop.
```
- When iterating over a data structure using a range loop, we must recall that all the
values are assigned to a unique variable with a single unique address.
```go
for i, v := range []Customers{
    // all the elements of Customers slice
    // are assigned to the same v variable 
    // created by the range loop.
    // so keep in mind that this variable is 
    // single and points to a single memory address.
}
```
- maps are unordered, if you want to keep the order, consider using **binary heap**
- If a map entry is created during iteration, it may be produced during the iteration or
skipped. The choice may vary for each entry created and from one iteration to the next.
- when using `break` in for loop in conjuction with `switch` or `select`, the break statement doesn't terminate the for loop, it terminates the switch statement. so again, One essential rule to keep in mind is that a `break statement terminates the execution of the innermost for, switch, or select statement.` so how to break the loop? **use labels** as so:
```go
loop:
    for i := 0; i < 5; i++ {
        fmt.Println(i)
        switch i {
            default: // do what you want
            case 2: break loop // breaking the loop attached to the label.
        }
    }
```
- The `defer` statement delays a call’s execution until the surrounding function returns.
- `defer` schedules a function call when the surrounding function returns.

#### chapter 5: Strings.
- encoding vs charset:
    - **charset**: A charset, as the name suggests, is a set of characters. For example, the Unicode
    charset contains 2^21 characters. 
    - **encoding**: An encoding is the translation of a character’s list in binary. For example, UTF-
    8 is an encoding standard capable of encoding all the Unicode characters in a
    variable number of bytes (from 1 to 4 bytes).
- We mentioned characters to simplify the charset definition. But in Unicode, we use
the concept of a `code point` to refer to an item represented by a single value. For example, the 汉 character is identified by the U+6C49 code point. Using UTF-8, 汉 is
encoded using three bytes: 0xE6, 0xB1, and 0x89. Why is this important? Because in
Go, a **rune** is a **Unicode code point**. Meanwhile, we mentioned that UTF-8 encodes characters into 1 to 4 bytes, hence, up to 32 bits. This is why in Go, a rune is an alias of int32.
- as the built-in `len()` function returns the number of bytes, we should be caution about charset characters
- A charset is a set of characters, whereas an encoding describes how to translate
a charset into binary.
- In Go, a string references an immutable slice of arbitrary bytes.
- A rune corresponds to the concept of a Unicode code point, meaning an item
represented by a single value.
- When iterating over a string, we don’t iterate over each rune; instead, we iterate over each starting index of a rune
- Using a range loop on a string returns two variables, the starting index of a rune and the rune itself
- In summary, if we want to iterate over a string’s runes, we can use the range loop on
the string directly. But we have to recall that the index corresponds not to the rune
index but rather to the starting index of the byte sequence of the rune
- `TrimRight` removes all the trailing runes contained in a given set.
- `Trim` applies both TrimLeft and TrimRight on a string.
- consider this code:
```go
func concat(values []string) string {
    s := ""
    for _, v := range values {
        s += v
    }
    return s
}
// because strings immutability
// this code lack perfomance due
// to reallocating s, in every iteration
```
- the solution to the code above: `strings.Builder{}` struct. def in golang package: A **Builder** is used to efficiently build a string using `Builder.Write` methods. It minimizes memory copying. The zero value is ready to use. Do not copy a non-zero Builder.
```go
func concat(values []string) string {
    sb := strings.Builder{}
    for _, v := range values {
        _, _ := sb.WriteString(v) // appends a string to the builder struct
    }
    sb.String()
}
```
we created a strings.Builder struct using its zero value. During each iteration,
we constructed the resulting string by calling the WriteString method that appends
the content of value to its **internal buffer**, hence minimizing memory copying.
- using `bytes` package instead of `strings`, prevent us of unnecessary conversions. so bother using bytes if possible.
- When doing a substring operation, the Go specification doesn’t specify whether
the resulting string and the one involved in the substring operation should share the
same data. However, the standard Go compiler does let them share the same backing
array, which is probably the best solution memory-wise and performance-wise as it prevents a new allocation and a copy.
- Because a string is mostly a <u>pointer</u>, calling a function to pass a string
doesn’t result in a deep copy of the bytes. The copied string will still reference
the same backing array.
- so, We need to keep two things in mind while using the substring operation in Go.
*First*, the interval provided is based on the number of bytes, not the number of runes.
*Second*, a substring operation may lead to a memory leak as the resulting substring
will share the same backing array as the initial string. The solutions to prevent this
case from happening are to perform a string copy manually or to use `strings.Clone`
from Go 1.18.
#### chapter 6. functions and methods
- In Go, input and output operations are achieved using primitives that model data as streams of bytes that can be `read` from or `written` to
- In most cases, using named result parameters in the context of an interface definition can increase readability without leading to any side effects.
- a nil pointer is a valid receiver
- In Go, a method is just syntactic sugar for a function whose first parameter is the receiver.
- in Go, having a nil receiver is allowed, and an interface converted from a nil pointer isn’t a nil  interface. For that reason, when we have to return an interface, we should return not a nil pointer but a nil value directly. Generally, having a nil pointer isn’t a desirable state and means a probable bug.
- We need to understand something crucial about argument evaluation in a `defer`
function: the arguments are evaluated **right away**, not once the surrounding function
returns.
- in the case of calling a closure as a defer statement. As a reminder, a
*closure* is an anonymous function value that references variables from outside its
body. The arguments passed to a defer function are evaluated *right away*. But we must
know that the variables referenced by a defer closure are evaluated during the closure
execution (hence, when the surrounding function returns).
```go
func main() {
    i, j := 0, 0
    defer func() {
        fmt.Println(i, j)
    }(i)
    i++
    j++
}
// as this closure accept i as a function argumen
// i is evaluated right away
// but it references j from outside of it's body
// so j gets evaluated during the execution of the closure.
// prints --> 0, 1
```
- In summary, when we call defer on a function or method, the call’s arguments are
evaluated immediately. If we want to mutate the arguments provided to defer afterward, we can use *pointers* or *closures*. For a method, the receiver is also evaluated
immediately; hence, the behavior depends on whether the receiver is a *value* or a
*pointer*
- The decision whether to use a value or a pointer receiver should be made based
on factors such as the type, whether it has to be mutated, whether it contains a
field that can’t be copied, and how large the object is. When in doubt, use a
pointer receiver.
- `When returning an interface, be cautious about returning not a nil pointer but
an explicit nil value`, as the nil pointer is not nil, it is a memory address that points to a nil value, so the caller never get's a nil value by calling it.
- returning nil pointer is not explicit `nil`.
- Designing functions to receive `io.Reader` types instead of filenames improves
the reusability of a function and makes testing easier.
- Passing a pointer to a defer function and wrapping a call inside a closure are
two possible solutions to overcome the immediate evaluation of arguments and
receivers.
#### chapter 7. Error management
- In Go, `panic` is a built-in function that stops the ordinary flow
-  errors are returned as normal return values.
-  Once a `panic` is triggered, it continues up the call stack until either the current goroutine has returned or panic is caught with `recover`.
-  Note that calling `recover()` to capture a goroutine panicking is only useful inside a
`defer` function; otherwise, the function would return nil and have no other effect. This
is because defer functions are also executed when the surrounding function panics.
- so in case of panicing, when there is a need to signal a programmer error or a mandatory dependency injection, on in case of a function call in `init` function of a package, there is a good reason to panic and signal that an error occured and exit the application. but in other cases, it is better to manage the error in a function returning a proper error type.
- Since Go 1.13, the `%w` directive allows us to wrap errors conveniently. Error wrapping is about wrapping or packing an error inside a wrapper container that also makes the source error available. In general, the two main use cases for error wrapping are the following:
      - 1. Adding additional context to an error
      - 2. Marking an error as a specific error
- To make sure our clients don’t rely on something that we consider implementation
details, the error returned should be transformed, not wrapped. In such a case, using
`%v` instead of `%w` can be the way to go.
- notes about %v which transforms and %w which wraps the error:
    - 1. `%v` directive --> `fmt.Errorf("failed. error: %v", err)` returns `*errors.errorString`
    - 2. `%w` directive --> `fmt.Errorf("failed. error %w, err)` returns `*fmt.wrapError`
- If we need to mark an error, we should create a custom error type. However, if we just want
to add extra context, we should use fmt.Errorf with the %w directive as it doesn’t
require creating a new error type. Yet, error wrapping creates potential coupling as it
makes the source error available for the caller. If we want to prevent it, we shouldn’t use
error wrapping but error transformation, for example, using fmt.Errorf with the %v
directive.
- Go 1.13 came with a directive to wrap an error and a way to check whether the wrapped error is of a certain type with `errors.As`. This function requires the second argument (the target error) to be a pointer. Otherwise, the function will compile but panic at runtime.
```go
// transietError is custom error type for temp error
type transietError struct {
    err error
}

// Error() implement the error interface for transietError type
func (t *transietError) Error() string {
    // 
}

func f() error {
    // some logic
    if err != nil {
        if errors.As(err, &transietError) {
            // some logic in case of the error being of type transietError
        } else {
            // some logic for another type of error
        }
    }
}
```
- A **sentinel error** is an error defined as a global variable. A sentinel error conveys an *expected* error. Conversely, situations like network issues and connection polling
errors are unexpected errors. It doesn’t mean we don’t want to handle unexpected
errors; it means that semantically, those errors convey a different meaning.
- as general guidelines:
    - 1. Expected errors should be designed as error values (sentinel errors): var
    ErrFoo = errors.New("foo").
    - 2. Unexpected errors should be designed as error types: `type BarError struct
    { … }`, with BarError implementing the error interface.
- We have seen how `errors.As` is used to check an error against a type. With error <u>values</u>, we can use its counterpart: `errors.Is`.
```go
// creating a sentinel error
var NoRowFound = errors.New("didn't fount the specified row")

err := Query()
if err != nil {
    if errors.Is(err, NoRowFound) {
        // some logic
    } else {
        // some logic
    }
}
```
- Using `errors.Is` instead of the `==` operator allows the comparison to work even if the
error is wrapped using `%w`.
- In summary, if we use error wrapping in our application with the %w directive and
fmt.Errorf, checking an error against a specific value should be done using
errors.Is instead of ==. Thus, even if the sentinel error is wrapped, errors.Is can
recursively unwrap it and compare each error in the chain against the provided value.
#### chapter 8, concurrency foundamental
- Unlike *parallelism*, which is about doing the same thing multiple times at once, *concurrency* is
about structure.
- **`concurrency enables parallelism.`**, concurrency provides a structure to solve a problem with parts that may be parallelized.
- so concurrency is about the desing of the *threads* which are going to be coordinated and be aware of each other's state, and these threads then can be paralleled and provide high through-put.
- In summary, concurrency and parallelism are different. Concurrency is about structure, and we can change a sequential implementation into a concurrent one by introducing different steps that separate concurrent threads can tackle. Meanwhile, parallelism is about *execution*, and we can use it at the step level by adding more parallel threads. Understanding these two concepts is fundamental to being a proficient Go developer.
- A **thread** is the smallest unit of processing that an OS can perform. If a process wants
to execute multiple actions simultaneously, it spins up multiple threads. these threads can be:
  - 1. **Concurrent**: Two or more threads can start, run, and complete in overlapping
    time periods, like the waiter thread and the coffee machine thread in the previous section.
  - 2. **Parallel**: The same task can be executed multiple times at once, like multiple
    waiter threads.
- The OS is responsible for scheduling the thread’s processes optimally. A CPU core executes different threads. When it switches from one thread to another, it executes an operation called *context switching*. The active thread consuming CPU cycles was in an *executing* state and moves to a *runnable* state, meaning it’s ready to be executed pending an available core.
- As Go developers, we can’t create threads directly, but we can create **goroutines**,
which can be thought of as application-level threads.
- as mentioned in the source code:
    - Goroutine scheduler: The scheduler's job is to distribute *ready-to-run* goroutines over worker threads.
    - the Go scheduler uses the following terminology:
      - G - Goroutine
      - M - Machine (OS thread)
      - P - CPU Core (Processor)
- Each OS thread (M) is assigned to a CPU core (P) by the OS scheduler. Then, each
goroutine (G) runs on an M. The `GOMAXPROCS` variable defines the limit of Ms in
charge of executing user-level code simultaneously. But if a thread is blocked in a system call (for example, I/O), the scheduler can spin up more Ms. As of Go 1.5, `GOMAXPROCS` is by default equal to the number of available CPU cores.
- **channels** are a communication mechanism. Internally, a channel is a pipe we can use to send and receive values and that allows us to *connect* concurrent goroutines. A channel can be either of the following:
    - 1. *UnBuffered*: The sender goroutine blocks until the receiver goroutine is ready.
    - 2. *Buffered*: The sender goroutine blocks only when the buffer is full.
- In general, parallel goroutines have to *synchronize*: for example, when they need to
access or mutate a shared resource such as a slice. Synchronization is enforced with
**mutexes** but not with any channel types (not with buffered channels). Hence, in general, synchronization between parallel goroutines should be achieved via mutexes.
- Conversely, in general, concurrent goroutines have to *coordinate and orchestrate*. For
example, if G3 needs to aggregate results from both G1 and G2(which are parallel with each other and concurrent with G3), G1 and G2 need to signal to G3 that a new intermediate result is available. This coordination falls under the scope of communication—therefore, **channels**.
- **Mutexes** and **channels** have different semantics. Whenever we want to share a state
or access a shared resource, mutexes ensure exclusive access to this resource. Conversely, channels are a mechanic for signaling with or without data (chan struct{} or
not). Coordination or ownership transfer should be achieved via channels. It’s important to know whether goroutines are parallel or concurrent because, in general, we need mutexes for parallel goroutines and channels for concurrent ones.
- A **data race** occurs when two or more goroutines simultaneously access the same memory location and at least one is writing.
- A *race condition* occurs when the behavior depends on the sequence or the timing of events that can’t be controlled. Ensuring a specific execution sequence among goroutines is a question of *coordination* and *orchestration*
- If we want to ensure that we first go from state 0 to state 1,
and then from state 1 to state 2, we should find a way to guarantee that the goroutines
are executed in order. Channels can be a way to solve this problem. Coordinating and
orchestrating can also ensure that a particular section is accessed by only one goroutine, which can also mean removing the mutex
- The [Go memory model](https://golang.org/ref/mem) is a specification that
defines the conditions under which a read from a variable in one goroutine can be
guaranteed to happen after a write to the same variable in a different goroutine. In
other words, it provides guarantees that developers should keep in mind to avoid data
races and force deterministic output.
- A send on a channel happens before the corresponding receive from that chan-
nel completes. In the next example, a parent goroutine increments a variable
before a send, while another goroutine reads it after a channel read:
```go
i := 0
ch := make(chan struct{})

go func() {
    <-ch
    fmt.Println(i)
}()
i++
ch <- struct{}{}

// execution order is as follows:
// variable increment, channel send, channel receive, variable read from println
```
so by transitivity, we can ensure that accesses to i are synchronized and hence free from data races.
- Closing a channel happens before a receive of this closure.
- We can use the `runtime.GOMAXPROCS(int)` function to update the value of GOMAXPROCS. Calling it with 0 as an argument doesn’t change the value; it just returns the current value:
```go
n := runtime.GOMAXPROCS(0) // return the current value of the number of logical CPU
```
- When implementing the worker-pooling pattern, we have seen that the optimal number of goroutines in the pool depends on the workload type. If the workload executed by the workers is I/O-bound, the value mainly depends on the external system. Conversely, if the workload is CPU-bound, the optimal number of goroutines is close to the number of available threads(`runtime.NumCPU()` or `runtime.GOMAXPROC()`). Knowing the workload type (I/O or CPU) is crucial when designing concurrent applications.
- A **Context** carries a deadline, a cancellation signal, and other values across API
boundaries.
- The context.Context type exports a `Done` method that returns a receive-only notification channel: `<-chan struct{}`. This channel is closed when the work associated with the context should be canceled.
- the internal channel should be closed when a context is canceled or has met a deadline, instead of when it receives a specific value, because the closure of a channel is the only channel action that all the consumer goroutines will
receive. This way, all the consumers will be notified once a context is canceled or a deadline is reached.
- `context.Context` exports an `Err` method that returns nil if the
Done channel isn’t yet closed. Otherwise, it returns a non-nil error explaining why the
Done channel was closed
- When in doubt about which context to use, we should use `context.TODO()` instead
of passing an empty context with context.Background. `context.TODO()` returns an
empty context, but semantically, it conveys that the context to be used is either
unclear or not yet available (not yet propagated by a parent, for example)
- Concurrency is about structure, whereas parallelism is about execution.
#### chapter 9. concurrency practices.
- A `context.Context` is an interface containing four methods:
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{} // recieve-only channel to signal the cancelation
    Err() error 
    Value(key any) any
}
```
- The context’s deadline is managed by the `Deadline` method and the cancellation signal is managed via the `Done` and `Err` methods. When a deadline has passed or the context has been canceled, `Done` should return a closed channel, whereas `Err` should return an error. Finally, the values are carried via the Value method.
- Channels are a mechanism for communicating across goroutines via signaling.
- An empty struct is a de facto standard to convey an absence of meaning.
```go
var m struct{}
fmt.Println(unsafe.Sizeof(m)) // 0
```
so, For example, if we need a hash set structure (a collection of unique elements), we should use an
empty struct as a value: `map[K]struct{}`. ( a **set** is map that the key type is any and the value could be anything and we don't care )
- A channel can be with or without data. If we want to design an idiomatic API in
regard to Go standards, let’s remember that a channel without data should be
expressed with a `chan struct{}` type. This way, it clarifies for receivers that they
shouldn’t expect any meaning from a message’s content, only the fact that they have
received a message. In Go, such channels are called **notification channels**.
- The select statement lets a goroutine wait on multiple operations at the same time.
- Receiving from a closed channel is a non-blocking operation.
- receiving from a nil channel will block forever. in general, waiting or sending to a nil channel is a blocking action, and this behavior isn’t useless. we can use nil channels to implement an elegant state machine that will remove one case from a select statement.
- In concurrency, synchronization
means we can guarantee that multiple goroutines will be in a known state at some
point. For example, a mutex provides synchronization because it ensures that only
one goroutine can be in a critical section at the same time. Regarding channels:
  - An unbuffered channel enables synchronization. We have the guarantee that
    two goroutines will be in a known state: one receiving and another sending a
    message.
  -  A buffered channel doesn’t provide any strong synchronization. Indeed, a pro-
    ducer goroutine can send a message and then continue its execution if the
    channel isn’t full. The only guarantee is that a goroutine won’t receive a message before it is sent. But this is only a guarantee because of causality (you don’t
    drink your coffee before you prepare it).
- There are other cases where unbuffered channels are preferable: for example, in
the case of a notification channel where the notification is handled via a channel clo-
sure (`close(ch)`). Here, using a buffered channel wouldn’t bring any benefits.
- 