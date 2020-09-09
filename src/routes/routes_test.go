package routes

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"rest-fs/src/fs"
	"testing"
)

func createManager() fs.FileManager {
	m := fs.NewMockFileManager()
	f, _ := m.Create("test.txt")
	w, _ := f.Writer()
	w.Write([]byte("test content"))

	return m
}

func checkOkResponse(t *testing.T, result *http.Response) {
	if result.StatusCode != http.StatusOK {
		t.Errorf("status code must be 200")
	} else if result.Header.Get("Content-Type") != "application/json" {
		t.Errorf("incorrent response content-type")
	}
}

func TestNewRouter(t *testing.T) {
	r := NewRouter(createManager())
	_, ok := r.(*httprouter.Router)
	if !ok {
		t.Errorf("instanse must be *httprouter.Router")
	}

}

func TestRedirectTraillingSlash(t *testing.T) {
	m := createManager()
	router := NewRouter(m)
	r := httptest.NewRequest("GET", "http://restfs.example.com", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	result := w.Result()
	if result.StatusCode != http.StatusMovedPermanently {
		t.Errorf("status code must be 301")
	} else if result.Header.Get("Location") != "http://restfs.example.com/" {
		t.Errorf("location must be http://restfs.example.com/")
	}
}

func TestListFiles(t *testing.T) {
	m := createManager()
	router := NewRouter(m)
	r := httptest.NewRequest("GET", "http://restfs.example.com/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	checkOkResponse(t, w.Result())
}

func TestReadFileMeta(t *testing.T) {
	m := createManager()
	router := NewRouter(m)

	t.Run("not found file meta", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://restfs.example.com/meta/notfound.ext", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		result := w.Result()
		if result.StatusCode != http.StatusNotFound {
			t.Errorf("status code must be 404")
		}
	})

	t.Run("get file meta", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://restfs.example.com/meta/test.txt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		checkOkResponse(t, w.Result())
	})
}

func TestReadFileData(t *testing.T) {
	m := createManager()
	router := NewRouter(m)

	t.Run("not found file", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://restfs.example.com/file/notfound.ext", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		result := w.Result()
		if result.StatusCode != http.StatusNotFound {
			t.Errorf("status code must be 404")
		}
	})

	t.Run("get file", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://restfs.example.com/file/test.txt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code must be 200")
		} else if resp.ContentLength != 12 {
			t.Errorf("bad content length: %d", resp.ContentLength)
		} else if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("bad content type: %s", resp.Header.Get("Content-Type"))
		}
	})

	t.Run("get file as attach", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://restfs.example.com/file/test.txt?download=1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code must be 200")
		} else if resp.ContentLength != 12 {
			t.Errorf("bad content length: %d", resp.ContentLength)
		} else if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("bad content type: %s", resp.Header.Get("Content-Type"))
		} else if resp.Header.Get("Content-Disposition") != "attachment; filename=test.txt" {
			t.Errorf("should be header Content-Disposition: attachment; filename=test1.txt")
		}
	})
}
