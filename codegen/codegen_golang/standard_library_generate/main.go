package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xplosunn/tenecs/typer/standard_library"
	goast "go/ast"
	goparser "go/parser"
	goprinter "go/printer"
	gotoken "go/token"
	"log"
	"os"
	"sort"
)

func main() {
	fmt.Println("Starting codegen golang standard_library")
	functionNames := handlePackage("", standard_library.StdLib)
	generateInit(functionNames)
}

func generateInit(functionNames []string) {
	filePath := fmt.Sprintf("../standard_library/%s.go", "init")

	functions := ""
	sort.Strings(functionNames)
	for _, functionName := range functionNames {
		functions += fmt.Sprintf(`"%s": %s(),`, functionName, functionName) + "\n"
	}

	fileContent := fmt.Sprintf(`package standard_library

// ###############################################
// # This file is generated via code-generation. #
// # Check gen.go                                #
// ###############################################

var Functions = map[string]Function{
%s}
`, functions)

	if fileExists(filePath) {
		err := os.Remove(filePath)
		if err != nil {
			fail(err)
		}
	}
	err := os.WriteFile(filePath, []byte(fileContent), os.ModePerm)
	if err != nil {
		fail(err)
	}
}

func handlePackage(namespace string, pkg standard_library.Package) []string {
	functionNames := []string{}
	for pkgName, innerPkg := range pkg.Packages {
		pkgNameSpace := namespace
		if pkgNameSpace != "" {
			pkgNameSpace += "_"
		}
		pkgNameSpace += pkgName
		functionNames = append(functionNames, handlePackage(pkgNameSpace, innerPkg)...)
	}

	if len(pkg.Variables) == 0 && len(pkg.Structs) == 0 {
		return functionNames
	}

	filePath := fmt.Sprintf("../standard_library/%s_src.go", namespace)
	if !fileExists(filePath) {
		os.WriteFile(filePath, []byte(baseFile()), os.ModePerm)
	}
	src, err := os.ReadFile(filePath)
	if err != nil {
		fail(err)
	}
	fset := gotoken.NewFileSet()
	parsedFile, err := goparser.ParseFile(fset, "", src, goparser.SkipObjectResolution)
	if err != nil {
		fail(err)
	}
	file := parsedFile

	for varName, variableType := range pkg.Variables {
		caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
		if caseTypeArgument != nil {
			failWithMessage("handlePackage caseTypeArgument")
		} else if caseList != nil {
			failWithMessage("handlePackage caseList")
		} else if caseKnownType != nil {
			failWithMessage("handlePackage caseKnownType")
		} else if caseFunction != nil {
			functionNames = append(functionNames, handleFunction(file, namespace, varName))
		} else if caseOr != nil {
			failWithMessage("handlePackage caseOr")
		} else {
			fail(fmt.Errorf("cases on %v", variableType))
		}
	}

	for varName, _ := range pkg.Structs {
		functionNames = append(functionNames, handleFunction(file, namespace, varName))
	}

	buf := new(bytes.Buffer)
	err = goprinter.Fprint(buf, fset, file)
	if err != nil {
		fail(err)
	}
	newFileContent := buf.String()
	err = os.Remove(filePath)
	if err != nil {
		fail(err)
	}
	err = os.WriteFile(filePath, []byte(newFileContent), os.ModePerm)
	if err != nil {
		fail(err)
	}

	return functionNames
}

func handleFunction(file *goast.File, namespace string, name string) string {
	functionName := namespace + "_" + name
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*goast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Name.Name == functionName {
			return functionName
		}
	}
	file.Decls = append(file.Decls, &goast.FuncDecl{
		Doc:  nil,
		Recv: nil,
		Name: &goast.Ident{
			Name: functionName,
		},
		Type: &goast.FuncType{
			TypeParams: nil,
			Params: &goast.FieldList{
				List: nil,
			},
			Results: &goast.FieldList{
				List: []*goast.Field{
					&goast.Field{
						Type: &goast.Ident{
							Name: "Function",
						},
					},
				},
			},
		},
		Body: &goast.BlockStmt{
			Lbrace: 0,
			List:   nil,
			Rbrace: 0,
		},
	})

	return functionName
}

func failWithMessage(msg string) {
	fail(errors.New(msg))
}

func fail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func baseFile() string {
	return `package standard_library

// ##################################################################
// # The signatures of this file are generated via code-generation. #
// # Check gen.go                                                   #
// ##################################################################
`
}
