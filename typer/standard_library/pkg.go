package standard_library

import "github.com/xplosunn/tenecs/typer/types"

type Package struct {
	Packages   map[string]Package
	Interfaces map[string]*types.Interface
	Variables  map[string]types.VariableType
}

func packageWith(opts ...func(*Package)) Package {
	pkg := &Package{
		Packages:   map[string]Package{},
		Interfaces: map[string]*types.Interface{},
		Variables:  map[string]types.VariableType{},
	}
	for _, opt := range opts {
		opt(pkg)
	}
	return *pkg
}

func withPackage(name string, pack Package) func(pkg *Package) {
	return func(pkg *Package) {
		pkg.Packages[name] = pack
	}
}

func withInterface(name string, interf *types.Interface) func(pkg *Package) {
	return func(pkg *Package) {
		pkg.Interfaces[name] = interf
	}
}

func withFunction(name string, function *types.Function) func(pkg *Package) {
	return func(pkg *Package) {
		pkg.Variables[name] = function
	}
}
