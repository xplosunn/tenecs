package codegen

import (
	"github.com/xplosunn/tenecs/typer/ast"
	"golang.org/x/exp/slices"
	"sort"
)

type _trackedDeclaration struct {
	Is        _isTrackedDeclaration
	VarName   ast.Ref
	TestSuite bool
}

type _isTrackedDeclaration string

const (
	isTrackedDeclarationGoMain            _isTrackedDeclaration = "go_main"
	isTrackedDeclarationWebWebApp         _isTrackedDeclaration = "web_webapp"
	isTrackedDeclarationUnitTest          _isTrackedDeclaration = "unit_test"
	isTrackedDeclarationGoIntegrationTest _isTrackedDeclaration = "go_integration_test"
)

type Runnables struct {
	GoMain    []ast.Ref
	WebWebApp []ast.Ref
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

	ast.SortRefs(runnables.GoMain)
	ast.SortRefs(runnables.WebWebApp)

	return runnables
}

type FoundTests struct {
	UnitTests          []ast.Ref
	UnitTestSuites     []ast.Ref
	GoIntegrationTests []ast.Ref
}

func FindTests(program *ast.Program) FoundTests {
	found := FoundTests{
		UnitTests:      []ast.Ref{},
		UnitTestSuites: []ast.Ref{},
	}

	programDeclarationNames := []ast.Ref{}
	for declarationName, _ := range program.Declarations {
		programDeclarationNames = append(programDeclarationNames, declarationName)
	}
	sort.Slice(programDeclarationNames, func(i, j int) bool {
		return programDeclarationNames[i].Package < programDeclarationNames[j].Package || programDeclarationNames[i].Name < programDeclarationNames[j].Name
	})

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
				} else if trackedDeclaration.Is == isTrackedDeclarationGoIntegrationTest {
					if trackedDeclaration.TestSuite {
						panic("go integration test shouldn't be a suite")
					} else {
						found.GoIntegrationTests = append(found.GoIntegrationTests, trackedDeclaration.VarName)
					}
				}
			}
		}
	}

	return found
}

type CachedTestCount struct {
	UnitTests          int
	UnitTestSuites     int
	GoIntegrationTests int
}

func RemoveCachedTests(tests FoundTests, cached FoundTests) (CachedTestCount, FoundTests) {
	result := FoundTests{
		UnitTests:          []ast.Ref{},
		UnitTestSuites:     []ast.Ref{},
		GoIntegrationTests: []ast.Ref{},
	}
	cachedCount := CachedTestCount{}
	for _, ref := range tests.UnitTests {
		if !slices.Contains(cached.UnitTests, ref) {
			result.UnitTests = append(result.UnitTests, ref)
		} else {
			cachedCount.UnitTests += 1
		}
	}
	for _, ref := range tests.UnitTestSuites {
		if !slices.Contains(cached.UnitTests, ref) {
			result.UnitTests = append(result.UnitTestSuites, ref)
		} else {
			cachedCount.UnitTestSuites += 1
		}
	}
	//TODO integration test caching
	result.GoIntegrationTests = tests.GoIntegrationTests
	return cachedCount, result
}

func checkTrackedDeclaration(declarationName ast.Ref, declarationExpression ast.Expression) *_trackedDeclaration {
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
		} else if caseKnownType.Name == "GoIntegrationTest" && caseKnownType.Package == "tenecs.test" {
			trackedDeclaration = &_trackedDeclaration{
				Is:      isTrackedDeclarationGoIntegrationTest,
				VarName: declarationName,
			}
		}
	}
	return trackedDeclaration
}
