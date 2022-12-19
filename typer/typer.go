package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"unicode"
)

func Validate(parsed parser.FileTopLevel) error {
	pkg, imports, modules := parser.FileTopLevelFields(parsed)
	err := validatePackage(pkg)
	if err != nil {
		return err
	}
	resolvedInterfaces, err := validateImports(imports, StdLib)
	if err != nil {
		return err
	}
	err = validateModules(modules, resolvedInterfaces)
	if err != nil {
		return err
	}

	return nil
}

func validatePackage(node parser.Package) *TypecheckError {
	identifier := parser.PackageFields(node)
	for _, r := range identifier {
		if !unicode.IsLower(r) {
			return &TypecheckError{Message: "package name should start with a lowercase letter"}
		} else {
			return nil
		}
	}
	return nil
}

func validateImports(nodes []parser.Import, knownPackages Package) (map[string]Interface, *TypecheckError) {
	resolvedInterfaces := map[string]Interface{}
	for _, node := range nodes {
		dotSeparatedNames := parser.ImportFields(node)
		if len(dotSeparatedNames) < 2 {
			return nil, &TypecheckError{Message: "all interfaces belong to a package"}
		}
		currPackage := knownPackages
		for i, name := range dotSeparatedNames {
			if i < len(dotSeparatedNames)-1 {
				p, ok := currPackage.Packages[name]
				if !ok {
					return nil, &TypecheckError{Message: "no package " + name + " found"}
				}
				currPackage = p
				continue
			}
			interf, ok := currPackage.Interfaces[name]
			if !ok {
				return nil, &TypecheckError{Message: "no interface " + name + " found"}
			}
			_, ok = resolvedInterfaces[name]
			if ok {
				return nil, &TypecheckError{Message: "already imported an interface with name " + name}
			}
			resolvedInterfaces[name] = interf
		}
	}
	return resolvedInterfaces, nil
}

