package webpage

import (
	logger "github.com/spf13/jwalterweatherman"
	c "net.vikesh/goshop/config"
	"net/http"
	"os"
)

// handler to write favicon file
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	s, se := os.Stat(c.Get().GetString(c.Favicon))
	if se != nil {
		logger.ERROR.Println("error reading stat of file", se)
		w.WriteHeader(404)
		return
	}
	f, e := os.Open(c.Get().GetString(c.Favicon))
	if e != nil {
		logger.ERROR.Println("error reading file", e)
		w.WriteHeader(404)
		return
	}

	http.ServeContent(w, r, r.URL.Path, s.ModTime(), f)
}

