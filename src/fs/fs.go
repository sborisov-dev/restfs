package fs

import (
	"errors"
	"io"
	"time"
)

var (
	ErrorFileExists = errors.New("file exists")
	ErrorOpenFile   = errors.New("can't open file")
)

type FileMeta struct {
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	UpdateAt time.Time `json:"updateAt"`
}

type File interface {
	Meta() *FileMeta
	Reader() (io.Reader, error)
	Writer() (io.Writer, error)
	Close()
}

type FileManager interface {
	Cwd() string
	Exists(path string) bool
	Get(path string) (File, error)
	Create(path string) (File, error)
	List() ([]*FileMeta, error)
}
