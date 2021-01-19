package api

import (
	"bytes"
	"github.com/google/uuid"
	"io"
	c "net.vikesh/goshop/config"
	"net/http"
	"os"
	"strings"
)

const maxFileSize = 8192 * 10

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Location string `json:"location"`
	}
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		logAndPanic(err, w, http.StatusInternalServerError)
	} else {
		file, header, fileError := r.FormFile("file")
		var mimeHeader string
		if header != nil {
			mimeHeader = strings.Split(header.Header.Get("Content-Type"), "/")[1]
		}
		if fileError != nil {
			logAndPanic(fileError, w, http.StatusInternalServerError)
		} else {
			defer file.Close()
			fileName := uuid.New().String()
			buf := bytes.Buffer{}
			_, fileError := io.Copy(&buf, file)
			if fileError != nil {
				logAndPanic(fileError, w, http.StatusInternalServerError)
			}
			fileUploadDir := c.Get().GetString(c.FileUploadDirectory)
			newFile, fileCreateError := os.Create(fileUploadDir + fileName + "." + mimeHeader)
			if fileCreateError != nil {
				logAndPanic(fileCreateError, w, http.StatusInternalServerError)
			}
			defer newFile.Close()
			newFile.Write(buf.Bytes())
			newFile.Sync()
			location := c.Get().GetString(c.ServerAppUrl) + c.Get().GetString(c.ServerUploadPath)
			response := &response{location + fileName + "." + mimeHeader}
			toJson(success(response), w, nil)
		}
	}
}
