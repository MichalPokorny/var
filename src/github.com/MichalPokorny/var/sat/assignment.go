package sat

type Assignment []bool

type PartialAssignment struct {
	Values []bool
	Assigned []bool
}

func MakeEmptyAssignment(formula Formula) PartialAssignment {
	varCount := formula.CountVariables()
	var assignment PartialAssignment
	assignment.Values = make([]bool, varCount)
	assignment.Assigned = make([]bool, varCount)
	return assignment
}

func (self PartialAssignment) String() string {
	var s = ""
	for i := range self.Values {
		if self.Assigned[i] {
			if self.Values[i] {
				s += "1"
			} else {
				s += "0"
			}
		} else {
			s += "?"
		}
	}
	return s
}

func (assignment PartialAssignment) GetLiteralValue(literal Literal) (known bool, value bool) {
	if !assignment.Assigned[literal.Variable] {
		return false, false
	} else {
		return true, (literal.Positive == assignment.Values[literal.Variable])
	}
}

func (assignment PartialAssignment) GetClauseValue(clause Clause) (known bool, value bool) {
	var anyUnassigned = false
	for _, literal := range clause.Literals {
		assigned, value := assignment.GetLiteralValue(literal)
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

func (assignment PartialAssignment) GetFormulaValue(formula Formula) (known bool, value bool) {
	for _, clause := range formula.Clauses {
		assigned, value := assignment.GetClauseValue(clause)
		if !assigned {
			return false, false
		}
		if !value {
			return true, false
		}
	}
	return true, true
}


func (self Assignment) String() string {
	var s = ""
	for _, v := range self {
		if v {
			s += "1"
		} else {
			s += "0"
		}
	}
	return s
}

func (self Assignment) MakeForbiddingClause() Clause {
	literals := make([]Literal, len(self))
	for i := 0; i < len(self); i++ {
		literals[i] = Literal{
			Variable: i,
			Positive: !self[i],
		}
	}
	return Clause{Literals: literals}
}
