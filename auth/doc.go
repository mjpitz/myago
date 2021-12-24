// Package auth provides common code for handling user authentication in a rather implementation agnostic way.
// Currently, we only provide basic auth backed by a CSV, but most components contain an interface that _should_ make it
// rather easy to swap out implementations.
package auth
