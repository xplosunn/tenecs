generate:
	go generate ./...

updateSublimeSyntax:
	go run syntaxhighlight/main.go && cp tenecs.sublime-syntax ~/Library/Application\ Support/Sublime\ Text/Packages/User/ && rm tenecs.sublime-syntax