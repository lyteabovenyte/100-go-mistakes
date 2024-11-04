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
-  