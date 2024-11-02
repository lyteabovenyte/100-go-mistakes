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
- <b>```the main caveat of interfaces, as a way to create abstraction; is that abstraction should be discovered, not created.```</b> so we shouldn't desing with interfaces and wait for a concrete need. Said differently, we should create an interface when we need it, not when we foresee that we could need it. so before introducint an interface type we should ask this question: `Why not call the implementation directly?`
- it's always a best practice to declare the interface on the client or consumer side and place the actual implementation on the producer side. so every client could declare it's own interface and import just the behavior or functionality that it needs, not the full funcitonality. it is also good in the sense that there could be any dependency between the package that implements the actual functionality and the package that is declaring interface, cause in Go interfaces are implemented implicitly
- in the concept of returning a struct (actual implementation) or interface there is a rule of thumb:
    - 1. returning a struct
    - 2. accepting interface if possible
- by using `any` we lose some of the core aspect and benefits of golang as a statically typed language.
- **generics**: In a nutshell, this allows writing code with types
that can be specified later and instantiated when needed.
- 