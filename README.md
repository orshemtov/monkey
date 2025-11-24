# README

The Monkey programming language.

Implemented as part of the book "Writing An Interpreter In Go" by Thorsten Ball.

And expanded with features from the follow-up book "Writing A Compiler In Go".

Two awesome books that I highly recommend!

![Monkey Demo](.github/demo.gif)

The compiler CAN and WILL compile this thing, and the VM will happily execute it:

```rust
let fibonacci = fn(x) {
  if (x == 0) {
    return 0;
  } else {
    if (x == 1) {
      return 1;
    } else {
      fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};
fibonacci(15);
```

## Benchmark

### Compiled vs. Interpreted - The numbers

```shell
go run benchmark/main.go --engine=vm
engine=vm, result=9227465, duration=3.785189584s

go run benchmark/main.go --engine=eval
engine=eval, result=9227465, duration=10.705003625s
```

## Components

- Tokens
- Lexer
- Parser
- `AST` (Abstract Syntax Tree)
- Object system
- Evaluator
- `REPL`

And extended with:

- Bytecode
- Compiler
- Virtual Machine

## Data types

- Integer
- Boolean
- Null
- String
- Array
- Hash

## Features

- Let statements
- If/else statements
- Return statements
- Function literals
- Function calls
- Higher-order functions
- Closures
- String concatenation
- Array literals
- Array indexing
- Hash literals
- Hash indexing
- Built-in functions
  - `len`
  - `first`
  - `last`
  - `rest`
  - `push`
  - `print`
