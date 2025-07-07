package util

func IfElse[T any](predicate bool, ifTrue T, ifFalse T) T {
	if predicate {
		return ifTrue
	}
	return ifFalse
}

func LazyIfElse[T any](predicate bool, ifTrue func() T, ifFalse func() T) T {
	if predicate {
		return ifTrue()
	}
	return ifFalse()
}
