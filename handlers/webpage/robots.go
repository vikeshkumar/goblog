package webpage

import (
	logger "github.com/spf13/jwalterweatherman"
	c "net.vikesh/goshop/config"
	"net/http"
	"os"
)

// handler to write robots.txt
func robotsHandler(w http.ResponseWriter, r *http.Request) {
	s, se := os.Stat(c.Get().GetString(c.Robots))
	if se != nil {
		logger.ERROR.Println("error reading stat of file", se)
		w.WriteHeader(404)
		return
	}
	f, e := os.Open(c.Get().GetString(c.Robots))
	if e != nil {
		logger.ERROR.Println("error reading file", e)
		w.WriteHeader(404)
		return
	}

	http.ServeContent(w, r, r.URL.Path, s.ModTime(), f)
}
