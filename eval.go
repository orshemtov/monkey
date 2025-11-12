package main

var (
	O_NULL  = &Null{}
	O_TRUE  = &Boolean{Value: true}
	O_FALSE = &Boolean{Value: false}
)

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
	case *BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)
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

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return O_TRUE
	}
	return O_FALSE
}
