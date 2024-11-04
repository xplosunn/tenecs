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

func FindMains(program *ast.Program) []string {
	trackedDeclarationMains := []string{}

	programDeclarationNames := []string{}
	for _, declaration := range program.Declarations {
		programDeclarationNames = append(programDeclarationNames, declaration.Name)
	}
	sort.Strings(programDeclarationNames)

	for _, declarationName := range programDeclarationNames {
		for _, declaration := range program.Declarations {
			if declaration.Name != declarationName {
				continue
			}
			trackedDeclaration := checkTrackedDeclaration(declaration)
			if trackedDeclaration != nil {
				if trackedDeclaration.Is == isTrackedDeclarationGoMain {
					trackedDeclarationMains = append(trackedDeclarationMains, trackedDeclaration.VarName)
				}
			}
		}
	}

	return trackedDeclarationMains
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
	for _, declaration := range program.Declarations {
		programDeclarationNames = append(programDeclarationNames, declaration.Name)
	}
	sort.Strings(programDeclarationNames)

	for _, declarationName := range programDeclarationNames {
		for _, declaration := range program.Declarations {
			if declaration.Name != declarationName {
				continue
			}
			trackedDeclaration := checkTrackedDeclaration(declaration)
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

func checkTrackedDeclaration(declaration *ast.Declaration) *_trackedDeclaration {
	var trackedDeclaration *_trackedDeclaration = nil
	varType := ast.VariableTypeOfExpression(declaration.Expression)
	_, caseKnownType, _, _ := varType.VariableTypeCases()
	if caseKnownType != nil {
		if caseKnownType.Name == "UnitTestSuite" && caseKnownType.Package == "tenecs.test" {
			trackedDeclaration = &_trackedDeclaration{
				Is:        isTrackedDeclarationUnitTest,
				VarName:   declaration.Name,
				TestSuite: true,
			}
		} else if caseKnownType.Name == "UnitTest" && caseKnownType.Package == "tenecs.test" {
			trackedDeclaration = &_trackedDeclaration{
				Is:      isTrackedDeclarationUnitTest,
				VarName: declaration.Name,
			}
		} else if caseKnownType.Name == "Main" && caseKnownType.Package == "tenecs.go" {
			trackedDeclaration = &_trackedDeclaration{
				Is:      isTrackedDeclarationGoMain,
				VarName: declaration.Name,
			}
		} else if caseKnownType.Name == "WebApp" && caseKnownType.Package == "tenecs.web" {
			trackedDeclaration = &_trackedDeclaration{
				Is:      isTrackedDeclarationWebWebApp,
				VarName: declaration.Name,
			}
		}
	}
	return trackedDeclaration
}
