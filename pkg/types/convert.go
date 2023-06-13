package types

func AsPointer[T any](v T) *T {
	return &v
}

func AsValue[T any](v *T) T {
	if v != nil {
		return *v
	}
	var empty T
	return empty
}
