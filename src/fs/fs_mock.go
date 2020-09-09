package fs

import (
	"bytes"
	"io"
	"time"
)

type fileMock struct {
	meta *FileMeta
	data *bytes.Buffer
}

func newFileMock(name string) *fileMock {
	t := time.Now()
	meta := &FileMeta{
		Name:     name,
		Size:     0,
		UpdateAt: t,
	}

	return &fileMock{
		meta: meta,
		data: bytes.NewBuffer(make([]byte, 0)),
	}
}

func (f *fileMock) Meta() *FileMeta {
	return &FileMeta{
		Name:     f.meta.Name,
		Size:     int64(f.data.Len()),
		UpdateAt: f.meta.UpdateAt,
	}
}

func (f *fileMock) Reader() (io.Reader, error) {
	//var buf bytes.Buffer
	return bytes.NewReader(f.data.Bytes()), nil
}

func (f *fileMock) Writer() (io.Writer, error) {
	return f.data, nil
}

func (f *fileMock) Close() {

}

type mockFileManager struct {
	files map[string]*fileMock
}

func NewMockFileManager() FileManager {
	return &mockFileManager{
		files: make(map[string]*fileMock),
	}
}

func (m *mockFileManager) Cwd() string {
	return "/mock"
}

func (m *mockFileManager) Exists(path string) bool {
	_, exists := m.files[path]

	return exists
}

func (m *mockFileManager) Create(path string) (File, error) {
	if m.Exists(path) {
		return nil, ErrorFileExists
	}
	f := newFileMock(path)
	m.files[path] = f
	return f, nil
}

func (m *mockFileManager) Get(path string) (File, error) {
	if !m.Exists(path) {
		return nil, ErrorOpenFile
	}
	f := m.files[path]
	return f, nil
}

func (m *mockFileManager) List() ([]*FileMeta, error) {
	size := len(m.files)
	res := make([]*FileMeta, size, size)

	idx := 0
	for _, f := range m.files {
		res[idx] = f.Meta()
		idx++
	}

	return res, nil
}
