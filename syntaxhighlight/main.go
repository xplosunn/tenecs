package main

import (
	"os"
	"path/filepath"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(wd, "tenecs.sublime-syntax")

	f, err := os.Create(filePath)
	defer f.Close()

	c, err := content()
	if err != nil {
		panic(err)
	}
	contentBytes := []byte(c)
	err = os.WriteFile(filePath, contentBytes, 0644)
	if err != nil {
		panic(err)
	}

}

func content() (string, error) {
	c := `%YAML 1.2
---
# See http://www.sublimetext.com/docs/syntax.html
file_extensions:
  - 10x
scope: source.tenecs
variables:
  ident: \b[[:alpha:]_][[:alnum:]_]*\b

contexts:
  main:
    # Strings begin and end with quotes, and use backslashes as an escape
    # character
    - match: '"'
      scope: punctuation.definition.string.begin.tenecs
      push: double_quoted_string

    # Note that blackslashes don't need to be escaped within single quoted
    # strings in YAML. When using single quoted strings, only single quotes
    # need to be escaped: this is done by using two single quotes next to each
    # other.
    - match: '\b(import|struct|interface|if|else)\b'
      scope: keyword.control.tenecs

    - match: \bpackage\b
      scope: keyword.declaration.namespace.tenecs
      push: pop-package-name

    # Numbers
    - match: '\b(-)?[0-9.]+\b'
      scope: constant.numeric.tenecs

    - match: ':='
      scope: keyword.operator.assignment.tenecs
    
    - match: ':'
      push: pop-type-name

    - match: 'implement'
      scope: keyword.control.tenecs
      push: pop-type-name

    - match: '<'
      push: type-name-list-generics

    - match: '\['
      push: pop-type-name

    - match: '\b(is)\b'
      scope: keyword.control.tenecs
      push: pop-type-name

    - match: '\|'
      scope: keyword.operator.bitwise.tenecs
      push: pop-type-name

  double_quoted_string:
    - meta_scope: string.quoted.double.tenecs
    - match: '\\.'
      scope: constant.character.escape.tenecs
    - match: '"'
      scope: punctuation.definition.string.end.tenecs
      pop: true

  pop-package-name:
    - match: '{{ident}}'
      scope: entity.name.namespace.tenecs
      pop: true

  pop-type-name:
    - match: '{{ident}}'
      scope: variable.function.tenecs
      pop: true

  type-name-list-generics:
    - match: '\,'
    - match: '{{ident}}'
      scope: variable.function.tenecs
    - match: '>'
      pop: true
`

	return c, nil
}
