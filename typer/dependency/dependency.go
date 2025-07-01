package dependency

import (
	"github.com/xplosunn/tenecs/desugar"
	"golang.org/x/exp/slices"
)

func DependenciesOfSinglePackage(parsedPackage map[string]desugar.FileTopLevel) []string {
	if len(parsedPackage) == 0 {
		panic("dependencies of empty package")
	}
	pkg := ""
	for _, fileTopLevel := range parsedPackage {
		pkgOfThisFile := ""
		for i, name := range fileTopLevel.Package.DotSeparatedNames {
			if i > 0 {
				pkgOfThisFile += "."
			}
			pkgOfThisFile += name.String
		}
		if pkg == "" {
			pkg = pkgOfThisFile
		} else if pkg != pkgOfThisFile {
			panic("multiple packages")
		}
	}
	nonStandardLibrabryDependencies := []string{}
	for _, fileTopLevel := range parsedPackage {
		for _, oneImport := range fileTopLevel.Imports {
			if oneImport.DotSeparatedVars[0].String == "tenecs" {
				continue
			}
			importedPkg := ""
			for i, name := range oneImport.DotSeparatedVars {
				if i < len(oneImport.DotSeparatedVars)-1 {
					if i > 0 {
						importedPkg += "."
					}
					importedPkg += name.String
				}
			}
			if !slices.Contains(nonStandardLibrabryDependencies, importedPkg) {
				nonStandardLibrabryDependencies = append(nonStandardLibrabryDependencies, importedPkg)
			}
		}
	}
	return nonStandardLibrabryDependencies
}
