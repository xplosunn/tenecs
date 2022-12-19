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

module app: Main {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}

`)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", res)

	err = typer.Validate(*res)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nJS:\n\n")

	js, err := javascript.Codegen(*res)
	if err != nil {
		panic(err)
	}

	fmt.Printf(js)

}
