package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

var locks = newLockmap()

type localFile struct {
	pathname string
	info     os.FileInfo
	descr    *os.File
	mtx      *sync.RWMutex
	read     bool
}

func (f *localFile) Meta() *FileMeta {
	return &FileMeta{
		Name:     f.info.Name(),
		Size:     f.info.Size(),
		UpdateAt: f.info.ModTime(),
	}
}

func (f *localFile) Reader() (io.Reader, error) {
	f.read = true
	f.mtx.RLock()
	return f.descr, nil
}

func (f *localFile) Writer() (io.Writer, error) {
	f.mtx.Lock()
	return f.descr, nil
}

func (f *localFile) Close() {
	if f.read {
		f.mtx.RUnlock()
	} else {
		f.mtx.Unlock()
	}
	if f.descr != nil {
		f.descr.Close()
		f.descr = nil
	}
}

type localFileManager struct {
	workDir string
}

func NewLocalFileManager(workDir string) FileManager {
	return &localFileManager{
		workDir: workDir,
	}
}

func (m *localFileManager) Cwd() string {
	return m.workDir
}

func (m *localFileManager) Exists(path string) bool {
	pathname := m.getFullPath(path)
	if _, err := os.Stat(pathname); os.IsNotExist(err) {
		return false
	}

	return true
}

func (m *localFileManager) Get(path string) (File, error) {
	pathname := m.getFullPath(path)

	d, err := os.OpenFile(pathname, os.O_RDONLY, 0)
	if err != nil {
		return nil, ErrorOpenFile
	}
	fi, _ := d.Stat()
	return &localFile{
		pathname: pathname,
		info:     fi,
		descr:    d,
		mtx:      locks.getMutex("pathname"),
	}, nil
}

func (m *localFileManager) Create(path string) (File, error) {
	pathname := m.getFullPath(path)

	d, err := os.OpenFile(pathname, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, ErrorOpenFile
	}
	fi, _ := d.Stat()
	return &localFile{
		pathname: pathname,
		info:     fi,
		descr:    d,
		mtx:      locks.getMutex(pathname),
	}, nil
}

func (m *localFileManager) List() ([]*FileMeta, error) {
	its, err := ioutil.ReadDir(m.workDir)
	if err != nil {
		return nil, err
	}

	res := make([]*FileMeta, len(its))
	i := 0
	for _, fi := range its {
		if !fi.IsDir() {
			res[i] = &FileMeta{
				Name:     fi.Name(),
				Size:     fi.Size(),
				UpdateAt: fi.ModTime(),
			}
		}
		i++
	}

	return res, nil
}

func (m *localFileManager) getFullPath(filename string) string {
	return path.Join(m.workDir, filename)
}
