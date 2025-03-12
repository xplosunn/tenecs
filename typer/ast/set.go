package ast

import "golang.org/x/exp/maps"

type Set[A comparable] map[A]bool

func (s Set[A]) Contains(a A) bool {
	return s[a]
}

func (s Set[A]) Put(a A) {
	s[a] = true
}

func (s Set[A]) PutAll(as []A) {
	for _, a := range as {
		s[a] = true
	}
}

func (s Set[A]) Remove(a A) {
	delete(s, a)
}

func (s Set[A]) Elements() []A {
	return maps.Keys(s)
}

func (s Set[A]) Copy() Set[A] {
	result := Set[A]{}
	for k, _ := range s {
		result[k] = true
	}
	return result
}
