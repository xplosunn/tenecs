package codegen_js

import (
	"errors"
	"fmt"
	"strings"
)

func NodeProgramToPrintWebAppExternalGenerate(pkgName string, programJs string, webAppVarName string) string {
	return programJs + fmt.Sprintf(`

for (const file of %s.external) {
  if (!file["$type"]) {
    console.log("ERROR: internal error: lack of $type");
    break;
  }
  if (file["$type"] == "CssUrl") {
    console.log("CSS:" + file.url);
  } else {
    console.log("ERROR: unexpected $type: " + file["$type"]);
    break;
  }
}

`, variableName(&pkgName, webAppVarName))
}

func NodeProgramToPrintWebAppExternalReadOutput(output string) ([]string, error) {
	cssFiles := []string{}
	for _, line := range strings.Split(output, "\n") {
		if len(strings.TrimSpace(line)) == 0 {
			// empty line
		} else if strings.HasPrefix(line, "CSS:") {
			cssFiles = append(cssFiles, strings.TrimPrefix(line, "CSS:"))
		} else if strings.HasPrefix(line, "ERROR:") {
			return nil, errors.New("Error reading external: " + strings.TrimPrefix(line, "ERROR:"))
		} else {
			return nil, errors.New("Error reading external, unexpected line: " + line)
		}
	}
	return cssFiles, nil
}
