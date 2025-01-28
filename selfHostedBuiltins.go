package main

func loadNativeBuiltins(env *Environment) {
	loadMap(env)
	loadRecude(env)
	loadFilter(env)
}

func loadMap(env *Environment) {
	input := `
	let map = 
		fn(arr, function) {
			let iter = fn(arr, acc) { 
				if (len(arr) == 0) { 
					acc 
				} else { 
				 	iter(rest(arr), push(acc, function(first(arr))))
				}
			} 
			iter(arr, [])
		}
`
	lexer := newLexer(input)
	parser := newParser(lexer)
	program := parser.parseProgram()
	eval(program, env)
}

func loadRecude(env *Environment) {
	input := `
	let reduce = fn(arr, initial, function) {
		let iter = fn(arr, result) {
			if (len(arr) == 0) {
			result
		} else {
			iter(rest(arr), function(result, first(arr)));
		}
	}
		iter(arr, initial)
	}
`
	lexer := newLexer(input)
	parser := newParser(lexer)
	program := parser.parseProgram()
	eval(program, env)
}

func loadFilter(env *Environment) {
	input := `
	let filter = fn(arr, function) {
		let iter = fn(arr, acc) {
			if (len(arr) == 0) {
				acc
			} else {
			 	let working = first(arr)
				if (function(working)) {
					iter(rest(arr), push(acc, working))
				} else {
					iter(rest(arr), acc) 
				}			
			}
		}
		iter(arr, [])
	}
`
	lexer := newLexer(input)
	parser := newParser(lexer)
	program := parser.parseProgram()
	eval(program, env)
}
