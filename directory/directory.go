package directory

import (
	"io"
)

type WriteCloseFlasher interface {
	io.WriteCloser
	Flush() error
}

type Directory interface {
	// TODO: Should be able to seek
	Read(path string) (io.ReadCloser, error)
	OpenWrite(path string) (WriteCloseFlasher, error)
	Exists(path string) (bool, error)
}

