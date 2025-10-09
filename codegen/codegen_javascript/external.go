package codegen_javascript

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xplosunn/tenecs/ir"
)

func NodeProgramToPrintWebAppExternalGenerate(programJs string, webAppVarName ir.Reference) string {
	return programJs + fmt.Sprintf(`

for (const file of %s()._external.value) {
  if (!file["$type"]) {
    console.log("ERROR: internal error: lack of $type");
    break;
  }
  if (file["$type"] == "tenecs.web.CssUrl") {
    console.log("CSS:" + file._url.value);
  } else {
    console.log("ERROR: unexpected $type: " + file["$type"]);
    break;
  }
}

`, webAppVarName.Name)
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
