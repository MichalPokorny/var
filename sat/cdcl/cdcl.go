package cdcl

import (
	"os"
	"strconv"
	"log"
	"sort"
	"github.com/MichalPokorny/var/sat"
)

type partialAssignment struct {
	Values []bool
	Times []int
	Assigned []bool
}

var logger *log.Logger

func initLogger() {
	logger = log.New(os.Stdout, "", log.Lshortfile)
}

const (
	CLAUSE_SATISFIED = iota
	CLAUSE_CONFLICTING = iota
	CLAUSE_UNRESOLVED = iota
)

type clauseState int

func (assignment *partialAssignment) getClauseState(clause sat.Clause) clauseState {
	foundUnresolved := false
	for _, literal := range clause.Literals {
		if !assignment.Assigned[literal.Variable] {
			foundUnresolved = true
			continue
		}
		if assignment.Values[literal.Variable] == literal.Positive {
			return CLAUSE_SATISFIED
		}
	}
	if foundUnresolved {
		return CLAUSE_UNRESOLVED
	}
	return CLAUSE_CONFLICTING
}

func (assignment *partialAssignment) findDirectlyImpliedLiteral(clause sat.Clause) *sat.Literal {
	// logger.Println("state of", clause, ":", assignment.getClauseState(clause))
	if assignment.getClauseState(clause) == CLAUSE_SATISFIED {
		return nil
	}
	var impliedLiteral *sat.Literal = nil
	for i, literal := range clause.Literals {
		if !assignment.Assigned[literal.Variable] {
			if impliedLiteral != nil {
				// Multiple unassigned literals.
				return nil
			}
			impliedLiteral = &clause.Literals[i]
		}
	}
	return impliedLiteral
}

func (assignment *partialAssignment) hasConflict(formula sat.Formula) bool {
	for _, clause := range formula.Clauses {
		if assignment.getClauseState(clause) == CLAUSE_CONFLICTING {
			return true
		}
	}
	return false
}

func (assignment *partialAssignment) selectStep() *sat.Literal {
	for i := 0; i < len(assignment.Values); i++ {
		if !assignment.Assigned[i] {
			return &sat.Literal{Variable: i, Positive: false}
		}
	}
	return nil
}

// antedescent clause: C unit and implies L => L implied, C is antedescent
// (if there were more: choose the one we used for unit propagation)

type implicationEdge struct {
	From int
	To int  // -1 represents conflict
		// (x,K) where x is in some conflict clause
	Clause *sat.Clause
}

type implicationVertex struct {
	IsDecision bool
}

// every cut separating decision vertices and conflict
// is a formula that implies conflict [take variable
// values assigned just before the cut]

type implicationGraph struct {
	// variables are implicit from partialAssignment

	vertices []implicationVertex  // only for variables
	edges []implicationEdge
}

func (graph implicationGraph) String() string {
	str := ""
	for i, vertex := range graph.vertices {
		str += strconv.Itoa(i) + ":";
		if vertex.IsDecision {
			str += "D"
		} else {
			str += "p"
		}
		str += " "
	}
	for _, edge := range graph.edges {
		str += strconv.Itoa(edge.From) + "->" + strconv.Itoa(edge.To) + "[" + edge.Clause.String() + "] "
	}
	return str
}

func (graph implicationGraph) getAntedescent(formula sat.Formula, variable int) sat.Clause {
	for i := len(graph.edges) - 1; i >= 0; i-- {
		if graph.edges[i].To == variable {
			return *graph.edges[i].Clause
		}
	}
	panic("no antedescent for " + strconv.Itoa(variable))
}

func booleanConstrainPropagation(formula sat.Formula, assignment *partialAssignment, graph *implicationGraph, time int) {
	continueSearch := true
	for continueSearch {
		continueSearch = false
		for i, clause := range formula.Clauses {
			if assignment.getClauseState(clause) == CLAUSE_CONFLICTING {
				// found conflict
				//logger.Println("bcp found conflict with", clause)
				for _, literal := range clause.Literals {
					edge := implicationEdge{
						From: literal.Variable,
						To: -1,
						Clause: &formula.Clauses[i],
					}
					graph.edges = append(graph.edges, edge)
				}
				return
			}
			literal := assignment.findDirectlyImpliedLiteral(clause)
			if literal != nil {
				// found implication
				for _, source := range clause.Literals {
					if source.Variable != literal.Variable {
						edge := implicationEdge{
							From: source.Variable,
							To: literal.Variable,
							Clause: &formula.Clauses[i],
						}
						// TODO: assert we aren't
						// erasing a decision
						graph.vertices[literal.Variable].IsDecision = false
						graph.edges = append(graph.edges, edge)
					}
				}
				//logger.Println("bcp:", clause, "implies", literal, "(time=", time, ")")
				if !assignment.Assigned[literal.Variable] {
					assignment.Times[literal.Variable] = time
					assignment.Assigned[literal.Variable] = true
					assignment.Values[literal.Variable] = literal.Positive
				}
				continueSearch = true
			}
		}
	}
}

func backtrack(assignment *partialAssignment, graph *implicationGraph, level int) {
	// TODO: figure out the level

	for i := 0; i < len(assignment.Times); i++ {
		if assignment.Times[i] >= level {
			assignment.Times[i] = -1
			assignment.Assigned[i] = false
		}
	}

	edges := make([]implicationEdge, 0)
	for _, edge := range graph.edges {
		if edge.To != -1 && !assignment.Assigned[edge.To] {
			edges = append(edges, edge)
		}
	}
	graph.edges = edges
}

