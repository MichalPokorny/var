package sat

// TODO: auto-derive?

// A xor B = Y
// TODO: optimize?
func XorConstrain(a int, b int, y int) []Clause {
	return []Clause{
		NewClause(true, a, true, b, false, y),
		NewClause(true, a, false, b, true, y),
		NewClause(false, a, true, b, true, y),
		NewClause(false, a, false, b, false, y),
	}
}

// (A <=> B) = Y
// TODO: optimize? (NOTE: <=> not XOR)
func EquivConstrain(a int, b int, y int) []Clause {
	return []Clause{
		NewClause(true, a, true, b, true, y),
		NewClause(true, a, false, b, false, y),
		NewClause(false, a, true, b, false, y),
		NewClause(false, a, false, b, true, y),
	}
}

// A and B = Y
// TODO: optimize?
func AndConstrain(a int, b int, y int) []Clause {
	return []Clause{
		NewClause(true, a, true, b, false, y),
		NewClause(true, a, false, b, false, y),
		NewClause(false, a, true, b, false, y),
		NewClause(false, a, false, b, true, y),
	}
}

// A or B = Y
// TODO: optimize?
func OrConstrain(a int, b int, y int) []Clause {
	return []Clause{
		NewClause(true, a, true, b, false, y),
		NewClause(true, a, false, b, true, y),
		NewClause(false, a, true, b, true, y),
		NewClause(false, a, false, b, true, y),
	}
}

// (A < B) = Y
// TODO: optimize?
func LtConstrain(a int, b int, y int) []Clause {
	return []Clause{
		NewClause(true, a, true, b, false, y),
		NewClause(true, a, false, b, true, y),
		NewClause(false, a, true, b, false, y),
		NewClause(false, a, false, b, false, y),
	}
}

// TODO: test
// TODO: optimize
func BitsAlwaysEqual(a int, b int) []Clause {
	return []Clause{
		NewClause(true, a, false, b),
		NewClause(false, a, true, b),
	}
}

func BitIsTrue(a int) []Clause {
	return []Clause{NewClause(true, a)}
}

func BitIsFalse(a int) []Clause {
	return []Clause{NewClause(false, a)}
}
