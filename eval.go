package main

import "fmt"

var FalseObject = Boolean{value: false}
var TrueObject = Boolean{value: true}
var NullObject = Null{}

// --------------------------------------------------------------------------------------------------------------------
// Evaluate parsed ast nodes
// --------------------------------------------------------------------------------------------------------------------

func eval(node Node, env *Environment) Object {
	switch node := node.(type) {
	// Statements
	case *Program:
		return evalProgram(node, env)
	case *BlockStatement:
		return evalBlockStatement(node, env)
	case *ExpressionStatement:
		return eval(node.expression, env)
	case *LetStatement:
		val := eval(node.value, env)
		if isError(val) {
			return val
		}
		env.set(node.name.value, val)
		return val
	case *ReturnStatement:
		val := eval(node.value, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{value: val}

	// Expressions
	case *ArrayLiteral:
		elements := evalExpressions(node.elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{elements: elements}
	case *BooleanLiteral:
		return nativeBoolToBoolObj(node.value)
	case *CallExpression:
		function := eval(node.function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(node.token, function, args)
	case *FloatLiteral:
		return &Float{value: node.value}
	case *FunctionLiteral:
		params := node.parameters
		body := node.body
		return &Function{parameters: params, env: env, body: body}
	case *Identifier:
		return evalIdentifier(node.token, node, env)
	case *IfExpression:
		return evalIfExpression(node, env)
	case *IndexExpression:
		left := eval(node.left, env)
		if isError(left) {
			return left
		}
		index := eval(node.index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(node.token, left, index)
	case *IntegerLiteral:
		return &Integer{value: node.value}
	case *InfixExpression:
		left := eval(node.left, env)
		if isError(left) {
			return left
		}
		right := eval(node.right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpr(node.token, left, right, node.operator)
	case *PrefixExpression:
		right := eval(node.right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.token, node.operator, right)
	case *StringLiteral:
		return &StringValue{value: node.tokenLiteral()}
	default:
		return &NullObject
	}
}

// --------------------------------------------------------------------------------------------------------------------
// Eval different types
// --------------------------------------------------------------------------------------------------------------------

func evalProgram(program *Program, env *Environment) Object {
	var result Object

	for _, stmt := range program.statements {
		result = eval(stmt, env)

		switch result := result.(type) {
		case *ReturnValue, *Error:
			return result
		}
	}

	return result
}

// --------------------------------------------------------------------------------------------------------------------

func evalBlockStatement(block *BlockStatement, env *Environment) Object {
	var result Object

	for _, stmt := range block.statements {
		result = eval(stmt, env)

		if result != nil {
			returnType := result.Type()
			if returnType == RETURN_OBJ || returnType == ERR_OBJ {
				return result
			}
		}
	}

	return result
}

// --------------------------------------------------------------------------------------------------------------------
// Call expressions
// --------------------------------------------------------------------------------------------------------------------

func evalExpressions(exprs []Expression, env *Environment) []Object {
	var result []Object

	for _, expr := range exprs {
		evaluated := eval(expr, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// --------------------------------------------------------------------------------------------------------------------
// If exprs
// --------------------------------------------------------------------------------------------------------------------

func evalArrayIndexExpression(array, index Object) Object {
	arrayObject := array.(*Array)
	idx := index.(*Integer).value
	max := int64(len(arrayObject.elements) - 1)

	if idx < 0 || idx > max {
		return &NullObject
	}

	return arrayObject.elements[idx]
}

// --------------------------------------------------------------------------------------------------------------------

func evalIdentifier(token Token, node *Identifier, env *Environment) Object {
	if val, ok := env.get(node.value); ok {
		return val
	}

	if builtin, ok := builtins[node.value]; ok {
		return builtin
	}

	return newError(
		"identifier not found {%v}. Found on line: %v, column: %v.",
		node.value,
		token.line,
		token.column,
	)
}

// --------------------------------------------------------------------------------------------------------------------

func evalIndexExpression(token Token, left, index Object) Object {
	switch {
	case left.Type() == ARRAY_OBJ && index.Type() == INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError(
			"index operator not supported: %v. On line %v, column: %v.",
			left.Type(),
			token.line,
			token.column,
		)
	}
}

// --------------------------------------------------------------------------------------------------------------------

func evalIfExpression(ifExpr *IfExpression, env *Environment) Object {
	condition := eval(ifExpr.condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return eval(ifExpr.consequence, env)
	} else if ifExpr.alternative != nil {
		return eval(ifExpr.alternative, env)
	} else {
		return &NullObject
	}
}

// --------------------------------------------------------------------------------------------------------------------

func isTruthy(object Object) bool {
	switch object {
	case &NullObject:
		return false
	case &TrueObject:
		return true
	case &FalseObject:
		return false
	default:
		return true
	}
}

// --------------------------------------------------------------------------------------------------------------------
// Infix exprs
// --------------------------------------------------------------------------------------------------------------------

func evalInfixExpr(token Token, left, right Object, operator string) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpr(token, left, right, operator)
	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		return evalFloatInfixExpr(token, left, right, operator)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(token, operator, left, right)
	case operator == "==":
		return nativeBoolToBoolObj(left == right)
	case operator == "!=":
		return nativeBoolToBoolObj(left != right)
	case left.Type() != right.Type():
		return newError(
			"mismatched types found when evaluating infix expression {%v, %v}. On line: %v, column: %v",
			left.Type(),
			right.Type(),
			token.line,
			token.column,
		)
	default:
		return newError(
			"invalid operator found when evaluating infix expression {%v %v %v}. On line %v, column %v.",
			left.Type(),
			operator,
			right.Type(),
			token.line,
			token.column,
		)
	}
}

// --------------------------------------------------------------------------------------------------------------------

func evalFloatInfixExpr(token Token, left, right Object, operator string) Object {
	leftVal := left.(*Float).value
	rightVal := right.(*Float).value

	switch operator {
	case "+":
		return &Float{value: leftVal + rightVal}
	case "-":
		return &Float{value: leftVal - rightVal}
	case "/":
		return &Float{value: leftVal / rightVal}
	case "*":
		return &Float{value: leftVal * rightVal}
	default:
		return evalOtherInfixOperators(token, leftVal, rightVal, operator)
	}
}

// --------------------------------------------------------------------------------------------------------------------

func evalIntegerInfixExpr(token Token, left, right Object, operator string) Object {
	leftVal := left.(*Integer).value
	rightVal := right.(*Integer).value

	switch operator {
	case "+":
		return &Integer{value: leftVal + rightVal}
	case "-":
		return &Integer{value: leftVal - rightVal}
	case "/":
		return &Integer{value: leftVal / rightVal}
	case "*":
		return &Integer{value: leftVal * rightVal}
	default:
		return evalOtherInfixOperators(token, leftVal, rightVal, operator)
	}
}

// --------------------------------------------------------------------------------------------------------------------

func evalOtherInfixOperators[T int64 | float64](token Token, left, right T, operator string) Object {
	switch operator {
	case "<":
		return nativeBoolToBoolObj(left < right)
	case "<=":
		return nativeBoolToBoolObj(left <= right)
	case ">":
		return nativeBoolToBoolObj(left > right)
	case ">=":
		return nativeBoolToBoolObj(left >= right)
	case "!=":
		return nativeBoolToBoolObj(left != right)
	case "==":
		return nativeBoolToBoolObj(left == right)
	default:
		return newError(
			"invalid operator found when evaluating infix expression {%v %v %v}. On line %v, column %v.",
			left,
			operator,
			right,
			token.line,
			token.column,
		)
	}
}

// --------------------------------------------------------------------------------------------------------------------

func evalStringInfixExpression(token Token, operator string, left, right Object) Object {
	if operator != "+" {
		return newError(
			"invalid operator found when evaluating infix expression {%v %v %v}. On line %v, column %v.",
			left.Type(),
			operator,
			right.Type(),
			token.line,
			token.column,
		)
	}

	leftVal := left.(*StringValue).value
	rightVal := right.(*StringValue).value

	return &StringValue{value: leftVal + rightVal}
}

// --------------------------------------------------------------------------------------------------------------------
// Prefix exprs
// --------------------------------------------------------------------------------------------------------------------

func evalPrefixExpression(token Token, operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpr(right)
	case "-":
		return evalMinusPrefixExpr(token, right)
	default:
		return newError(
			"invalid operator in prefix position {%v}. On line: %v, column: %v",
			operator,
			token.line,
			token.column,
		)
	}
}

// --------------------------------------------------------------------------------------------------------------------

func evalBangOperatorExpr(right Object) Object {
	switch right {
	case &TrueObject:
		return &FalseObject
	case &FalseObject:
		return &TrueObject
	case &NullObject:
		return &TrueObject
	default:
		return &FalseObject
	}
}

// --------------------------------------------------------------------------------------------------------------------

func evalMinusPrefixExpr(token Token, right Object) Object {
	if right.Type() == INTEGER_OBJ {
		return &Integer{value: -right.(*Integer).value}
	}
	if right.Type() == FLOAT_OBJ {
		return &Float{value: -right.(*Float).value}
	}

	return newError(
		"invalid operator {%v}. On line: %v, column: %v.",
		token.literal,
		token.line,
		token.column,
	)
}

// --------------------------------------------------------------------------------------------------------------------
// Helpers
// --------------------------------------------------------------------------------------------------------------------

func applyFunction(token Token, fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := eval(fn.body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *Builtin:
		return fn.fn(args...)
	default:
		return newError(
			"attempted to eval something that wasn't a function {%v}. On line: %v, column: %v.",
			fn.Type(),
			token.line,
			token.column,
		)
	}
}

// --------------------------------------------------------------------------------------------------------------------

func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := newEnclosedEnvironment(fn.env)

	for idx, param := range fn.parameters {
		env.set(param.value, args[idx])
	}

	return env
}

// --------------------------------------------------------------------------------------------------------------------

func unwrapReturnValue(object Object) Object {
	if returnValue, ok := object.(*ReturnValue); ok {
		return returnValue.value
	}

	return object
}

// --------------------------------------------------------------------------------------------------------------------

func isError(object Object) bool {
	if object != nil {
		return object.Type() == ERR_OBJ
	}

	return false
}

// --------------------------------------------------------------------------------------------------------------------

func nativeBoolToBoolObj(input bool) *Boolean {
	if input {
		return &TrueObject
	}
	return &FalseObject
}

// --------------------------------------------------------------------------------------------------------------------

func newError(format string, vars ...interface{}) *Error {
	return &Error{message: fmt.Sprintf(format, vars...)}
}

// --------------------------------------------------------------------------------------------------------------------
