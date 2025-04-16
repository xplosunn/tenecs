package parser_typer_test

import (
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestExpectedGenericFunctionInvoked4(t *testing.T) {
	program := validProgram(t, testcode.GenericFunctionInvoked4)
	snaps.MatchStandaloneSnapshot(t, program)
}

func TestExpectedGenericFunctionDoubleInvoked(t *testing.T) {
	program := validProgram(t, testcode.GenericFunctionDoubleInvoked)
	snaps.MatchStandaloneSnapshot(t, program)
}

func TestWrongGeneric(t *testing.T) {
	invalidProgram(t, `
package mypackage

struct Tuple<L, R>(left: L, right: R)

leftAs := <L, R, T>(tuple: Tuple<L, R>, as: T): Tuple<T, R> => {
  result := Tuple<T, T>(as, as)
  result
}
`, "expected type mypackage.Tuple<T, R> but found mypackage.Tuple<T, T>")
}

func TestWrongGeneric2(t *testing.T) {
	invalidProgram(t, `
package mypackage

struct Tuple<L, R>(left: L, right: R)

leftAs := <L, R, T>(tuple: Tuple<L, R>, as: T): Tuple<T, R> => {
  Tuple<T, T>(as, as)
}
`, "expected type mypackage.Tuple<T, R> but found mypackage.Tuple<T, T>")
}

func TestGenericFunctionInvocation(t *testing.T) {
	validProgram(t, `
package mypackage

takeList := <A>(arr: List<A>): Void => {
  null
}

usage := (): Void => {
  takeList<String | Int>(<Int | String>["", 1])
  null
}
`)
}

func TestGenericFunctionInvocation2(t *testing.T) {
	validProgram(t, `
package mypackage

take := <A>(a: A): Void => {
  null
}

usage := (): Void => {
  take<List<String> | String>(<String>[])
  null
}
`)
}

func TestGenericFunctionInvocation3(t *testing.T) {
	validProgram(t, `
package mypackage

struct Parser<T>()

parseList := <Of>(parserOf: Parser<Of>): Parser<List<Of>> => {
  Parser<List<Of>>()
}

parseString := (): Parser<String> => {
  Parser<String>()
}

takeParser := <Of>(parser: Parser<Of>): Void => {
  null
}

usage := (): Void => {
  takeParser<List<List<String>>>(parseList<List<String>>(parseList<String>(parseString())))
}

`)
}

func TestGenericFunctionInvocation4(t *testing.T) {
	validProgram(t, `
package mypackage

wrapFunction := <R>(f: () ~> R): () ~> R => {
  (): R => {
    f()
  }
}

usage := (): Void => {
  f := wrapFunction<Void>(() => null)
  f()
}
`)
}

func TestGenericFunctionInvocation5(t *testing.T) {
	validProgram(t, `
package mypackage

apply := <A, B>(a: A, f: (A) ~> B): B => {
  f(a)
}

usage := (): String => {
  apply(1, (int: Int): String => {""})
}
`)
}

func TestGenericFunctionWrongInvocation(t *testing.T) {
	invalidProgram(t, `
package mypackage

take := <A>(arg: A): Void => {
  null
}

usage := (): Void => {
  take<String>(1)
  null
}

`, "expected type String but found Int")
}

func TestGenericFunctionWrongInvocation2(t *testing.T) {
	invalidProgram(t, `
package mypackage

takeList := <A>(arr: List<A>): Void => {
  null
}

usage := (): Void => {
  takeList<String>(<Int>[1])
  null
}

`, "expected List<String> but got List<Int>")
}

func TestGenericFunctionWrongInvocation3(t *testing.T) {
	invalidProgram(t, `
package mypackage

takeList := <A>(arr: List<A>): Void => {
  null
}

usage := (): Void => {
  takeList<String>(<String | Int>[""])
  null
}

`, "expected List<String> but got List<String | Int>")
}

func TestGenericFunctionWrongInvocation4(t *testing.T) {
	invalidProgram(t, `
package mypackage

takeList := <A>(arr: List<A>): Void => {
  null
}

usage := (): Void => {
  takeList<List<String>>(<String>[])
  null
}

`, "expected List<List<String>> but got List<String>")
}

func TestGenericFunctionWrongInvocation5(t *testing.T) {
	invalidProgram(t, `
package mypackage

assertEqual := <T> (a: T, b: T): Void => {
  null
}

listOfStringOrString := (): List<String> | String => {
  ""
}

usage := (): Void => {
	assertEqual<List<String>>(<String>[], listOfStringOrString())
}

`, "expected type List<String> but found List<String> | String")
}
