The expr package provides a simple evaluator for arithmetic integer
expressions. The syntax and operations are the same as in Go. Operands are
the native "int" type, except that unlike in Go, boolean values, which are
created by comparisons, are integer 1 (true) and 0 (false). Create a parsed
expression using Parse, and then evaluate it with Eval.
