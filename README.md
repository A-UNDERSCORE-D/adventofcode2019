# Advent of Code
This repo is my answers to the [advent of code](https://adventofcode.com). Solutions will be published once the leader boards are full

## Intcode ASM
I have created an intcode assembler. It uses a custom assembly language to 
provide a simpler method to writing intcode. Currently it supports the following operations:

|Operation|Desc|
|:---:|:---:|
|`add a b target`|Adds a to b, and stores the result in target|
|`mul a b target`|Multiplies a by b, and stores the result in target|
|`ipt target`|accepts input and stores it in target|
|`out a`|outputs a|
|`jnz a target`|Jumps to target if a is not zero (intcode jump if true)|
|`jez a target`|Jumps to target if a is zero (intcode jump if false)| 
|`lt a b target`|If a is less than b, set target to 1, 0 otherwise|
|`eq a b target`|If a equals b, set target to 1, 0 otherwise|
|`halt`| stops the program|

The syntax for arguments is inspired by intel ASM.

|Prefix|Description|
|:---:|:---:|
|None |The number as it is literally|
|`r`  |A register|
|`$`  |A pointer / positional|
|`.`  |The current IP (can be added or subtracted from using `.1` and `.-1` respectively|
|`!`  |The current Instruction number excluding arguments (replaced with the correct IP for jumps etc). Can be added to or subtracted from using `!10` and `!-10` respectively|

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
