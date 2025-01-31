# Creating an Interpreter in GO

This project follow the book: "Writing an Interpreter in Go" by Thorsten Ball

We will be building an interpreter with:

    - Lexer
    - Parser
    - Abstract Syntax Tree (AST)
    - Evaluator

We will understand what "tokens" are, what an abstract syntaxt tree is, how to build it, how to evaluate it and how to extend it with new data structures and built-in functions

We will be interpreting the made-up programming language "Monkey" which has the following features:

- test

- C-like syntax
- variable bindings
- integers and booleans
- arithmetic expressions
- build-in functions
- first-class and higher-order functions
- closures
- a string data structure
- an array data structure
- a hash data structure

An example of binding values to variable names in the Monkey programming language:

```
let age = 1;
let name = "Monkey";
let result = 10*(20/2);

// binding an array of integers
let myArray = [1,2,3,4,5];

// creating a hash
let thorsten = {"name": "Thorsten", "age": 28};

// Accessing Elements in an array
myArray[0] // => 1
thorsten["name"] // => "Thorsten"

// Binding functions to names
let add = fn(a,b) {return a+b;};

// Func with implied "return"
let add = fn(a,b) {a + b;};

```

"Monkey" supports recursion and higher-order functions (first-class functions too)


## LEXING

### Lexical Analysis (Converting Source Code to Tokens)

Source Code -> Tokens -> AST

Example:

```
let x = 5 + 5;
```

The example above if input to a LEXER may output as:

```
[
    LET,
    IDENTIFIER("x"),
    EQUAL_SIGN,
    INTEGER(5),
    PLUS_SIGN,
    INTEGER(5),
    SEMICOLON
]
```


### Defining Tokens

We start by defining our Tokens (we will extend them later)

To start we will treat all numbers as type "number" and all variable names as "identifiers"

For special characters that aren't variable names we will treat them as "keywords" (let, fn, etc.)

### The Lexer

The goal of the Lexer is to take source code as input and output the Tokens

We initialize the lexer with our source code and then repeatedly call NextToken() on it to go through the source code, token by token, character by character.

### Extending Token Set and Lexer

Now we will add functionality for: ==, !, !=, -, /, *, <, >, true, false, if, else, return

Add these characters to the test file for testing functionality

```
!-/*5;
5 < 10 > 5;
```

NOTE: the lexer's job is **NOT** to tell us if the code makes sense, works or contains errors.  It is simply to recognize valid alphanumeric inputs and turn them into Tokens.

It is a good practice in testing to cover:
- all tokens 
- produce off-by-one errors
- edge cases at end-of-file
- newline handling
- multi-digit number parsing
- etc.

### Start of a REPL

