package main

func loadMap(env *Environment) {
	input := `
	let map = 
		fn(arr, f) {
			let iter = fn(arr, acc) { 
				if (len(arr) == 0) { 
					acc 
				} else { 
				 	iter(rest(arr), push(acc, f(first(arr))))}
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
	let reduce = fn(arr, initial, f) {
		let iter = fn(arr, result) {
			if (len(arr) == 0) {
			result
		} else {
			iter(rest(arr), f(result, first(arr)));
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