func validateModules(nodes []parser.Module, resolvedInterfaces map[string]Interface) *TypecheckError {
	moduleNames := map[string]bool{}
	for _, node := range nodes {
		name, implements, declarations := parser.ModuleFields(node)
		_, ok := moduleNames[name]
		if ok {
			return &TypecheckError{Message: "another module declared with name " + name}
		}
		moduleNames[name] = true

		implementedInterfaces, err := validateImplementedInterfaces(implements, resolvedInterfaces)
		err = validateDeclarations(declarations, implementedInterfaces, resolvedInterfaces)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateImplementedInterfaces(implements []string, resolvedInterfaces map[string]Interface) ([]Interface, *TypecheckError) {
	implementedInterfaces := []Interface{}
	for _, implement := range implements {
		interf, ok := resolvedInterfaces[implement]
		if !ok {
			return implementedInterfaces, &TypecheckError{Message: "not found interface with name " + implement}
		}
		implementedInterfaces = append(implementedInterfaces, interf)
	}
	allInterfaceVariableNames := map[string]bool{}
	for _, implementedInterface := range implementedInterfaces {
		for varName, _ := range implementedInterface.Variables {
			_, ok := allInterfaceVariableNames[varName]
			if ok {
				return nil, &TypecheckError{Message: "imcompatible interfaces implemented because both shared a variable name"}
			}
			allInterfaceVariableNames[varName] = true
		}
	}
	return implementedInterfaces, nil
}

func validateDeclarations(nodes []parser.Declaration, implementedInterfaces []Interface, resolvedInterfaces map[string]Interface) *TypecheckError {
	for _, node := range nodes {
		public, name, lambda := parser.DeclarationFields(node)
		var typeOfInterfaceVariableWithSameName *VariableType
	typeOfInterfaceVariableWithSameNameLoop:
		for _, implementedInterface := range implementedInterfaces {
			for varName, varType := range implementedInterface.Variables {
				if varName == name {
					typeOfInterfaceVariableWithSameName = &varType
					break typeOfInterfaceVariableWithSameNameLoop
				}
			}
		}

		lambdaType, err := validateLambda(lambda, resolvedInterfaces)
		if err != nil {
			return err
		}

		if typeOfInterfaceVariableWithSameName == nil && public {
			return &TypecheckError{Message: "variable shouldn't be public: " + name}
		}
		if typeOfInterfaceVariableWithSameName != nil {
			if !public {
				return &TypecheckError{Message: "variable should be public: " + name}
			}
			if variableTypeEquals(lambdaType, *typeOfInterfaceVariableWithSameName) {
				return &TypecheckError{Message: "variable should be of the same type as the one on the implemented interface: " + name + ", expected " + printableName(*typeOfInterfaceVariableWithSameName) + ", got " + printableName(lambdaType)}
			}
		}
	}
	return nil
}

func variableTypeEquals(o1 VariableType, o2 VariableType) bool {
	return o1 == o2
}

func printableName(varType VariableType) string {
	caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
	if caseInterface != nil {
		result := "Interface with variables ("
		needsComma := false
		for key, _ := range caseInterface.Variables {
			if !needsComma {
				result = result + key
				needsComma = true
			} else {
				result = result + ", " + key
			}
		}
		return result + ")"
	} else if caseFunction != nil {
		result := "("
		for i, argumentType := range caseFunction.ArgumentTypes {
			if i == 0 {
				result = result + printableName(argumentType)
			} else {
				result = result + ", " + printableName(argumentType)
			}
		}
		return result + ") => " + printableName(caseFunction.ReturnType)
	} else if caseBasicType != nil {
		return caseBasicType.Type
	} else if caseVoid != nil {
		return "Void"
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}

func validateLambda(lambda parser.Lambda, resolvedInterfaces map[string]Interface) (*Function, *TypecheckError) {
	lambdaArgumentTypes := []VariableType{}
	var lambdaReturnType *VariableType
	parameters, block := parser.LambdaFields(lambda)
	for _, parameter := range parameters {
		interf, ok := resolvedInterfaces[parameter.Type]
		if !ok {
			return nil, &TypecheckError{Message: "not found type: " + parameter.Type}
		}
		lambdaArgumentTypes = append(lambdaArgumentTypes, interf)
	}
	for _, invocation := range block {
		dotSeparatedVarName, arguments := parser.InvocationFields(invocation)
		var currentContext Interface
		for i, varName := range dotSeparatedVarName {
			if i == 0 {
				foundLambdaParameterWithSameName := false
				for _, lambdaParameter := range parameters {
					if lambdaParameter.Name == varName {
						interf, ok := resolvedInterfaces[lambdaParameter.Type]
						if !ok {
							return nil, &TypecheckError{Message: "not found type: " + lambdaParameter.Type}
						}
						currentContext = interf
						foundLambdaParameterWithSameName = true
						break
					}
				}
				if !foundLambdaParameterWithSameName {
					return nil, &TypecheckError{Message: "not found a lambda parameter with name: " + varName}
				}
				continue
			}
			if i < len(dotSeparatedVarName)-1 {
				varType, ok := currentContext.Variables[varName]
				if !ok {
					return nil, &TypecheckError{Message: "not found variable: " + varName}
				}
				caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
				if caseInterface != nil {
					currentContext = *caseInterface
					continue
				} else if caseFunction != nil {
					return nil, &TypecheckError{Message: "expected interface but found function: " + varName}
				} else if caseBasicType != nil {
					return nil, &TypecheckError{Message: "expected interface but found basic type: " + varName}
				} else if caseVoid != nil {
					return nil, &TypecheckError{Message: "expected interface but found void: " + varName}
				} else {
					panic(fmt.Errorf("cases on %v", varType))
				}
				continue
			}

			varType, ok := currentContext.Variables[varName]
			if !ok {
				return nil, &TypecheckError{Message: "not found variable: " + varName}
			}
			caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			if caseInterface != nil {
				return nil, &TypecheckError{Message: "expected function but found interface: " + varName}
			} else if caseFunction != nil {
				argumentTypes, returnType := FunctionFields(*caseFunction)
				if len(arguments) != len(argumentTypes) {
					return nil, &TypecheckError{Message: fmt.Sprintf("Expected %d arguments but got %d", len(argumentTypes), len(arguments))}
				}
				for i2, argument := range arguments {
					expectedType := argumentTypes[i2]
					err := isOfExpectedType(argument, expectedType)
					if err != nil {
						return nil, err
					}
				}
				lambdaReturnType = &returnType
			} else if caseBasicType != nil {
				return nil, &TypecheckError{Message: "expected function but found basic type: " + varName}
			} else if caseVoid != nil {
				return nil, &TypecheckError{Message: "expected function but found void: " + varName}
			} else {
				panic(fmt.Errorf("cases on %v", varType))
			}
		}
	}

	if lambdaReturnType == nil {
		return nil, &TypecheckError{Message: "could not resolve lambda return type"}
	}

	lambdaType := Function{
		ArgumentTypes: lambdaArgumentTypes,
		ReturnType:    *lambdaReturnType,
	}
	return &lambdaType, nil
}

func isOfExpectedType(argument parser.Literal, expectedType VariableType) *TypecheckError {
	caseInterface, caseFunction, caseBasicType, caseVoid := expectedType.Cases()
	if caseInterface != nil {
		panic("TODO")
	} else if caseFunction != nil {
		panic("TODO")
	} else if caseBasicType != nil {
		basicType := *caseBasicType
		expectBasicType := func(typeName string) *TypecheckError {
			if basicType.Type != typeName {
				return &TypecheckError{Message: "expected type " + typeName + " but found " + basicType.Type}
			}
			return nil
		}
		return parser.LiteralFold[*TypecheckError](
			argument,
			func(arg float64) *TypecheckError {
				return expectBasicType("Float")
			},
			func(arg int) *TypecheckError {
				return expectBasicType("Int")
			},
			func(arg string) *TypecheckError {
				return expectBasicType("String")
			},
			func(arg bool) *TypecheckError {
				return expectBasicType("Boolean")
			},
		)
	} else if caseVoid != nil {
		panic("TODO")
	} else {
		panic(fmt.Errorf("cases on %v", expectedType))
	}
}

type TypecheckError struct {
	Message string
}

func (t TypecheckError) Error() string {
	return t.Message
}

type Package struct {
	Packages   map[string]Package
	Interfaces map[string]Interface
}

type VariableType interface {
	sealedVariableType()
	Cases() (*Interface, *Function, *BasicType, *Void)
}

type Interface struct {
	Variables map[string]VariableType
}

func (i Interface) sealedVariableType() {}
func (i Interface) Cases() (*Interface, *Function, *BasicType, *Void) {
	return &i, nil, nil, nil
}

type Function struct {
	ArgumentTypes []VariableType
	ReturnType    VariableType
}

func (f Function) sealedVariableType() {}
func (f Function) Cases() (*Interface, *Function, *BasicType, *Void) {
	return nil, &f, nil, nil
}

func FunctionFields(function Function) ([]VariableType, VariableType) {
	return function.ArgumentTypes, function.ReturnType
}

type BasicType struct {
	Type string
}

func (b BasicType) sealedVariableType() {}
func (b BasicType) Cases() (*Interface, *Function, *BasicType, *Void) {
	return nil, nil, &b, nil
}

type Void struct {
}

func (v Void) sealedVariableType() {}
func (v Void) Cases() (*Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, &v
}
