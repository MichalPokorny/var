package sat

type Assignment []bool

type PartialAssignment struct {
	Values []bool
	Assigned []bool
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
