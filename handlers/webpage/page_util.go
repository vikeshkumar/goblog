package webpage

import (
	"bytes"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"html/template"
	c "net.vikesh/goshop/config"
	"net/http"
	"sync"
)

var templateLock sync.Mutex

func join(pageTemplateParts []string, commonEssentials []string) []string {
	s := make([]string, len(pageTemplateParts)+len(commonEssentials))
	i := 0
	for _, p := range pageTemplateParts {
		s[i] = c.Get().GetString(c.TemplateRoot) + p + c.Get().GetString(c.TemplateSuffix)
		i++
	}
	for _, p := range commonEssentials {
		s[i] = c.Get().GetString(c.TemplateRoot) + p + c.Get().GetString(c.TemplateSuffix)
		i++
	}
	return s
}

func dispatch(t *template.Template, err error, w http.ResponseWriter, data interface{}) {
	if err != nil {
		logError(w, err, nil, http.StatusInternalServerError)
		return
	}
	var b bytes.Buffer
	execError := t.Execute(&b, data)
	if execError != nil {
		logError(w, execError, nil, http.StatusInternalServerError)
	} else {
		_, _ = b.WriteTo(w)
	}
}

func logError(w http.ResponseWriter, e error, data []byte, statusCode int) {
	log.Errorf("%v", e)
	w.WriteHeader(statusCode)
	if data != nil {
		_, _ = w.Write(data)
	}
}

func parseHtml(path ...string) (*template.Template, error) {
	var common = []string{
		"common/footer",
		"common/head",
		"common/header",
		"common/js",
		"common/style",
		"common/sme",
	}
	return template.ParseFiles(join(path, common)...)
}

func AddHandlers(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/index").HandlerFunc(homePageHandler)
	router.Methods(http.MethodGet).Path("/").HandlerFunc(homePageHandler)
	router.Methods(http.MethodGet).Path("/page/{page}").HandlerFunc(contentPageHandler)
	router.Methods(http.MethodGet).Path("/index.html").HandlerFunc(homePageHandler)
	router.Methods(http.MethodGet).Path("/post/{id}").HandlerFunc(articleHandler)
	router.Methods(http.MethodGet).Path("/favicon.ico").HandlerFunc(faviconHandler)
	router.Methods(http.MethodGet).Path("/robots.txt").HandlerFunc(robotsHandler)
}
