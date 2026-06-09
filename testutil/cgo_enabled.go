//go:build cgo

package testutil

// CgoEnabled reports whether this binary was built with CGO enabled.
func CgoEnabled() bool { return true }
