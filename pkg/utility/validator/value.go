package validator

// IsValueEmpty checks if a value of any comparable type is empty.
func IsValueEmpty[T comparable](data T) bool {
	var t T
	return data == t
}

// IsValueNotEmpty checks if a value of any comparable type is not empty.
func IsValueNotEmpty[T comparable](data T) bool {
	var t T
	return data != t
}
