package webpage

import (
	"github.com/gorilla/mux"
	"html/template"
	"net.vikesh/goshop/db"
	"net.vikesh/goshop/dto"
	"net/http"
	"strings"
)

var contentPageService = &db.ContentPageService{}
var contentPageTemplate *template.Template

func contentPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	input := vars["page"]
	var err error
	if contentPageTemplate == nil {
		articleTemplate, err = parseHtml("page", "page_content")
	}
	page, _ := contentPageService.GetContentPageByUrl(input)
	if page == nil {
		http.Redirect(w, r, "/page/not-found", http.StatusTemporaryRedirect)
		return
	}
	articleForWeb := parseContentPageForWeb(page)
	dispatch(articleTemplate, err, w, articleForWeb)
}

func parseContentPageForWeb(page *dto.Page) interface{} {
	type webArticle struct {
		ID      interface{}
		Title   string
		Content template.HTML
	}
	var addingLineNumbersClass string
	if len(page.Content) != 0 {
		addingLineNumbersClass = strings.ReplaceAll(page.Content, "<code>", "<code class=\"line-numbers\">")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "<html>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "</html>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "<body>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "</body>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "<head>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "</head>", "")
		addingLineNumbersClass = strings.ReplaceAll(addingLineNumbersClass, "<!DOCTYPE html>", "")
	}
	html := template.HTML(addingLineNumbersClass)
	return &webArticle{ID: page.ID, Title: page.Title, Content: html}
}
