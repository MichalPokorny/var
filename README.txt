Bit vector equation & SAT solver in Go.

License: GPL2

How to install and use:
1) Set up Go as described in https://golang.org/doc/code.html
2) Run `go get github.com/MichalPokorny/var`
3) To run tests: `go test -v github.com/MichalPokorny/var`

The following works:
	* Addition
	* Bitwise operations (&, |, ^)
	* Ordering (> >= < <=)
	* Multiplication
	* Division
	* Literal assignment
	* Left and right shift (<< >>); large bits in amount are ignored (as in x86)
	* Nonzero
	* Hamming weight (i.e., popcount)

Equality and nonequality can be composed from bitwise XOR and nonzero.