func resolve(a sat.Clause, b sat.Clause, variable int) sat.Clause {
	y := sat.Clause{Literals: make([]sat.Literal, 0)}
	containsLiteral := func(x sat.Literal) bool {
		for _, l := range y.Literals {
			if l == x {
				return true
			}
		}
		return false
	}
	for _, l := range a.Literals {
		if l.Variable != variable && !containsLiteral(l) {
			y.Literals = append(y.Literals, l)
		}
	}
	for _, l := range b.Literals {
		if l.Variable != variable && !containsLiteral(l) {
			y.Literals = append(y.Literals, l)
		}
	}
	return y
}

func secondHighestOrMinusOne(x []int) int {
	sort.Ints(x)
	if len(x) < 2 {
		return -1
	}
	for i := len(x) - 2; i >= 0; i-- {
		if x[i + 1] != x[i] {
			return x[i]
		}
	}
	return -1
}

func analyzeConflict(level int, formula *sat.Formula, graph implicationGraph, assignment *partialAssignment) int {
	var clause sat.Clause
	found := false
	for _, c := range formula.Clauses {
		// TODO: how to pick conflicting clause? pick all?
		if assignment.getClauseState(c) == CLAUSE_CONFLICTING {
			clause = c
			found = true
			break
		}
	}
	if !found {
		panic("no conflict")
	}
	stopCriterion := func(c sat.Clause) bool {
		onLevel := 0
		for _, l := range c.Literals {
			if assignment.Times[l.Variable] == level {
				onLevel++
			}
		}
		//logger.Println("clauses on current level (", level, "):", onLevel)
		return onLevel == 1
	}
	mostRecentNondecisionLiteral := func(c sat.Clause) sat.Literal {
		found := false
		var mostRecent sat.Literal
		for _, l := range c.Literals {
			if graph.vertices[l.Variable].IsDecision {
				continue
			}
			if !found || assignment.Times[l.Variable] > assignment.Times[mostRecent.Variable] {
				found = true
				mostRecent = l
			}
		}
		if !found {
			panic("all literals are decision")
		}
		return mostRecent
	}
	//logger.Println("assignment:", assignment)
	//logger.Println("graph:", graph)
	//logger.Println("analyzing conflict starting with", clause)

	for !stopCriterion(clause) {
		mostRecent := mostRecentNondecisionLiteral(clause)
		//rlevel := assignment.Times[mostRecent.Variable]
		//logger.Println("resolving through", mostRecent, "(level=", rlevel, ")")
		antedescent := graph.getAntedescent(*formula, mostRecent.Variable)
		//logger.Println("antedescent:", antedescent)
		clause = resolve(clause, antedescent, mostRecent.Variable)
		//logger.Println("=>", clause)
	}

	//logger.Println("learning clause:", clause)
	formula.Clauses = append(formula.Clauses, clause)

	// TODO: backtrack to 2nd highest level

	levels := make([]int, 0)
	for _, l := range clause.Literals {
		levels = append(levels, assignment.Times[l.Variable])
	}
	return secondHighestOrMinusOne(levels)
	/*
	lowestLevel := assignment.Times[clause.Literals[0].Variable]
	for _, l := range clause.Literals {
		if assignment.Times[l.Variable] < lowestLevel {
			lowestLevel = assignment.Times[l.Variable]
		}
	}
	return lowestLevel - 1
	*/
}

func Solve(formula sat.Formula) sat.Assignment {
	initLogger()

	//logger.Println("cdcl on:", formula)

	n := formula.CountVariables()
	assignment := partialAssignment{
		Values: make([]bool, n),
		Times: make([]int, n),  /* -1 for each value */
		Assigned: make([]bool, n),
	}
	graph := implicationGraph{
		vertices: make([]implicationVertex, n),
		edges: make([]implicationEdge, 0),
	}
	for i := 0; i < n; i++ {
		assignment.Times[i] = -1
		assignment.Assigned[i] = false
		graph.vertices[i] = implicationVertex{IsDecision: false}
	}

	time := -1
	booleanConstrainPropagation(formula, &assignment, &graph, time)
	if assignment.hasConflict(formula) {
		//logger.Println("conflict after first bcp")
		return nil
	}
	time = 0

	for {
		step := assignment.selectStep()
		if step == nil {
			// good, no more variables to assign
			//logger.Println("no more vars to assign =>", assignment.Values)
			return assignment.Values
		}

		//logger.Println("time=", time, "assigning", step)
		//logger.Println(formula)

		assignment.Times[step.Variable] = time
		assignment.Assigned[step.Variable] = true
		assignment.Values[step.Variable] = step.Positive
		graph.vertices[step.Variable].IsDecision = true

		// time := number of decisions taken
		// variable was decided ==> Time[v] is how many decisions were taken before assignment
		// variable was inferred ==> Time[v] is time of the variable that triggered the inference

		booleanConstrainPropagation(formula, &assignment, &graph, time)

		for assignment.hasConflict(formula) {
		//	if time == 0 {
			if time < 0 {
				// Conflict on level 0.
				// Cannot backtrack.
				//logger.Print("conflict in formula, cannot satisfy")
				return nil
			}

			nextTime := analyzeConflict(time, &formula, graph, &assignment)
			//logger.Println("backtrack to level", nextTime)
			// logger.Println("TODO")
			// return nil
			/*
			if level < 0 {
				return nil
			} else {
				backtrack(&assignment, &graph, level)
				time := level
			}
			*/
			backtrack(&assignment, &graph, nextTime)
			time = nextTime
			booleanConstrainPropagation(formula, &assignment, &graph, time)
		}

		time++
	}
}

/*
func Solve(formula sat.Formula) sat.Assignment {
}
*/
