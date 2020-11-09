# assembly
An interpreter for the Assembly programming language.

	go get github.com/dkmccandless/assembly

Usage:

	assembly [resolution filename]

NB: This interpreter is a work in progress, and the informal specification below will change.

### Resolution structure

Assembly programs are known as _resolutions_. A resolution consists of, in order,
* a title,
* one or more Whereas clauses, and
* one or more Resolved clauses.

While it is not currently enforced by the interpreter, resolutions are required to adhere to proper parliamentary resolution form. In particular, excluding its title, a resolution should consist of a single English sentence with proper grammar, and its composition should be undertaken in accordance with exacting standards of thoroughness, clarity, and propriety.

#### Whereas clauses

Each Whereas clause may contain one (1) variable declaration, and it is further encouraged that Whereas clauses be utilized to provide context and explanation regarding the purpose of the resolution. Whereas clauses should serve as documentation—not just of what the resolution does, but why, and indeed why the resolution is necessary and appropriate.

#### Resolved clauses

After the motivation for the resolution has been properly established, its function is manifested in the Resolved clauses. Each may contain one (1) statement, which must not be a variable declaration.

### Literals

#### Strings

A string literal is denoted by enclosing quotation marks in the conventional manner.

#### Integers

For the sake of clarity, integers are expressed in the form of a cardinal followed by a parenthesized numeral, which must be properly delimited. Assembly supports signed integers with an internal representation size of sixty-four (64) bits.

### Variables

Variable identifiers must consist of a single capitalized word. Variables are declared via the `hereinafter`. 

Each variable must be declared before it is used in a Resolved clause or another variable declaration, each variable must be declared exactly once, and each declared variable must be used.

### Operators

#### Numeric

The following operators are recognized, according to the indicated order of precedence:

* Postfix operators: `squared`, `cubed`
* Unary prefix operators: `twice`, `thrice`
* Binary prefix operators: `sum`, `product`, `quotient`, `remainder`
* Infix operators: `less`

#### Relational

Within `if` statements only, expressions can be compared via the following operators:
* `equals`
* `exceeds` (numeric expressions only)

### Keywords

The following statement keywords are recognized:

In Whereas clauses:
Keyword|Function|Syntax example
-|-|-
`hereinafter`|variable declaration|`WHEREAS the Customary Greeting (hereinafter the Greeting) is "Hello, World!",`

In Resolved clauses:
Keyword|Function|Syntax example
-|-|-
`assume`|variable assignment|`BE IT RESOLVED that this Assembly directs Total to assume the value Total less one (1)`
`if`|conditional execution|`BE IT RESOLVED that if Quorum exceeds Attendance, the Secretary shall publish "This Assembly lacks a quorum."`
`publish`|print|`BE IT RESOLVED that the Secretary is instructed to publish the aforesaid Greeting.`

### Comments

For the facilitation of a resolution's clear and thorough documentation, the interpreter allows any text not recognized as a literal, keyword, or variable identifier to stand as commentary.
