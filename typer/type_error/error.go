package type_error

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"strings"
)

type TypecheckError struct {
	Node    parser.Node
	Message string
}

func (t TypecheckError) Error() string {
	return t.Message
}

func PtrOnNodef(node parser.Node, format string, a ...any) *TypecheckError {
	return &TypecheckError{
		Node:    node,
		Message: fmt.Sprintf(format, a...),
	}
}

func Render(program string, err *TypecheckError) (string, error) {
	errLineIndex := err.Node.Pos.Line - 1

	programLines := strings.Split(program, "\n")

	prevLines := safeSlice(programLines, errLineIndex-3, errLineIndex)
	nextLines := safeSlice(programLines, errLineIndex+1, errLineIndex+3)

	pad := digitsLen(err.Node.Pos.Line + len(nextLines))
	prevFrom := err.Node.Pos.Line - len(prevLines)
	errFrom := err.Node.Pos.Line
	nextFrom := err.Node.Pos.Line + 1

	errLine := programLines[errLineIndex]

	errReportLine := strings.Repeat(" ", err.Node.Pos.Column+pad+4)
	errReportLine += "^ " + err.Error()

	result := strings.Join(prefixLinesWithLineNumber(prevLines, pad, prevFrom), "\n") + "\n"
	result += prefixLineWithLineNumber(errLine, pad, errFrom) + "\n"
	result += errReportLine + "\n"
	result += strings.Join(prefixLinesWithLineNumber(nextLines, pad, nextFrom), "\n")

	return result, nil
}

func prefixLinesWithLineNumber(lines []string, pad int, from int) []string {
	result := []string{}
	for i, line := range lines {
		result = append(result, prefixLineWithLineNumber(line, pad, from+i))
	}
	return result
}

func prefixLineWithLineNumber(line string, pad int, from int) string {
	return fmt.Sprintf("| %-"+fmt.Sprintf("%d", pad)+"d | %s", from, line)
}

func digitsLen(i int) int {
	if i >= 1e18 {
		return 19
	}
	x, count := 10, 1
	for x <= i {
		x *= 10
		count++
	}
	return count
}

func safeSlice[T any](arr []T, from int, to int) []T {
	if to < from {
		panic(fmt.Sprintf("safeSlice requires to >= from, but %d is smaller than %d", to, from))
	}
	if from >= len(arr) || to <= 0 {
		return []T{}
	}
	if from < 0 {
		from = 0
	}
	if to >= len(arr) {
		to = len(arr) - 1
	}
	return arr[from:to]
}
