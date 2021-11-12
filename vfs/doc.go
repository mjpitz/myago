/*
Package vfs provides utilities for managing virtual file systems on contexts to avoid direct calls to the built-in `os`
interface. This is particularly useful for testing. Currently, this wraps the afero virtual file system which provides
OS and in memory implementations.
*/
package vfs
