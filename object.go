package main

import (
	"bytes"
	"fmt"
	"strings"
)

// --------------------------------------------------------------------------------------------------------------------
// Allowed objects
// --------------------------------------------------------------------------------------------------------------------

type ObjectType string

type BuiltinFunc func(args ...Object) Object

const (
	ARRAY_OBJ    = "ARRAY"
	BOOL_OBJ     = "BOOLEAN"
	BUILTIN_OBJ  = "BUILTIN"
	ERR_OBJ      = "ERROR_OBJ"
	FLOAT_OBJ    = "FLOAT"
	FUNCTION_OBJ = "FUNCTION"
	INTEGER_OBJ  = "INTEGER"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN_VALUE"
	STRING_OBJ   = "STRING_OBJ"
)

// --------------------------------------------------------------------------------------------------------------------
// Object constraint
// --------------------------------------------------------------------------------------------------------------------

type Object interface {
	Type() ObjectType
	inspect() string
}

// --------------------------------------------------------------------------------------------------------------------
// Objects
// --------------------------------------------------------------------------------------------------------------------

type Array struct {
	elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }

func (a *Array) inspect() string {
	var buffer bytes.Buffer
	elements := make([]string, 0)

	for _, elem := range a.elements {
		elements = append(elements, elem.inspect())
	}

	buffer.WriteString("[")
	buffer.WriteString(strings.Join(elements, ", "))
	buffer.WriteString("]")

	return buffer.String()
}

// --------------------------------------------------------------------------------------------------------------------

type Boolean struct {
	value bool
}

func (b *Boolean) Type() ObjectType { return BOOL_OBJ }

func (b *Boolean) inspect() string { return fmt.Sprintf("%v", b.value) }

// --------------------------------------------------------------------------------------------------------------------

type Builtin struct {
	fn BuiltinFunc
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

func (b *Builtin) inspect() string { return "builtin function" }

// --------------------------------------------------------------------------------------------------------------------

type Error struct {
	message string
}

func (e *Error) Type() ObjectType { return ERR_OBJ }

func (e *Error) inspect() string { return "Error: " + e.message }

// --------------------------------------------------------------------------------------------------------------------

type Float struct {
	value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }

func (f *Float) inspect() string { return fmt.Sprintf("%v", f.value) }

// --------------------------------------------------------------------------------------------------------------------

type Function struct {
	parameters []*Identifier
	body       *BlockStatement
	env        *Environment
}

func (f *Function) Type() ObjectType { return FLOAT_OBJ }

func (f *Function) inspect() string {
	var buffer bytes.Buffer
	params := make([]string, 0)

	for _, param := range f.parameters {
		params = append(params, param.toString())
	}

	buffer.WriteString("fn(")
	buffer.WriteString(strings.Join(params, ", "))
	buffer.WriteString("){\n")
	buffer.WriteString("  " + f.body.toString())
	buffer.WriteString("\n}")

	return buffer.String()
}

// --------------------------------------------------------------------------------------------------------------------

type Integer struct {
	value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

func (i *Integer) inspect() string { return fmt.Sprintf("%v", i.value) }

// --------------------------------------------------------------------------------------------------------------------

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }

func (n *Null) inspect() string { return "null" }

// --------------------------------------------------------------------------------------------------------------------

type ReturnValue struct {
	value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_OBJ }

func (r *ReturnValue) inspect() string { return r.value.inspect() }

// --------------------------------------------------------------------------------------------------------------------

type StringValue struct {
	value string
}

func (s *StringValue) Type() ObjectType { return STRING_OBJ }

func (s *StringValue) inspect() string { return s.value }

// --------------------------------------------------------------------------------------------------------------------
// Environment
// --------------------------------------------------------------------------------------------------------------------

type Environment struct {
	store map[string]Object
	outer *Environment
}

// --------------------------------------------------------------------------------------------------------------------

func newEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: nil}
}

// --------------------------------------------------------------------------------------------------------------------

func newEnclosedEnvironment(outer *Environment) *Environment {
	env := newEnvironment()
	env.outer = outer

	return env
}

// --------------------------------------------------------------------------------------------------------------------

func (e *Environment) get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.get(name)
	}

	return obj, ok
}

// --------------------------------------------------------------------------------------------------------------------

func (e *Environment) set(name string, val Object) Object {
	e.store[name] = val

	return val
}

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
			return newError("push: wrong number of arguments. Got %v, want %v", len(args), 2)
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
}

// --------------------------------------------------------------------------------------------------------------------
