package main

import "fmt"

// --------------------------------------------------------------------------------------------------------------------
// Language builtins
// --------------------------------------------------------------------------------------------------------------------

var builtins = map[string]*Builtin{
	"len": {fn: func(args ...Object) Object {
		if len(args) != 1 {
			return newError("wrong number of args passed to len(). Got %v, want 1", len(args))
		}
		switch arg := args[0].(type) {
		case *Array:
			return &Integer{value: int64(len(arg.elements))}
		case *StringValue:
			return &Integer{value: int64(len(arg.value))}
		default:
			return newError("argument to 'len()' not supported, got %v", args[0].Type())
		}
	},
	},
	"first": {fn: func(args ...Object) Object {
		if len(args) != 1 {
			return newError("first: wrong number of arguments. Got %v, want 1", len(args))
		}
		if args[0].Type() != ARRAY_OBJ {
			return newError("first: arguments to first must be an Array, got %v.", args[0].Type())
		}
		arr := args[0].(*Array)
		if len(arr.elements) > 0 {
			return arr.elements[0]
		}

		return &NullObject
	},
	},
	"last": {fn: func(args ...Object) Object {
		if len(args) != 1 {
			return newError("last: wrong number of arguments. Got %v, want 1", len(args))
		}
		if args[0].Type() != ARRAY_OBJ {
			return newError("last: arguments to last must be an Array, got %v.", args[0].Type())
		}

		arr := args[0].(*Array)
		length := len(arr.elements)
		if length > 0 {
			return arr.elements[length-1]
		}

		return &NullObject
	},
	},
	"rest": {fn: func(args ...Object) Object {
		if len(args) != 1 {
			return newError("rest: wrong number of arguments. Got %v, want 1", len(args))
		}
		if args[0].Type() != ARRAY_OBJ {
			return newError("rest: arguments to rest must be an Array, got %v.", args[0].Type())
		}

		arr := args[0].(*Array)
		length := len(arr.elements)
		if length > 0 {
			newElements := make([]Object, length-1, length-1)
			copy(newElements, arr.elements[1:length])
			return &Array{elements: newElements}
		}

		return &NullObject
	},
	},
	"push": {fn: func(args ...Object) Object {
		if len(args) != 2 {
			return newError("push: wrong number of arguments. Got %v, want 2", len(args))
		}
		if args[0].Type() != ARRAY_OBJ {
			return newError("push: first argument to push must be an Array, got %v.", args[0].Type())
		}

		arr := args[0].(*Array)
		length := len(arr.elements)

		newElements := make([]Object, length+1, length+1)
		copy(newElements, arr.elements)
		newElements[length] = args[1]

		return &Array{elements: newElements}
	},
	},
	"range_array": {fn: func(args ...Object) Object {
		if len(args) != 2 {
			return newError("range_array: wrong number of arguments. Got %v, want 2", len(args))
		}
		if args[0].Type() != args[1].Type() || args[0].Type() != INTEGER_OBJ {
			return newError("range_array: invalid types provided: (%v, %v). This function only accepts INTEGERS.", args[0].Type(), args[1].Type())
		}
		start := args[0].(*Integer).value
		stop := args[1].(*Integer).value
		arr := make([]Object, 0)

		for i := int64(start); i <= stop; i++ {
			arr = append(arr, &Integer{value: i})
		}

		return &Array{elements: arr}
	},
	},
	"puts": {fn: func(args ...Object) Object {
		for _, arg := range args {
			fmt.Println(arg.inspect())
		}
		return &NullObject
	},
	},
}

// --------------------------------------------------------------------------------------------------------------------
