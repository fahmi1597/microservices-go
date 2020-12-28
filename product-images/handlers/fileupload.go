package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

// UploadREST is a handler for uploading a file in restful way
func (fh *File) UploadREST(resp http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	fn := vars["filename"]

	fh.log.Printf("[INFO] Handle PUT request: id=%s filename=%s", id, fn)
	fh.saveFiles(id, fn, resp, req.Body)

}

// UploadMultipart is a handler for uploading a file in multipart way
func (fh *File) UploadMultipart(resp http.ResponseWriter, req *http.Request) {

	err := req.ParseMultipartForm(128 * 1024)
	if err != nil {
		fh.log.Println("[ERROR] Bad request:", err)
		http.Error(resp, "Expected multipart form", http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(req.FormValue("id")); err != nil {
		fh.log.Println("[ERROR] Bad request:", err)
		http.Error(resp, "Expected id integer", http.StatusBadRequest)
		return
	}

	mf, mfh, err := req.FormFile("file")
	if err != nil {
		fh.log.Println("[ERROR] Bad request:", err)
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
	fh.log.Println("[INFO] Saving file")

	fp := filepath.Join(id, fn)

	if err := fh.store.Write(fp, mf); err != nil {
		fh.log.Println("[ERROR] Unable to save file", err)
		http.Error(resp, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fh.log.Printf("[INFO] File saved: %s", fp)
}
