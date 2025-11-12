package main

func Eval(node Node) Object {
	switch node := node.(type) {
	// Statements
	case *Program:
		return evalStatements(node.Statements)
	case *ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *IntegerLiteral:
		return &Integer{Value: node.Value}
	}

	return nil
}

func evalStatements(statements []Statement) Object {
	var result Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}
