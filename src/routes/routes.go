package routes

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"rest-fs/src/api"
	"rest-fs/src/fs"
	"strconv"
)

func ListFiles(manager fs.FileManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		files, err := manager.List()
		if err != nil {
			api.NewInternalServerErrorResponse().Write(w)
			return
		}
		api.
			NewFileListResponse(files).
			Write(w)
	}
}

func ReadFileMeta(manager fs.FileManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filename := ps.ByName("name")
		if !manager.Exists(filename) {
			api.NewFileNotFoundResponse(filename).Write(w)
			return
		}
		file, err := manager.Get(filename)
		if err != nil {
			api.
				NewInternalServerErrorResponse().
				Write(w)
			return
		}
		api.
			NewFileMetaResponse(file.Meta()).
			Write(w)
	}
}

func StreamFile(manager fs.FileManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filename := ps.ByName("name")
		if !manager.Exists(filename) {
			api.NewFileNotFoundResponse(filename).Write(w)
			return
		}
		file, err := manager.Get(filename)
		if err != nil {
			api.
				NewInternalServerErrorResponse().
				Write(w)
			return
		}
		dwld, _ := strconv.ParseBool(r.URL.Query().Get("download"))
		api.
			NewFileStreamResponse(file, dwld).
			Write(w)
	}
}

func CreateFile(manager fs.FileManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filename := ps.ByName("name")

		if manager.Exists(filename) {
			api.
				NewFileCreateErrorResponse(filename).
				Write(w)
			return
		}

		if err := writeFile(filename, manager, r.Body); err != nil {
			log.Printf("create file error: %s", err)
			api.
				NewInternalServerErrorResponse().
				Write(w)
			return
		}

		api.
			NewOkResponse().
			Write(w)
	}
}

func UpsertFile(manager fs.FileManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filename := ps.ByName("name")

		if !manager.Exists(filename) {
			api.
				NewFileUpdateErrorResponse(filename).
				Write(w)
			return
		}

		if err := writeFile(filename, manager, r.Body); err != nil {
			log.Printf("create file error: %s", err)
			api.
				NewInternalServerErrorResponse().
				Write(w)
			return
		}

		api.
			NewOkResponse().
			Write(w)
	}
}

func NewRouter(manager fs.FileManager) http.Handler {
	router := httprouter.New()

	router.GET("/", ListFiles(manager))
	router.GET("/meta/:name", ReadFileMeta(manager))
	router.GET("/file/:name", StreamFile(manager))
	router.PUT("/file/:name", CreateFile(manager))
	router.POST("/file/:name", UpsertFile(manager))

	return router
}

func writeFile(fname string, m fs.FileManager, payload io.ReadCloser) error {
	f, err := m.Create(fname)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
		payload.Close()
	}()

	dest, err := f.Writer()
	if err != nil {
		return err
	}

	_, err = io.Copy(dest, payload)
	if err != nil {
		return err
	}

	return nil
}
