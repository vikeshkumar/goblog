package api

import (
	"github.com/gorilla/mux"
	"net.vikesh/goshop/db"
	"net.vikesh/goshop/dto"
	"net/http"
)

var articleService = &db.ArticleService{}
func newArticleHandler(w http.ResponseWriter, r *http.Request) {
	newArticle, err := articleService.NewArticle()
	toJson(successWithCode(newArticle, http.StatusCreated), w, failed(err))
}

func updateArticleHandler(w http.ResponseWriter, r *http.Request) {
	article := &dto.Article{}
	err := decode(article, r)
	if err != nil {
		toJson(success(nil), w, failedWithStatus(err, http.StatusNotAcceptable))
	} else {
		toJson(success(nil), w, failed(articleService.SaveArticle(article)))
	}
}

func listArticlesHandler(w http.ResponseWriter, r *http.Request) {
	request := &dto.SearchParams{}
	err := decode(request, r)
	if err != nil {
		toJson(success(nil), w, failedWithStatus(err, http.StatusNotAcceptable))
	} else {
		articles, err := articleService.ListArticles(request)
		toJson(success(articles), w, failed(err))
	}
}

func getArticleByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleId := vars["id"]
	articles, err := articleService.GetArticleById(articleId)
	toJson(success(articles), w, failed(err))
}

func deleteArticleByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleId := vars["id"]
	err := articleService.DeleteArticleById(articleId)
	toJson(success(nil), w, failed(err))
}
