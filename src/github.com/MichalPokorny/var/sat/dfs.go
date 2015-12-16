package sat

import (
//	"fmt"
//	"strconv"
)

func countVariables(formula Formula) int {
	var maxVariable = 0
	for _, clause := range formula.Clauses {
		for _, literal := range clause.Literals {
			if literal.Variable > maxVariable {
				maxVariable = literal.Variable
			}
		}
	}
	return maxVariable + 1
}

func getLiteralValue(literal Literal, assignment PartialAssignment) (known bool, value bool) {
	if !assignment.Assigned[literal.Variable] {
		return false, false
	} else {
		return true, (literal.Positive == assignment.Values[literal.Variable])
	}
}

func getClauseValue(clause Clause, assignment PartialAssignment) (known bool, value bool) {
	var anyUnassigned = false
	for _, literal := range clause.Literals {
		assigned, value := getLiteralValue(literal, assignment)
		if !assigned {
			anyUnassigned = true
		}
		if value {
			return true, true
		}
	}
	if anyUnassigned {
		return false, false
	} else {
		return true, false
	}
}

func getFormulaValue(formula Formula, assignment PartialAssignment) (known bool, value bool) {
	for _, clause := range formula.Clauses {
		assigned, value := getClauseValue(clause, assignment)
		if !assigned {
			// fmt.Println("formula:" + formula.String() + " assgn:" + assignment.String() + " => don't know")
			return false, false
		}
		if !value {
			// fmt.Println("formula:" + formula.String() + " assgn:" + assignment.String() + " => false")
			return true, false
		}
	}
	// fmt.Println("formula:" + formula.String() + " assgn:" + assignment.String() + " => true")
	return true, true
}

// TODO: We should at least do DPLL
func Solve(formula Formula) Assignment {
	varCount := countVariables(formula)
	var assignment PartialAssignment
	assignment.Values = make([]bool, varCount)
	assignment.Assigned = make([]bool, varCount)

	var firstToChange = 0
	for {
		// fmt.Println("firstToChange=" + strconv.Itoa(firstToChange) + ", current assignment: " + assignment.String())
		if !assignment.Assigned[firstToChange] {
			assignment.Assigned[firstToChange] = true
			assignment.Values[firstToChange] = false
		} else {
			if !assignment.Values[firstToChange] {
				assignment.Values[firstToChange] = true
			} else {
				// does not work, backtrack
				// fmt.Println("backtrack")
				assignment.Assigned[firstToChange] = false
				firstToChange -= 1
				if firstToChange < 0 {
					return nil
				}
				continue
			}
		}

		known, result := getFormulaValue(formula, assignment)
		if known {
			if result {
				return assignment.Values
			} else {
				// does not work, will backtrack
			}
		} else {
			firstToChange += 1
		}
	}
	return nil
}
