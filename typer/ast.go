package typer

type Program struct {
	Modules []Module
}

type Module struct {
	Name       string
	Implements []Interface
	Variables  map[string]VariableType
}
