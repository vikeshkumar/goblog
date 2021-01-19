package api

import (
	"github.com/gorilla/mux"
	"net.vikesh/goshop/db"
	"net.vikesh/goshop/dto"
	"net/http"
)

var pageService = &db.ContentPageService{}
func newContentHandler(w http.ResponseWriter, r *http.Request) {
	newPage, err := pageService.NewContentPage()
	toJson(successWithCode(newPage, http.StatusCreated), w, failed(err))
}

func updateContentPageHandler(w http.ResponseWriter, r *http.Request) {
	page := &dto.Page{}
	err := decode(page, r)
	if err != nil {
		toJson(success(nil), w, failedWithStatus(err, http.StatusNotAcceptable))
	} else {
		toJson(success(nil), w, failed(pageService.SaveContentPage(page)))
	}
}

func listContentPagesHandler(w http.ResponseWriter, r *http.Request) {
	request := &dto.SearchParams{}
	err := decode(request, r)
	if err != nil {
		toJson(success(nil), w, failedWithStatus(err, http.StatusNotAcceptable))
	} else {
		pages, err := pageService.ListContentPages(request)
		toJson(success(pages), w, failed(err))
	}
}

func getContentPageByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleId := vars["id"]
	pages, err := pageService.GetContentPageById(articleId)
	toJson(success(pages), w, failed(err))
}

func deleteContentPageByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleId := vars["id"]
	err := pageService.DeleteContentPageById(articleId)
	toJson(success(nil), w, failed(err))
}
