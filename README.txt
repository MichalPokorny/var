Bit vector equation & SAT solver in Go.

License: GPL2

The following works:
	* Addition
	* Bitwise operations (&, |, ^)
	* Ordering (> >= < <=)
	* Multiplication
	* Division
	* Literal assignment
	* Left and right shift (<< >>); large bits in amount are ignored (as in x86)
	* Nonzero

Equality and nonequality can be composed from bitwise XOR and nonzero.
