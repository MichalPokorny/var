package main

import (
	"strconv"
	"fmt"
	"sort"
	"testing"
	"github.com/MichalPokorny/var/sat"
	"github.com/MichalPokorny/var/sat/dfs"
	"github.com/MichalPokorny/var/sat/dpll"
	"github.com/MichalPokorny/var/sat/cdcl"
	"github.com/MichalPokorny/var/bitvecsat"
)

func solve(formula sat.Formula) sat.Assignment {
	if false {
		// DPLL
		return dpll.Solve(formula, sat.MakeEmptyAssignment(formula))
	}

	if false {
		// DFS
		return dfs.Solve(formula)
	}

	// CDCL
	return cdcl.Solve(formula)
}

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

func makeUnique(a [][]int) [][]int {
	r := make([][]int, 0)
	r = append(r, a[0])
	for i := 1; i < len(a); i++ {
		if !intSlicesEqual(a[i - 1], a[i]) {
			r = append(r, a[i])
		}
	}
	return r
}

func resultSetsEqual(a [][]int, b [][]int) bool {
	sort.Sort(resultSorter{Elements: a})
	sort.Sort(resultSorter{Elements: b})

	// For example, division may generate multiple solutions.
	a = makeUnique(a)
	b = makeUnique(b)

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

func resultSetContains(results [][]int, result []int) bool {
	for i := 0; i < len(results); i++ {
		if intSlicesEqual(results[i], result) {
			return true
		}
	}
	return false
}

func getResultsDifference(a [][]int, b [][]int) [][]int {
	difference := make([][]int, 0)
	for _, slice := range(a) {
		if !resultSetContains(b, slice) {
			difference = append(difference, slice)
		}
	}
	return difference
}

type binaryOperator func(a, b int, width uint) int;

func testBinaryOperator(t *testing.T, width uint, a int, b int, c int, problem bitvecsat.Problem, operator binaryOperator) {
	formula := problem.MakeSatFormula()
	t.Logf("Formula: %v", formula)
	forbidders := make([]sat.Clause, 0)

	// Contains 3-tuples of (A, B, C)
	foundSolutions := make([][]int, 0)

	for {
		formula.Clauses = append(formula.Clauses, forbidders...)
		solution := solve(formula)
		if solution == nil {
			break
		}

		aValue := problem.GetValueInAssignment(solution, a)
		bValue := problem.GetValueInAssignment(solution, b)
		cValue := problem.GetValueInAssignment(solution, c)

		if true {
			aString := problem.GetBitsInAssignment(solution, a)
			bString := problem.GetBitsInAssignment(solution, b)
			cString := problem.GetBitsInAssignment(solution, c)

			t.Log("A=" + strconv.Itoa(aValue) + "=" + aString + " B=" + strconv.Itoa(bValue) + "=" + bString + " C=" + strconv.Itoa(cValue) + "=" + cString);
		}

		foundSolutions = append(foundSolutions, []int{aValue, bValue, cValue})
		forbidders = append(forbidders, solution.MakeForbiddingClause())
	}

	expectedSolutions := make([][]int, 0)
	for a := 0; a < (1 << width); a++ {
		for b := 0; b < (1 << width); b++ {
			c := operator(a, b, width)
			if c >= 0 {
				// c < 0 means "outside defined inputs"
				expectedSolutions = append(expectedSolutions, []int{a, b, c})
			}
		}
	}

	if !resultSetsEqual(foundSolutions, expectedSolutions) {
		t.Errorf("unexpected solutions: found %v, expected %v", foundSolutions, expectedSolutions)
		extra := getResultsDifference(foundSolutions, expectedSolutions)
		if len(extra) > 0 {
			t.Errorf("extra: %v", extra)
		}
		missing := getResultsDifference(expectedSolutions, foundSolutions)
		if len(missing) > 0 {
			t.Errorf("missing: %v", missing)
		}
	}
}

func operatorPlus(a int, b int, width uint) int {
	return (a + b) % (1 << width);
}

func TestAddition(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		plus_constrain := bitvecsat.PlusConstrain{AIndex: a, BIndex: b, SumIndex: c}
		plus_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorPlus)
	}
}

func operatorOr(a int, b int, width uint) int {
	return a | b;
}

func TestOr(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		or_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.OrConstrain}
		or_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorOr)
	}
}

func operatorAnd(a int, b int, width uint) int {
	return a & b;
}

func TestAnd(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		and_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.AndConstrain}
		and_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorAnd)
	}
}

func operatorXor(a int, b int, width uint) int {
	return a ^ b;
}

func TestXor(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		xor_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.XorConstrain}
		xor_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorXor)
	}
}

func operatorEquiv(a int, b int, width uint) int {
	return (1 << width - 1) - (a ^ b);
}

func TestEquiv(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		equiv_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.EquivConstrain}
		equiv_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorEquiv)
	}
}

type ternaryRelation func(a, b, c int, width uint) bool;

func testTernaryRelation(t *testing.T, width uint, a, b, c int, problem bitvecsat.Problem, relation ternaryRelation) {
	formula := problem.MakeSatFormula()
	forbidders := make([]sat.Clause, 0)

	// Contains tuples of (A, B, C)
	foundSolutions := make([][]int, 0)

	for {
		formula.Clauses = append(formula.Clauses, forbidders...)
		t.Log(formula)
		solution := solve(formula)

		if solution == nil {
			break
		}

		aValue := problem.GetValueInAssignment(solution, a)
		bValue := problem.GetValueInAssignment(solution, b)
		cValue := problem.GetValueInAssignment(solution, c)

		t.Log(solution)
		fmt.Println(aValue, bValue, cValue)

		foundSolutions = append(foundSolutions, []int{aValue, bValue, cValue})
		forbidders = append(forbidders, solution.MakeForbiddingClause())
	}

	expectedSolutions := make([][]int, 0)
	for a := 0; a < (1 << width); a++ {
		for b := 0; b < (1 << width); b++ {
			for c := 0; c < (1 << width); c++ {
				if relation(a, b, c, width) {
					expectedSolutions = append(expectedSolutions, []int{a, b, c})
				}
			}
		}
	}

	if !resultSetsEqual(foundSolutions, expectedSolutions) {
		t.Fatalf("unexpected solutions: found %v, expected %v", foundSolutions, expectedSolutions)
	}
}

