package codegen

import (
	"github.com/xplosunn/tenecs/typer/ast"
	"sort"
)

type _trackedDeclaration struct {
	Is        _isTrackedDeclaration
	VarName   string
	TestSuite bool
}

type _isTrackedDeclaration string

const (
	isTrackedDeclarationGoMain    _isTrackedDeclaration = "go_main"
	isTrackedDeclarationWebWebApp _isTrackedDeclaration = "web_webapp"
	isTrackedDeclarationUnitTest  _isTrackedDeclaration = "unit_test"
)

type Runnables struct {
	GoMain    []string
	WebWebApp []string
}

func FindRunnables(program *ast.Program) Runnables {
	runnables := Runnables{}

	for declarationName, declarationExpression := range program.Declarations {
		trackedDeclaration := checkTrackedDeclaration(declarationName, declarationExpression)
		if trackedDeclaration != nil {
			if trackedDeclaration.Is == isTrackedDeclarationGoMain {
				runnables.GoMain = append(runnables.GoMain, trackedDeclaration.VarName)
			} else if trackedDeclaration.Is == isTrackedDeclarationWebWebApp {
				runnables.WebWebApp = append(runnables.WebWebApp, trackedDeclaration.VarName)
			}
		}
	}

	sort.Strings(runnables.GoMain)
	sort.Strings(runnables.WebWebApp)

	return runnables
}

type FoundTests struct {
	UnitTests      []string
	UnitTestSuites []string
}

func FindTests(program *ast.Program) FoundTests {
	found := FoundTests{
		UnitTests:      []string{},
		UnitTestSuites: []string{},
	}

	programDeclarationNames := []string{}
	for declarationName, _ := range program.Declarations {
		programDeclarationNames = append(programDeclarationNames, declarationName)
	}
	sort.Strings(programDeclarationNames)

	for _, declarationName := range programDeclarationNames {
		for decName, decExp := range program.Declarations {
			if decName != declarationName {
				continue
			}
			trackedDeclaration := checkTrackedDeclaration(decName, decExp)
			if trackedDeclaration != nil {
				if trackedDeclaration.Is == isTrackedDeclarationUnitTest {
					if trackedDeclaration.TestSuite {
						found.UnitTestSuites = append(found.UnitTestSuites, trackedDeclaration.VarName)
					} else {
						found.UnitTests = append(found.UnitTests, trackedDeclaration.VarName)
					}
				}
			}
		}
	}

	return found
}

func checkTrackedDeclaration(declarationName string, declarationExpression ast.Expression) *_trackedDeclaration {
	var trackedDeclaration *_trackedDeclaration = nil
	varType := ast.VariableTypeOfExpression(declarationExpression)
	_, _, caseKnownType, _, _ := varType.VariableTypeCases()
	if caseKnownType != nil {
		if caseKnownType.Name == "UnitTestSuite" && caseKnownType.Package == "tenecs.test" {
			trackedDeclaration = &_trackedDeclaration{
				Is:        isTrackedDeclarationUnitTest,
				VarName:   declarationName,
				TestSuite: true,
			}
		} else if caseKnownType.Name == "UnitTest" && caseKnownType.Package == "tenecs.test" {
			trackedDeclaration = &_trackedDeclaration{
				Is:      isTrackedDeclarationUnitTest,
				VarName: declarationName,
			}
		} else if caseKnownType.Name == "Main" && caseKnownType.Package == "tenecs.go" {
			trackedDeclaration = &_trackedDeclaration{
				Is:      isTrackedDeclarationGoMain,
				VarName: declarationName,
			}
		} else if caseKnownType.Name == "WebApp" && caseKnownType.Package == "tenecs.web" {
			trackedDeclaration = &_trackedDeclaration{
				Is:      isTrackedDeclarationWebWebApp,
				VarName: declarationName,
			}
		}
	}
	return trackedDeclaration
}
