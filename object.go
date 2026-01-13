package main

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

// --------------------------------------------------------------------------------------------------------------------
// Allowed objects
// --------------------------------------------------------------------------------------------------------------------

type ObjectType string

type Hashable interface {
	HashKey() HashKey
}

type BuiltinFunc func(args ...Object) Object

const (
	ARRAY_OBJ    = "ARRAY"
	BOOL_OBJ     = "BOOLEAN"
	BUILTIN_OBJ  = "BUILTIN"
	ERR_OBJ      = "ERROR_OBJ"
	FLOAT_OBJ    = "FLOAT"
	FUNCTION_OBJ = "FUNCTION"
	HASH_OBJ     = "HASH_OBJ"
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

	for idx, elem := range a.elements {
		elements = append(elements, elem.inspect())
		if idx == 5 {
			buffer.WriteString("[")
			buffer.WriteString(strings.Join(elements, ", "))
			buffer.WriteString(", ...")
			buffer.WriteString("]")

			return buffer.String()
		}
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

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{keyType: b.Type(), value: value}
}

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

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

func (f *Function) inspect() string {
	params := make([]string, 0)

	for _, param := range f.parameters {
		params = append(params, param.toString())
	}

	return fmt.Sprintf("fn(%v)", strings.Join(params, ", "))
}

// --------------------------------------------------------------------------------------------------------------------

type Hash struct {
	pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

func (h *Hash) inspect() string {
	var buffer bytes.Buffer
	pairs := make([]string, 0)

	for _, pair := range h.pairs {
		pairs = append(pairs, fmt.Sprintf("%v: %v", pair.key.inspect(), pair.value.inspect()))
	}

	buffer.WriteString("{")
	buffer.WriteString(strings.Join(pairs, ", "))
	buffer.WriteString("}")

	return buffer.String()
}

// --------------------------------------------------------------------------------------------------------------------

type HashKey struct {
	keyType ObjectType
	value   uint64
}

// --------------------------------------------------------------------------------------------------------------------

type HashPair struct {
	key   Object
	value Object
}

// --------------------------------------------------------------------------------------------------------------------

type Integer struct {
	value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

func (i *Integer) inspect() string { return fmt.Sprintf("%v", i.value) }

func (i *Integer) HashKey() HashKey {
	return HashKey{keyType: i.Type(), value: uint64(i.value)}
}

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

func (s *StringValue) HashKey() HashKey {
	hash := fnv.New64a()
	hash.Write([]byte(s.value))

	return HashKey{keyType: s.Type(), value: hash.Sum64()}
}

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
