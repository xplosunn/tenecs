package codegen_golang

const (
	colorReset = "\u001B[0m"

	colorGreen = "\u001B[32m"
	colorRed   = "\u001B[31m"
)

func Green(text string) string {
	return colorGreen + text + colorReset
}

func Red(text string) string {
	return colorRed + text + colorReset
}
