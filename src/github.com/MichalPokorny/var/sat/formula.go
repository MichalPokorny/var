package sat

import ("fmt")

type Literal struct {
	Variable int
	Positive bool
}

type Clause struct {
	Literals []Literal
}

// Call as:
//   true, 1, false, 2
func NewClause(args ...interface{}) Clause {
	if len(args) % 2 != 0 {
		panic("odd arguments")
	}
	literals := make([]Literal, len(args) / 2)
	for i := 0; i < len(args) / 2; i++ {
		literals[i] = Literal{Variable: args[(i * 2) + 1].(int), Positive: args[i * 2].(bool)}
	}
	return Clause{Literals: literals}
}

type Formula struct {
	Clauses []Clause
}

func (self Literal) String() string {
	if self.Positive {
		return fmt.Sprintf("%d", self.Variable)
	} else {
		return fmt.Sprintf("-%d", self.Variable)
	}
}

func (self Clause) String() string {
	var s = "("
	for idx, literal := range self.Literals {
		s += literal.String()
		if idx != len(self.Literals) - 1 {
			s += " | "
		}
	}
	return s + ")"
}

func (self Formula) String() string {
	var s = ""
	for idx, clause := range self.Clauses {
		s += clause.String()
		if idx != len(self.Clauses) - 1 {
			s += " & "
		}
	}
	return s
}
