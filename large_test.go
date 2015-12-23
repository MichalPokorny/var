package main

import (
	"testing"
	"github.com/MichalPokorny/var/bitvecsat"
)

func relationIsSquare(a int, b int, width uint) bool {
	return ((a * a) % (1 << width)) == b;
}

func TestSquare(t *testing.T) {
	for width := uint(1); width <= 4; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)

		multiply_constrain := bitvecsat.MultiplyConstrain{AIndex: a, BIndex: a, ProductIndex: b}
		multiply_constrain.AddToProblem(&problem)

		testBinaryRelation(t, width, a, b, problem, relationIsSquare)
	}
}
