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

func BitIfThenElse(result int, condition int, ifTrue int, ifFalse int) []Clause {
	// TODO: optimize
	clauses := make([]Clause, 0)
	for _, resultValue := range([]bool{false, true}) {
		for _, conditionValue := range([]bool{false, true}) {
			for _, ifTrueValue := range([]bool{false, true}) {
				for _, ifFalseValue := range([]bool{false, true}) {
					var expected bool
					if conditionValue {
						expected = ifTrueValue
					} else {
						expected = ifFalseValue
					}
					isOk := (resultValue == expected)
					if !isOk {
						clauses = append(clauses, NewClause(!resultValue, result, !conditionValue, condition, !ifTrueValue, ifTrue, !ifFalseValue, ifFalse))
					}
				}
			}
		}
	}
	return clauses
	/*
	return []Clause{
		NewClause(true, result, true, condition, true, ifTrue, false, ifFalse),
		NewClause(true, result, true, condition, false, ifTrue, false, ifFalse),
		NewClause(true, result, false, condition, true, ifTrue, false, ifFalse),
		NewClause(true, result, false, condition, false, ifTrue, true, ifFalse),
		NewClause(true, result, false, condition, false, ifTrue, false, ifFalse),
		NewClause(false, result, true, condition, true, ifTrue, true, ifFalse),
		NewClause(false, result, true, condition, true, ifTrue, false, ifFalse),
		NewClause(false, result, false, condition, true, ifTrue, true, ifFalse),
		NewClause(false, result, false, condition, true, ifTrue, false, ifFalse),
	}
	*/
}
