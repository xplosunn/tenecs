package main

import (
	"fmt"
	"github.com/xplosunn/tenecs/codegen/javascript"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
)

func main() {
	res, err := parser.ParseString(`
package main

import tenecs.os.Runtime
import tenecs.os.Main

implementing Main module app() {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}

`)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", res)

	program, err := typer.Typecheck(*res)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nProgram:\n%v\n", *program)

	fmt.Printf("\nJS:\n\n")

	js, err := javascript.Codegen(*program)
	if err != nil {
		panic(err)
	}

	fmt.Printf(js)

}
