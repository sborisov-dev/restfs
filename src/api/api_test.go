package api

import (
	"net/http"
	"testing"
)

const (
	okResult = "\"ok\""
)

func TestNewOkResponse(t *testing.T) {
	r := NewOkResponse()
	sr, ok := r.(*SuccessResponse)
	if !ok {
		t.Errorf("instance should be *SucccessResponse")
	} else if string([]byte(*sr.Result)) != okResult {
		t.Error("result must be 'ok'")
	}
}

func TestNewFileMetaResponse(t *testing.T) {
	r := NewFileMetaResponse(nil)
	sr, ok := r.(*SuccessResponse)
	if !ok {
		t.Errorf("instance should be *SucccessResponse")
	} else if sr.StatusCode != http.StatusOK {
		t.Error("status must be 200")
	}
}

func TestNewFileListResponse(t *testing.T) {
	r := NewFileListResponse(nil)
	sr, ok := r.(*SuccessResponse)
	if !ok {
		t.Errorf("instance should be *SucccessResponse")
	} else if sr.StatusCode != http.StatusOK {
		t.Error("status must be 200")
	}
}

func TestNewInternalServerErrorResponse(t *testing.T) {
	r := NewInternalServerErrorResponse()
	er, ok := r.(*ErrorResponse)
	if !ok {
		t.Errorf("instance should be *ErrorResponse")
	} else if er.StatusCode != http.StatusInternalServerError {
		t.Error("status must be 500")
	} else if er.ErrorCode != http.StatusInternalServerError {
		t.Error("error code must be 500")
	} else if er.ErrorMessage != "internal server error" {
		t.Error("incorrect error message")
	}
}

func TestNewFileNotFoundResponse(t *testing.T) {
	r := NewFileNotFoundResponse("f.ext")
	er, ok := r.(*ErrorResponse)
	if !ok {
		t.Errorf("instance should be *ErrorResponse")
	} else if er.StatusCode != http.StatusNotFound {
		t.Error("status must be 404")
	} else if er.ErrorCode != ErrorCodeFileNotFound {
		t.Error("error code must be ErrorCodeFileNotFound")
	} else if er.ErrorMessage != "File 'f.ext' not found" {
		t.Error("incorrect error message")
	}
}

func TestNewFileCreateResponse(t *testing.T) {
	r := NewFileCreateErrorResponse("f.ext")
	er, ok := r.(*ErrorResponse)
	if !ok {
		t.Errorf("instance should be *ErrorResponse")
	} else if er.StatusCode != http.StatusInternalServerError {
		t.Error("status must be 500")
	} else if er.ErrorCode != ErrorCodeFileCreateFail {
		t.Error("error code must be ErrorCodeFileCreateFail")
	} else if er.ErrorMessage != "Can't create file 'f.ext'" {
		t.Error("incorrect error message")
	}
}

func TestNewFileUpdateResponse(t *testing.T) {
	r := NewFileUpdateErrorResponse("f.ext")
	er, ok := r.(*ErrorResponse)
	if !ok {
		t.Errorf("instance should be *ErrorResponse")
	} else if er.StatusCode != http.StatusInternalServerError {
		t.Error("status must be 500")
	} else if er.ErrorCode != ErrorCodeFileUpdateFail {
		t.Error("error code must be ErrorCodeFileUpdateFail")
	} else if er.ErrorMessage != "Can't update file 'f.ext'" {
		t.Error("incorrect error message")
	}
}

func TestNewFileStreamResponse(t *testing.T) {
	r := NewFileStreamResponse(nil, false)
	_, ok := r.(*BinaryResponse)
	if !ok {
		t.Errorf("instance should be *BinaryResponse")
	}
}
