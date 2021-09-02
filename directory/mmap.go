package directory

import (
	"bytes"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"io"
	"os"
)

type mmapDirectory struct {
	rootPath  string
}

func NewMMapDirectory(rootPath string) *mmapDirectory {
	return &mmapDirectory{
		rootPath: rootPath,
	}
}

func (m *mmapDirectory) Read(path string) (io.ReadCloser, error) {
	f, err := os.Open(m.buildPath(path))
	if err != nil {
		return nil, err
	}
	// should we use mmap for read?
	return f, err
}

func (m *mmapDirectory) OpenWrite(path string) (WriteCloseFlasher, error) {
	f, err := os.OpenFile(m.buildPath(path), os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}
	// Empty file causes mmap error
	if stat, err := f.Stat(); err != nil {
		return nil, err
	} else {
		if stat.Size() == 0 {
			if _, err := f.WriteString("\n"); err != nil {
				return nil, err
			}
		}
	}
	mem, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return nil, err
	}

	return newMmapIO(mem), nil
}

func (m *mmapDirectory) Exists(path string) (bool, error) {
	_, err := os.Stat(m.buildPath(path))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *mmapDirectory) buildPath(path string) string {
	return fmt.Sprintf("%s/%s", m.rootPath, path)
}

type mmapIO struct {
	mmap mmap.MMap
}

func newMmapIO(m mmap.MMap) *mmapIO {
	return &mmapIO{
		mmap: m,
	}
}

func (m *mmapIO) Read(p []byte) (n int, err error) {
	return bytes.NewReader(m.mmap).Read(p)
}

func (m *mmapIO) Write(p []byte) (n int, err error) {
	return bytes.NewBuffer(m.mmap).Write(p)
}

func (m *mmapIO) Flush() error {
	return m.mmap.Flush()
}

func (m *mmapIO) Close() error {
	return m.mmap.Unmap()
}
