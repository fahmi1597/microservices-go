package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

// Upload is a handler for uploading a file
func (fh *File) Upload(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	fn := vars["filename"]

	fh.log.Println("[INFO] Handle PUT request", id, fn)
	fh.saveFiles(id, fn, resp, req)
}

func (fh *File) saveFiles(id string, fn string, resp http.ResponseWriter, req *http.Request) {
	fh.log.Println("[INFO] Saving file")

	fp := filepath.Join(id, fn)

	if err := fh.store.Write(fp, req.Body); err != nil {
		fh.log.Println("[ERROR] Unable to save file", err)
		http.Error(resp, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fh.log.Printf("[INFO] File saved: %s", fp)
}
