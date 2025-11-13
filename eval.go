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
		return evalProgram(node)
	case *ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *IntegerLiteral:
		return &Integer{Value: node.Value}
	case *BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)
	case *PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *BlockStatement:
		return evalBlockStatement(node)
	case *IfExpression:
		return evalIfExpression(node)
	case *ReturnStatement:
		val := Eval(node.ReturnValue)
		return &ReturnValue{Value: val}
	}

	return nil
}

func evalProgram(program *Program) Object {
	var result Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalBlockStatement(block *BlockStatement) Object {
	var result Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

		if result != nil && result.Type() == RETURN_VALUE_OBJ {
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return O_TRUE
	}
	return O_FALSE
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return O_NULL
	}
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return O_NULL
	}
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case O_TRUE:
		return O_FALSE
	case O_FALSE:
		return O_TRUE
	case O_NULL:
		return O_TRUE
	default:
		return O_FALSE
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return O_NULL
	}
	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		return &Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return O_NULL
	}
}

func evalIfExpression(ie *IfExpression) Object {
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return O_NULL
	}
}

func isTruthy(obj Object) bool {
	switch obj {
	case O_NULL:
		return false
	case O_TRUE:
		return true
	case O_FALSE:
		return false
	default:
		return true
	}
}
