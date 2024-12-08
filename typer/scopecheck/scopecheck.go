package scopecheck

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/types"
)

type ScopeCheckError interface {
	Error() string
	Node() parser.Node
	scopeCheckErrorImpl() scopeCheckError
}

type scopeCheckError struct {
	node    parser.Node
	message string
}

func (e scopeCheckError) Error() string {
	return e.message
}

func (e scopeCheckError) Node() parser.Node {
	return e.node
}

func (e scopeCheckError) scopeCheckErrorImpl() scopeCheckError {
	return e
}

func ptrScopeCheckError(node parser.Node, message string) ScopeCheckError {
	return scopeCheckError{
		node:    node,
		message: message,
	}
}

func ValidateTypeAnnotationInScope(typeAnnotation parser.TypeAnnotation, file string, scope binding.Scope) (types.VariableType, ScopeCheckError) {
	switch len(typeAnnotation.OrTypes) {
	case 0:
		return nil, ptrScopeCheckError(typeAnnotation.Node, "unexpected error ValidateTypeAnnotationInScope no types found")
	case 1:
		elem := typeAnnotation.OrTypes[0]
		return ValidateTypeAnnotationElementInScope(elem, file, scope)
	default:
		elements := []types.VariableType{}
		for _, element := range typeAnnotation.OrTypes {
			newElement, err := ValidateTypeAnnotationElementInScope(element, file, scope)
			if err != nil {
				return nil, err
			}
			elements = append(elements, newElement)
		}
		return &types.OrVariableType{
			Elements: elements,
		}, nil
	}
}

func ValidateTypeAnnotationElementInScope(typeAnnotationElement parser.TypeAnnotationElement, file string, scope binding.Scope) (types.VariableType, ScopeCheckError) {
	var varType types.VariableType
	var err ScopeCheckError
	parser.TypeAnnotationElementExhaustiveSwitch(
		typeAnnotationElement,
		func(underscoreTypeAnnotation parser.SingleNameType) {
			err = ptrScopeCheckError(underscoreTypeAnnotation.Node, "Generic inference not allowed here")
		},
		func(typeAnnotation parser.SingleNameType) {
			genericTypes := []types.VariableType{}
			for _, generic := range typeAnnotation.Generics {
				genericVarType, err2 := ValidateTypeAnnotationInScope(generic, file, scope)
				if err2 != nil {
					err = err2
					return
				}
				genericTypes = append(genericTypes, genericVarType)
			}
			varType2, err2 := binding.GetTypeByTypeName(scope, file, typeAnnotation.TypeName.String, genericTypes)
			varType = varType2
			if err2 != nil {
				err = scopeCheckError{
					node:    typeAnnotation.TypeName.Node,
					message: err2.Problem,
				}
			}
		},
		func(typeAnnotation parser.FunctionType) {
			localScope := scope
			var bindingErr *binding.ResolutionError
			for _, generic := range typeAnnotation.Generics {
				localScope, bindingErr = binding.CopyAddingTypeToFile(localScope, file, generic, &types.TypeArgument{Name: generic.String})
				if bindingErr != nil {
					err = scopeCheckError{
						node:    generic.Node,
						message: bindingErr.Problem,
					}
					return
				}
			}
			arguments := []types.FunctionArgument{}
			for _, argAnnotatedType := range typeAnnotation.Arguments {
				varType, err = ValidateTypeAnnotationInScope(argAnnotatedType.Type, file, localScope)
				if err != nil {
					return
				}
				name := "_"
				if argAnnotatedType.Name != nil {
					name = argAnnotatedType.Name.String
				}
				arguments = append(arguments, types.FunctionArgument{
					Name:         name,
					VariableType: varType,
				})
			}
			var returnType types.VariableType
			returnType, err = ValidateTypeAnnotationInScope(typeAnnotation.ReturnType, file, localScope)
			if err != nil {
				return
			}
			generics := []string{}
			for _, generic := range typeAnnotation.Generics {
				generics = append(generics, generic.String)
			}
			if typeAnnotation.Generics == nil {
				generics = nil
			}
			varType = &types.Function{
				Generics:   generics,
				Arguments:  arguments,
				ReturnType: returnType,
			}
		},
	)
	return varType, err
}
