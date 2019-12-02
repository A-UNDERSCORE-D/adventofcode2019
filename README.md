# Advent of Code
This repo is my answers to the [advent of code](https://adventofcode.com). Solutions will be published once the leader boards are full

## Intcode ASM
I have created an intcode assembler. It uses a custom assembly language to 
provide a simpler method to writing intcode. Currently it supports the following operations:

|Operation|Desc|
|:---:|:---:|
|`add a b target`|Adds a to b, and stores the result in target|
|`mul a b target`|Multiplies a by b, and stores the result in target|
|`halt`| stops the program|

Literal numbers are supported for arguments, they will be added at the end of the program and pointed to.


Additionally, there are two kinds of special arguments. Numbers prefixed with `%` point to registers, these are also added to the end,
and their pointers are auto resolved. As a special case, `%0` points to position zero (it is equivalent to `$0`)

Numbers prefixed with `$` are raw pointers, they will be copied verbatim into the emitted intcode.
`$` also supports negative indexes, which will be resolved from the end of the program (but before the literal and register storage)

Running the assembler as simple as `go run intcodeasm.go`

### Example 
 
The following code sets position 0 to the number 1337 by multiplying 13 by 100 and then adding 37
```
% go run intcodeasm.go << EOF
heredoc> mul 13 100 %0
heredoc> add %0 37 %0
heredoc> EOF
Intcode for ASM: 2,9,10,0,1,0,11,0,99,13,100,37
```
