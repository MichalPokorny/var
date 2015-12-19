package dpll

import (
	//"fmt"
	"github.com/MichalPokorny/var/sat"
)

func onlyUnassignedLiteral(clause sat.Clause, assignment sat.PartialAssignment) (found bool, literal sat.Literal) {
	var unassignedLiteral sat.Literal
	found = false

	for _, literal := range(clause.Literals) {
		known, _ := assignment.GetLiteralValue(literal)
		if !known {
			// TODO: and unassignedLiteral != literal...
			if found {
				return false, literal
			} else {
				found = true
				unassignedLiteral = literal
			}
		}
	}
	return found, unassignedLiteral
}

func Prune(formula sat.Formula, assignment sat.PartialAssignment) (possible bool, result sat.Formula) {
	result.Clauses = make([]sat.Clause, 0)
	for _, clause := range(formula.Clauses) {
		if known, satisfied := assignment.GetClauseValue(clause); known {
			if satisfied {
				continue  // do not add this clause
			} else {
				return false, result
			}
		}

		newClause := sat.Clause{Literals: make([]sat.Literal, 0)}
		for _, literal := range(clause.Literals) {
			if known, result := assignment.GetLiteralValue(literal); known {
				if result {
					panic("clause satisfied")
				} else {
					// trim this literal, it's not satisfied
					continue
				}
			}
			newClause.Literals = append(newClause.Literals, literal)
		}
		result.Clauses = append(result.Clauses, newClause)
	}
	return true, result
}

func CopyAssignment(assignment sat.PartialAssignment) sat.PartialAssignment {
	return sat.PartialAssignment{
		Values: append([]bool(nil), assignment.Values...),
		Assigned: append([]bool(nil), assignment.Assigned...),
	}
}

func Solve(formula sat.Formula, assignment sat.PartialAssignment) sat.Assignment {
	varCount := formula.CountVariables()

	//fmt.Println("DPLL(", formula, "):")
	//fmt.Println(assignment)

	known, result := assignment.GetFormulaValue(formula)
	if known {
		if result {
			//fmt.Println("value known (true)")
			return assignment.Values
		} else {
			//fmt.Println("value known (false)")
			return nil
		}
	}
	//fmt.Println("value unknown")

	// TODO: Assign pure variables.

	// Propagate unit clauses and remove true clauses.
	changed := true
	for changed {
		changed = false
		for _, clause := range(formula.Clauses) {
			found, literal := onlyUnassignedLiteral(clause, assignment)
			if found {
				assignment = CopyAssignment(assignment)
				assignment.Assigned[literal.Variable] = true
				assignment.Values[literal.Variable] = literal.Positive
				//fmt.Println("clause", clause, "implies", literal)
				changed = true
				// Remove the finished clause.
				var possible bool
				possible, formula = Prune(formula, assignment)
				if !possible {
					//fmt.Println("result now false")
					return nil
				}
				if known, result := assignment.GetFormulaValue(formula); known {
					if result {
						//fmt.Println("result now true")
						return assignment.Values
					} else {
						//fmt.Println("result now false")
						return nil
					}
				}
				break
			}
		}
	}

	// Select some variable to assign.
	var varToAssign int
	for i := 0; i < varCount; i++ {
		if !assignment.Assigned[i] {
			varToAssign = i
			break
		}
	}
	assignment = CopyAssignment(assignment)
	assignment.Assigned[varToAssign] = true
	assignment.Values[varToAssign] = true

	//fmt.Println("try", varToAssign, "<- true")
	if possible, prunedFormula := Prune(formula, assignment); possible {
		positiveAssignment := Solve(prunedFormula, assignment)
		if positiveAssignment != nil {
			return positiveAssignment
		}
	} else {
		//fmt.Println("pruning not possible")
	}

	//fmt.Println("try", varToAssign, "<- false")
	assignment = CopyAssignment(assignment)
	assignment.Assigned[varToAssign] = true
	assignment.Values[varToAssign] = false

	if possible, prunedFormula := Prune(formula, assignment); possible {
		negativeAssignment := Solve(prunedFormula, assignment)
		if negativeAssignment != nil {
			return negativeAssignment
		}
	}
	assignment.Assigned[varToAssign] = false

	return nil
}
