package typer

import "github.com/xplosunn/tenecs/typer/types"

type Program struct {
	Modules []Module
}

type Module struct {
	Name       string
	Implements types.Interface
	Variables  map[string]types.VariableType
}
