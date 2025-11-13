package main

import "fmt"

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
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *BlockStatement:
		return evalBlockStatement(node)
	case *IfExpression:
		return evalIfExpression(node)
	case *ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	}

	return nil
}

func evalProgram(program *Program) Object {
	var result Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *BlockStatement) Object {
	var result Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
				return result
			}
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
		return newError("unknown operator: %s%s", operator, right.Type())
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
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: -%s", right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *IfExpression) Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}
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

func newError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}