type binaryRelation func(a, b int, width uint) bool;

func testBinaryRelation(t *testing.T, width uint, a int, b int, problem bitvecsat.Problem, relation binaryRelation) {
	formula := problem.MakeSatFormula()
	forbidders := make([]sat.Clause, 0)

	// Contains tuples of (A, B)
	foundSolutions := make([][]int, 0)

	for {
		formula.Clauses = append(formula.Clauses, forbidders...)
		solution := solve(formula)
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

type unaryRelation func(a int, width uint) bool;

func testUnaryRelation(t *testing.T, width uint, a int, problem bitvecsat.Problem, relation unaryRelation) {
	formula := problem.MakeSatFormula()
	forbidders := make([]sat.Clause, 0)

	foundSolutions := make([][]int, 0)

	for {
		formula.Clauses = append(formula.Clauses, forbidders...)
		solution := solve(formula)
		if solution == nil {
			break
		}

		aValue := problem.GetValueInAssignment(solution, a)

		foundSolutions = append(foundSolutions, []int{aValue})
		forbidders = append(forbidders, solution.MakeForbiddingClause())
	}

	expectedSolutions := make([][]int, 0)
	for a := 0; a < (1 << width); a++ {
		if relation(a, width) {
			expectedSolutions = append(expectedSolutions, []int{a})
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
	for width := uint(1); width < 4; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)

		lte_constrain := bitvecsat.OrderingConstrain{AIndex: a, BIndex: b, Type: bitvecsat.LTE}
		lte_constrain.AddToProblem(&problem)

		testBinaryRelation(t, width, a, b, problem, relationLte)
	}
}

func relationLt(a int, b int, width uint) bool {
	return a < b;
}

func TestLt(t *testing.T) {
	for width := uint(1); width < 4; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)

		lte_constrain := bitvecsat.OrderingConstrain{AIndex: a, BIndex: b, Type: bitvecsat.LT}
		lte_constrain.AddToProblem(&problem)

		testBinaryRelation(t, width, a, b, problem, relationLt)
	}
}

func operatorMultiply(a int, b int, width uint) int {
	return (a * b) % (1 << width);
}

func TestMultiplication(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		mutliply_constrain := bitvecsat.MultiplyConstrain{AIndex: a, BIndex: b, ProductIndex: c}
		mutliply_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorMultiply)
	}
}

func TestLiteral(t *testing.T) {
	width := uint(8)
	problem := bitvecsat.Problem{}
	a := problem.AddNewVector(width)

	constrain := bitvecsat.LiteralConstrain{AIndex: a, Value: 193}
	constrain.AddToProblem(&problem)

	equals93 := func(a int, width uint) bool {
		return a == 193;
	}

	testUnaryRelation(t, width, a, problem, equals93)
}

func TestDivision(t *testing.T) {
	relationDivide := func(a, b, c int, width uint) bool {
		return (b != 0) && (a / b == c);
	}
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		ratio := problem.AddNewVector(width)

		c := bitvecsat.DivideConstrain{AIndex: a, BIndex: b, RatioIndex: ratio}
		c.AddToProblem(&problem)

		testTernaryRelation(t, width, a, b, ratio, problem, relationDivide)
	}
}

func TestModulo(t *testing.T) {
	relationModulo := func(a, b, c int, width uint) bool {
		return (b != 0) && (a % b == c);
	}
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		remainder := problem.AddNewVector(width)

		c := bitvecsat.DivideConstrain{AIndex: a, BIndex: b, RemainderIndex: remainder}
		c.AddToProblem(&problem)

		testTernaryRelation(t, width, a, b, remainder, problem, relationModulo)
	}
}

func TestShiftLeft(t *testing.T) {
	shift := uint(2)
	width := uint(1 << shift)

	operatorShiftLeft := func(a int, b int, width uint) int {
		// Upper bits of actualShift are ignored, as on x86 CPUs
		// or in the C language.
		actualShift := (uint(b) % width)
		return (a << actualShift) % (1 << width);
	}

	problem := bitvecsat.Problem{}
	a := problem.AddNewVector(width)
	b := problem.AddNewVector(width)
	c := problem.AddNewVector(width)

	shift_constrain := bitvecsat.ShiftLeftConstrain{AIndex: a, AmountIndex: b, YIndex: c}
	shift_constrain.AddToProblem(&problem)

	testBinaryOperator(t, width, a, b, c, problem, operatorShiftLeft)
}

func TestShiftRight(t *testing.T) {
	shift := uint(2)
	width := uint(1 << shift)

	operatorShiftRight := func(a int, b int, width uint) int {
		actualShift := (uint(b) % width)
		return (a >> actualShift) % (1 << width);
	}

	problem := bitvecsat.Problem{}
	a := problem.AddNewVector(width)
	b := problem.AddNewVector(width)
	c := problem.AddNewVector(width)

	shift_constrain := bitvecsat.ShiftRightConstrain{AIndex: a, AmountIndex: b, YIndex: c}
	shift_constrain.AddToProblem(&problem)

	testBinaryOperator(t, width, a, b, c, problem, operatorShiftRight)
}
