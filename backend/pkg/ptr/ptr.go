// Package ptr provides utility functions for working with pointers.
package ptr

// To returns a pointer to v.
func To[T any](v T) *T {
	return &v
}
