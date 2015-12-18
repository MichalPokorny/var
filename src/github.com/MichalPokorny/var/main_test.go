package main

import (
	"sort"
	"testing"
	"github.com/MichalPokorny/var/sat"
	"github.com/MichalPokorny/var/bitvecsat"
)

type resultSorter struct {
	Elements [][]int
}

func (sorter resultSorter) Len() int {
	return len(sorter.Elements)
}

func (sorter resultSorter) Less(i, j int) bool {
	for idx := 0; idx < len(sorter.Elements[i]); idx++ {
		if sorter.Elements[i][idx] < sorter.Elements[j][idx] {
			return true
		}
		if sorter.Elements[i][idx] > sorter.Elements[j][idx] {
			return false
		}
	}
	return false
}

func (sorter resultSorter) Swap(i, j int) {
	tmp := sorter.Elements[i]
	sorter.Elements[i] = sorter.Elements[j]
	sorter.Elements[j] = tmp
}

func intSlicesEqual(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func resultSetsEqual(a [][]int, b [][]int) bool {
	sort.Sort(resultSorter{Elements: a})
	sort.Sort(resultSorter{Elements: b})
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !intSlicesEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

type binaryOperator func(a, b int, width uint) int;

func testBinaryOperator(t *testing.T, width uint, a int, b int, c int, problem bitvecsat.Problem, operator binaryOperator) {
	problem.PrepareSat(int(width))
	formula := problem.MakeSatFormula()
	forbidders := make([]sat.Clause, 0)

	// Contains 3-tuples of (A, B, C)
	foundSolutions := make([][]int, 0)

	for {
		formula.Clauses = append(formula.Clauses, forbidders...)
		solution := sat.Solve(formula)
		if solution == nil {
			break
		}

		aValue := problem.GetValueInAssignment(solution, a)
		bValue := problem.GetValueInAssignment(solution, b)
		cValue := problem.GetValueInAssignment(solution, c)

		foundSolutions = append(foundSolutions, []int{aValue, bValue, cValue})
		forbidders = append(forbidders, solution.MakeForbiddingClause())
	}

	expectedSolutions := make([][]int, 0)
	for a := 0; a < (1 << width); a++ {
		for b := 0; b < (1 << width); b++ {
			c := operator(a, b, width)
			expectedSolutions = append(expectedSolutions, []int{a, b, c})
		}
	}

	if !resultSetsEqual(foundSolutions, expectedSolutions) {
		t.Errorf("unexpected solutions: found %v, expected %v", foundSolutions, expectedSolutions)
	}
}

func operatorPlus(a int, b int, width uint) int {
	return (a + b) % (1 << width);
}

func TestAddition(t *testing.T) {
	problem := bitvecsat.Problem{}
	a := problem.AddNewVector()
	b := problem.AddNewVector()
	c := problem.AddNewVector()

	plus_constrain := bitvecsat.PlusConstrain{AIndex: a, BIndex: b, SumIndex: c}
	plus_constrain.AddToProblem(&problem)

	testBinaryOperator(t, 3, a, b, c, problem, operatorPlus)
}

func operatorOr(a int, b int, width uint) int {
	return a | b;
}

func TestOr(t *testing.T) {
	problem := bitvecsat.Problem{}
	a := problem.AddNewVector()
	b := problem.AddNewVector()
	c := problem.AddNewVector()

	or_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.OrConstrain}
	or_constrain.AddToProblem(&problem)

	testBinaryOperator(t, 3, a, b, c, problem, operatorOr)
}

func operatorAnd(a int, b int, width uint) int {
	return a & b;
}

func TestAnd(t *testing.T) {
	problem := bitvecsat.Problem{}
	a := problem.AddNewVector()
	b := problem.AddNewVector()
	c := problem.AddNewVector()

	and_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.AndConstrain}
	and_constrain.AddToProblem(&problem)

	testBinaryOperator(t, 3, a, b, c, problem, operatorAnd)
}

func operatorXor(a int, b int, width uint) int {
	return a ^ b;
}

func TestXor(t *testing.T) {
	problem := bitvecsat.Problem{}
	a := problem.AddNewVector()
	b := problem.AddNewVector()
	c := problem.AddNewVector()

	xor_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.XorConstrain}
	xor_constrain.AddToProblem(&problem)

	testBinaryOperator(t, 3, a, b, c, problem, operatorXor)
}

func operatorEquiv(a int, b int, width uint) int {
	return (1 << width - 1) - (a ^ b);
}

func TestEquiv(t *testing.T) {
	problem := bitvecsat.Problem{}
	a := problem.AddNewVector()
	b := problem.AddNewVector()
	c := problem.AddNewVector()

	equiv_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.EquivConstrain}
	equiv_constrain.AddToProblem(&problem)

	testBinaryOperator(t, 3, a, b, c, problem, operatorEquiv)
}

type binaryRelation func(a, b int, width uint) bool;

func testBinaryRelation(t *testing.T, width uint, a int, b int, problem bitvecsat.Problem, relation binaryRelation) {
	problem.PrepareSat(int(width))
	formula := problem.MakeSatFormula()
	forbidders := make([]sat.Clause, 0)

	t.Log(problem)
	t.Log(formula)

	// Contains 3-tuples of (A, B, C)
	foundSolutions := make([][]int, 0)

	for {
		formula.Clauses = append(formula.Clauses, forbidders...)
		solution := sat.Solve(formula)
		if solution == nil {
			break
		}

		aValue := problem.GetValueInAssignment(solution, a)
		bValue := problem.GetValueInAssignment(solution, b)

		foundSolutions = append(foundSolutions, []int{aValue, bValue})
		forbidders = append(forbidders, solution.MakeForbiddingClause())
	}

	expectedSolutions := make([][]int, 0)
	for a := 0; a < (1 << width); a++ {
		for b := 0; b < (1 << width); b++ {
			if relation(a, b, width) {
				expectedSolutions = append(expectedSolutions, []int{a, b})
			}
		}
	}

	if !resultSetsEqual(foundSolutions, expectedSolutions) {
		t.Errorf("unexpected solutions: found %v, expected %v", foundSolutions, expectedSolutions)
	}
}

func relationLte(a int, b int, width uint) bool {
	return a <= b;
}

func TestLte(t *testing.T) {
	for width := 1; width < 4; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector()
		b := problem.AddNewVector()

		lte_constrain := bitvecsat.LTEConstrain{AIndex: a, BIndex: b}
		lte_constrain.AddToProblem(&problem)

		testBinaryRelation(t, uint(width), a, b, problem, relationLte)
	}
}
