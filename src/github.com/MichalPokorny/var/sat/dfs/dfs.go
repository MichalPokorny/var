package dfs

import "github.com/MichalPokorny/var/sat"


// TODO: We should at least do DPLL
func Solve(formula sat.Formula) sat.Assignment {
	assignment := sat.MakeEmptyAssignment(formula)

	var firstToChange = 0
	for {
		if !assignment.Assigned[firstToChange] {
			assignment.Assigned[firstToChange] = true
			assignment.Values[firstToChange] = false
		} else {
			if !assignment.Values[firstToChange] {
				assignment.Values[firstToChange] = true
			} else {
				// does not work, backtrack
				assignment.Assigned[firstToChange] = false
				firstToChange -= 1
				if firstToChange < 0 {
					return nil
				}
				continue
			}
		}

		known, result := assignment.GetFormulaValue(formula)
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
