package webpage

import (
	"html/template"
	"net.vikesh/goshop/dto"
	"net/http"
)

var homePageTemplate *template.Template

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	type page struct {
		Website string
		Title   string
		Meta    map[string]string
		Posts   []dto.ArticleListing
		PageId  string
	}
	var err error
	if homePageTemplate == nil {
		homePageTemplate, err = parseHtml("index", "article_list")
	}
	articles, err := articleService.ListPublishedArticles(nil)
	if articles != nil {
		dispatch(homePageTemplate, err, w, page{Website: "vikesh.net", Title: "Listing", Posts: articles.Listings, PageId: "article-listing"})
	} else {
		http.Redirect(w, r, "/page/not-found", http.StatusTemporaryRedirect)
		return
	}
}
