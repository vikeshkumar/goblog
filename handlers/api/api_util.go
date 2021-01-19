package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"net.vikesh/goshop/config"
	"net.vikesh/goshop/dto"
	"net/http"
	"runtime"
)

var decoder = schema.NewDecoder()
var cfg = config.Get()

func AddHandlers(r *mux.Router) {
	r.Methods(http.MethodGet).Path("/api/navigations").HandlerFunc(navigationHandler)
	//Upload file
	r.Methods(http.MethodPost).Path("/api/tinymce/upload").HandlerFunc(fileUploadHandler)
	//Articles
	r.Methods(http.MethodPost).Path("/api/new/article").HandlerFunc(newArticleHandler)
	r.Methods(http.MethodPost).Path("/api/article/update").HandlerFunc(updateArticleHandler)
	r.Methods(http.MethodGet).Path("/api/articles/{id}").HandlerFunc(getArticleByIdHandler)
	r.Methods(http.MethodDelete).Path("/api/articles/{id}").HandlerFunc(deleteArticleByIdHandler)
	r.Methods(http.MethodGet, http.MethodPost).Path("/api/articles/").HandlerFunc(listArticlesHandler)
	//Content Pages
	r.Methods(http.MethodPost).Path("/api/new/page").HandlerFunc(newContentHandler)
	r.Methods(http.MethodPost).Path("/api/page/update").HandlerFunc(updateContentPageHandler)
	r.Methods(http.MethodGet).Path("/api/page/{id}").HandlerFunc(getContentPageByIdHandler)
	r.Methods(http.MethodDelete).Path("/api/page/{id}").HandlerFunc(deleteContentPageByIdHandler)
	r.Methods(http.MethodGet, http.MethodPost).Path("/api/pages/").HandlerFunc(listContentPagesHandler)
	//User Management
	r.Methods(http.MethodPost).Path("/api/register").HandlerFunc(registerUserHandler)
	r.Methods(http.MethodPost).Path("/api/authenticate").HandlerFunc(authenticationHandler)
	//Who is the user
	r.Methods(http.MethodGet, http.MethodPost).Path("/api/.whoami").HandlerFunc(whoamiHandler)
}
func logAndPanic(err error, w http.ResponseWriter, status int) {
	if err != nil {
		pc, file, line, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			log.Errorf("error in file : %v, %v:%v : %v", file, details.Name(), line, err)
		}
		w.WriteHeader(status)
		panic(err)
	}
}

func toJson(success *dto.SuccessResponse, w http.ResponseWriter, err *dto.ErrorResponse) {
	if err != nil {
		pc, file, line, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			log.Errorf("error in file : %v, %v:%v : %v", file, details.Name(), line, err.Error)
		}
		b, e := json.MarshalIndent(err, "", "    ")
		if e != nil {
			panic(e)
		}
		w.Header().Add("content-type", "text/json")
		w.WriteHeader(err.Status)
		bytesWritten, writeError := w.Write(b)
		defer logWriteError(writeError, bytesWritten)
	} else if success != nil {
		if success.Result != nil {
			b, e := json.MarshalIndent(success.Result, "", "    ")
			if e != nil {
				panic(e)
			}
			w.Header().Add("content-type", "text/json")
			bytesWritten, writeError := w.Write(b)
			defer logWriteError(writeError, bytesWritten)
		} else {
			w.WriteHeader(success.Status)
		}
	}
}

func logWriteError(writeError error, bytesWritten int) {
	if writeError != nil {
		log.Errorf("error in writing: %v", writeError)
	}
	log.Debugf("written %b bytes", bytesWritten)
}

func decode(payload interface{}, r *http.Request) error {
	err := r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	if err != nil {
		return errors.New("received unsupported data")
	}
	decodeError := decoder.Decode(payload)
	if decodeError != nil {
		return errors.New("received unsupported data")
	}
	return nil
}

func success(result interface{}) *dto.SuccessResponse {
	if result != nil {
		return &dto.SuccessResponse{
			Status: http.StatusOK,
			Result: result,
		}
	}
	return &dto.SuccessResponse{
		Status: http.StatusOK,
	}
}

func successWithCode(result interface{}, code int) *dto.SuccessResponse {
	return &dto.SuccessResponse{
		Status: code,
		Result: result,
	}
}

func failed(err error) *dto.ErrorResponse {
	if err != nil {
		return &dto.ErrorResponse{
			Error:    err,
			Status:   http.StatusInternalServerError,
			Response: err.Error(),
		}
	}
	return nil
}

func failedWithStatus(err error, code int) *dto.ErrorResponse {
	if err != nil {
		return &dto.ErrorResponse{
			Error:    err,
			Status:   code,
			Response: err.Error(),
		}
	}
	return nil
}
