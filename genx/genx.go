package genx

// func IsZero[T any](v T) bool {
// 	var zero T
//
// 	return v == zero
// }

func Zero[T any]() T {
	var zero T

	return zero
}
