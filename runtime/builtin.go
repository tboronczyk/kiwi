package runtime

import (
	"fmt"
)


// built-in functions, [name]func{implementation}
var builtins = map[string]func(Params){
	// write - prints a value
	"write": func(p Params) {
		fmt.Print(p[0].Value)
	},
}
