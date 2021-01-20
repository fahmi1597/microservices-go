package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

// UploadREST is a handler for uploading a file in restful approach
func (fh *File) UploadREST(resp http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	fn := vars["filename"]

	fh.log.Debug("Handle PUT request", "id", id, "filename", fn)
	fh.saveFiles(id, fn, resp, req.Body)

}

// UploadMultipart is a handler for uploading a file in multipart approach
func (fh *File) UploadMultipart(resp http.ResponseWriter, req *http.Request) {

	err := req.ParseMultipartForm(128 * 1024)
	if err != nil {
		http.Error(resp, "Expected multipart form", http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(req.FormValue("id")); err != nil {
		http.Error(resp, "Expected id integer", http.StatusBadRequest)
		return
	}

	mf, mfh, err := req.FormFile("file")
	if err != nil {
		http.Error(resp, "Expected file", http.StatusBadRequest)
		return
	}

	fh.saveFiles(
		req.FormValue("id"),
		mfh.Filename,
		resp,
		mf,
	)

}

func (fh *File) saveFiles(id string, fn string, resp http.ResponseWriter, mf io.ReadCloser) {
	fh.log.Debug("Saving file", "filename", fn)

	fp := filepath.Join(id, fn)

	if err := fh.store.Write(fp, mf); err != nil {
		fh.log.Debug("Unable to save file", "error", err)
		http.Error(resp, "Unable to save file", http.StatusInternalServerError)
		return
	}
}
