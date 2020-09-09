package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"rest-fs/src/fs"
	"strconv"
)

const (
	ErrorCodeMarshal        = 1
	ErrorCodeFileNotFound   = 2
	ErrorCodeFileCreateFail = 3
	ErrorCodeFileUpdateFail = 4
)

type Response interface {
	Write(dest http.ResponseWriter)
}

type SuccessResponse struct {
	StatusCode int              `json:"-"`
	Result     *json.RawMessage `json:"result,omitempty"`
}

func (r *SuccessResponse) Write(dest http.ResponseWriter) {
	payload, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("marshal response error: %s", err)
		newErrorResponse(ErrorCodeMarshal, err.Error()).Write(dest)
		return
	}

	dest.Header().Set("Content-Type", "application/json")
	dest.WriteHeader(r.StatusCode)
	if _, err := dest.Write(payload); err != nil {
		fmt.Printf("write payload error: %s", err)
	}
}

type ErrorResponse struct {
	StatusCode   int    `json:"-"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (r *ErrorResponse) Write(dest http.ResponseWriter) {
	dest.Header().Set("Content-Type", "application/json")
	dest.WriteHeader(r.StatusCode)
	p, _ := json.Marshal(r)
	if _, err := dest.Write(p); err != nil {
		fmt.Printf("error response write error: %s", err)
	}
}

type BinaryResponse struct {
	file     fs.File
	download bool
}

func (r *BinaryResponse) Write(dest http.ResponseWriter) {
	defer r.file.Close()

	meta := r.file.Meta()
	headers := dest.Header()
	headers.Set("Last-Modified", meta.UpdateAt.UTC().Format(http.TimeFormat))
	headers.Set("Content-Type", mime.TypeByExtension(path.Ext(meta.Name)))
	headers.Set("Content-Length", strconv.FormatInt(meta.Size, 10))
	if r.download {
		headers.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", meta.Name))
	}
	content, err := r.file.Reader()
	if err != nil {
		log.Printf("write binary error: %s", err)
		NewInternalServerErrorResponse().Write(dest)
		return
	}
	io.CopyN(dest, content, meta.Size)
}

func newErrorResponseWithStatus(status, code int, message string) Response {
	return &ErrorResponse{status, code, message}
}

func newErrorResponse(code int, message string) Response {
	return newErrorResponseWithStatus(http.StatusInternalServerError, code, message)
}

func newSuccessResponse(res json.RawMessage) Response {
	return &SuccessResponse{
		StatusCode: http.StatusOK,
		Result:     &res,
	}
}

func NewOkResponse() Response {
	return newSuccessResponse(json.RawMessage("\"ok\""))
}

func NewFileMetaResponse(meta *fs.FileMeta) Response {
	res, _ := json.Marshal(meta)

	return newSuccessResponse(json.RawMessage(res))
}

func NewFileListResponse(list []*fs.FileMeta) Response {
	res, _ := json.Marshal(list)

	return newSuccessResponse(json.RawMessage(res))
}

func NewInternalServerErrorResponse() Response {
	return newErrorResponse(http.StatusInternalServerError, "internal server error")
}

func NewFileNotFoundResponse(path string) Response {
	return newErrorResponseWithStatus(http.StatusNotFound,
		ErrorCodeFileNotFound,
		fmt.Sprintf("File '%s' not found", path))
}

func NewFileCreateErrorResponse(path string) Response {
	return newErrorResponse(ErrorCodeFileCreateFail,
		fmt.Sprintf("Can't create file '%s'", path))
}

func NewFileUpdateErrorResponse(path string) Response {
	return newErrorResponse(ErrorCodeFileUpdateFail,
		fmt.Sprintf("Can't update file '%s'", path))
}

func NewFileStreamResponse(f fs.File, download bool) Response {
	return &BinaryResponse{file: f, download: download}
}
