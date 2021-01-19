package webpage

import (
	"github.com/gorilla/mux"
	"html/template"
	"net.vikesh/goshop/db"
	"net.vikesh/goshop/dto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var articleService = &db.ArticleService{}
var articleTemplate *template.Template
func articleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	input := vars["id"]
	var err error
	if articleTemplate == nil {
		articleTemplate, err = parseHtml("post", "article")
	}
	_, inputNotInteger := strconv.Atoi(input)
	var article *dto.Article
	if inputNotInteger == nil {
		article, _ = articleService.GetArticleById(input)
	} else {
		article, _ = articleService.GetArticleByTitle(input)
	}
	if article == nil {
		http.Redirect(w, r, "/page/not-found", http.StatusTemporaryRedirect)
		return
	}
	articleForWeb := parseArticleForWeb(article)
	dispatch(articleTemplate, err, w, articleForWeb)
}

func parseArticleForWeb(article *dto.Article) interface{} {
	type webArticle struct {
	ID            interface{}
	Title         string
	Content       template.HTML
	ModifiedTime  time.Time
	CreationTime  time.Time
	PublishedTime time.Time
}
	var addingLineNumbersClass string
	if len(article.Content) != 0 {
		addingLineNumbersClass = strings.ReplaceAll(article.Content, "<code>", "<code class=\"line-numbers\">")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "<html>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "</html>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "<body>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "</body>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "<head>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "</head>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "<!DOCTYPE html>", "")
	}
	html := template.HTML(addingLineNumbersClass)
	return &webArticle{ID: article.ID, Title: article.Title, Content: html, ModifiedTime: article.ModifiedDate, CreationTime: article.CreationTime, PublishedTime: article.PublishedTime}
}
