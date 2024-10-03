package formatter

import "github.com/xplosunn/tenecs/parser"

func DisplayFileTopLevel(parsed parser.FileTopLevel) string {
	return displayFileTopLevel(parsed, false)
}

func DisplayFileTopLevelIgnoringComments(parsed parser.FileTopLevel) string {
	return displayFileTopLevel(parsed, true)
}

func DisplayExpression(expression parser.Expression) string {
	return displayExpression(expression)
}

func DisplayDeclaration(declaration parser.Declaration) string {
	return displayDeclaration(declaration)
}
