package cdcl

import (
	"testing"
	"github.com/MichalPokorny/var/sat"
)

func TestClauseState(t *testing.T) {
	assignment := partialAssignment{
		Values: []bool{false, false, false, true, true},
		Times: []int  {   -1,     1,     2,    3,   -1},
		Assigned: []bool{false, true, true, true, false},
	}
	clause := sat.NewClause(true, 2, true, 0, true, 1)
	if assignment.getClauseState(clause) != CLAUSE_UNRESOLVED {
		t.Error("clause should be unresolved")
	}
	clause = sat.NewClause(true, 1, false, 2, false, 3)
	if assignment.getClauseState(clause) != CLAUSE_SATISFIED {
		t.Error("clause should be satisfied")
	}
	clause = sat.NewClause(true, 1, true, 2, false, 3)
	if assignment.getClauseState(clause) != CLAUSE_CONFLICTING {
		t.Error("clause should be conflicting")
	}
}

func TestFindDirectlyImpliedLiteral(t *testing.T) {
	assignment := partialAssignment{
		Values: []bool{false, false, false, true, false},
		Times: []int  {   -1,     1,     2,    3,   -1},
		Assigned: []bool{false, true, true, true, false},
	}
	clause := sat.NewClause(true, 1, false, 2, false, 3)
	if assignment.findDirectlyImpliedLiteral(clause) != nil {
		t.Error("should not have implied literal, clause is true")
	}

	clause = sat.NewClause(true, 0, false, 2, false, 4)
	if assignment.findDirectlyImpliedLiteral(clause) != nil {
		t.Error("should not have implied literal, clause has two unassigned")
	}

	clause = sat.NewClause(true, 0, true, 2, false, 3)
	literal := assignment.findDirectlyImpliedLiteral(clause)
	if literal == nil || literal.Variable != 0 || literal.Positive != true {
		t.Error("failed to find implied literal")
	}

	clause = sat.NewClause(false, 0, true, 2, false, 3)
	literal = assignment.findDirectlyImpliedLiteral(clause)
	if literal == nil || literal.Variable != 0 || literal.Positive != false {
		t.Error("failed to find implied literal in", clause)
	}
}
