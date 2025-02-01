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

## Parsers

A software component that takes input data (text or tokens) and builds a data structure (AST).

Also called "syntactic analysis"

It gives a structural representation of the input, checking for correct syntax in the process.

Example using JS:

```javascript
> var input = '{"name":"Thorsten", "age":28}';
> var output = JSON.parse(input);
> output
{name: 'Thorsten', age: 28}
> output.name
'Thorsten'
> output.age
28
```

Abstract Syntax Tree: The 'abstract' is based on the fact that certain details visible in the source code are omitted in the AST

Semicolons, newlines, whitespace, comments, braces, bracket and parentheses (depending on the language)

**Example Workflow**

```go
if (3 * 5 > 10) {
return "hello";
} else {
return "goodbye";
}
```

Assume we have a MagicLexer, MagicParser and the AST is built with JS objects ->

```go
> var input = 'if (3*5 > 10) {return "hello"; } else {return "goodbye";}';
> var tokens = MagicLexer.parse(input);
> MagicParser.parse(tokens);

{
	type: "if-statement",
	condition: {
		type: "operator-expression",
		operator: ">",
		left: {
			type: "operator-expression",
			operator: "*",
			left: {type: "integer-literal", value: 3},
			right: {type: "integer-literal", value: 5},
		},
		right: {type: "integer-literal", value:10}
	},
	consequence: {
		type: "return-statement",
		returnValue: {type: "string-literal", value: "hello"}
		},
	alternative: {
		type: "return-statement",
		returnValue: {type: "string-literal", value: "goodbye"}
	}
}
```

GOAL is to build a PARSER that accepts the Tokens generated by our LEXER to build an AST and construct instances of the AST while RECURSIVELY parsing tokens.

### Why Not Use a Parser Generator:

See the text and author explanation regarding the benefits of Parser Generators and:

- Context-Free Grammar (CFG)
- Backus-Naur Form (BNF)
- Extended Backus-Naur Form (EBNF)

### This Parser:

This parser will be a "top down operator precedence" parser (Pratt parser)

NOTE: Top-down vs Bottom-up Parsing:

- Top-down: Start with a root node of the AST and descend
- Bottom-up: Start with a leaf node of the AST and ascend

#### Trade-offs:

- Won't be blazingly fast
- won't have formal proof of correctness
- error-recovery process and detection of erroneous syntax will be flawed

### Parsing "Let" statements

variable bindings are STATEMENTS:

```javascript
let x = 10;
let y = 5;

let add = fn(a,b){
    return a + b;
}; 
```

let statements will consist of two parts: **Identifier** and **Expression**

IDENTIFIERS (above): ***x***, ***y***, and ***add***
EXPRESSIONS (above): **10**, **15** and **fn(a,b){return a+b;}**


STATEMENTS: don't produce values
EXPRESSIONS: produce values

What does and doesn't produce values depends on design choices for the programming language

In some, function literals are expressions and can be used in any place where any other expression is allowed
In others, function literals can only be a part of a function declaration statement (top-level of the program)

Note: the Object of this parser is to repeatedly advance the tokens and check the current token to decide what to do next.

It will either call another parsing function or throw an error.

In summary: when parsing statements we process tokens from left to right, expect or rejct the next tokens and if everything fits we return an AST Node

### Parsing Expressions

Parsing expressions is more challenging.

#### Operator Precedence

`5 * 5 + 10` for this expression the AST should represent something like `((5 *5) + 10)`

So multiplication should take precedence over addition

Similarly the parser must understand that in `5*(5+10)` the addition should take precedence

So parenthesis have a higher precedence than multiplication which has a higher precedence than addition

#### Reoccurring Tokens

In expressions, tokens of the same type can appear in multiple positions

Also in the expression `5 * (add(2,3) + 10)` the outer parenthesis indicate a grouped expression but the inner parenthesis represent a "call expression"

So now, the token's position depends on the context, the tokens before and after and their precedence

#### Expressions in Monkey

Aside from 'let' and 'return' everything is an expression

It has expressions with prefix operators:

```go
-5
!true
!false
```

Infix operators (binary operators):

```go
5 + 5
5 - 5
5 / 5
5 * 5
5 ** 5
```

Comparison operators:

```go
foo == bar
foo != bar
foo < bar
foo > bar
```

Parenthesis groups:

```go
5 * (5 + 5)
((5 + 5) * 5) * 5
```

Call expressions:

```go
add(2,3)
add(add(2,3), add(5,10))
max(5, add(5, (5 * 5)))
```
Identifiers as expressions:

```go
foo * bar / foobar
add(foo, bar)
```
#### Terminology

**prefix operator** is an operator "in front of" its operand `--5`

Here the operator `--` (decrement), the operand is the integer literal 5 and the operator is in the prefix position

**postix operator** is an operator "after" its operand `foobar++`

> Note the Monkey language doesn't include postfix operators for simplicity.  Consider adding.

**infix operator** is an operator that sits between its operands `5 * 8`

