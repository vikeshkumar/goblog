package main

import (
	"context"
	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net.vikesh/goshop/config"
	"net.vikesh/goshop/db"
	"net.vikesh/goshop/handlers/api"
	"net.vikesh/goshop/handlers/webpage"
	"net.vikesh/goshop/middlewares"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	cfg, err := config.Init("./config.toml", ".", "toml")
	if err != nil {
		log.Errorf("error reading required configuration file %v", err)
		os.Exit(1)
	}
	router := mux.NewRouter()
	connectError := db.Connect()
	if connectError != nil {
		log.Fatal("failed to connect to database", err)
	}
	webpage.AddHandlers(router)
	api.AddHandlers(router)
	deploymentDir := cfg.GetString(config.ServerExecDir)
	//admin
	adminPathDir := deploymentDir + cfg.GetString(config.AdminResourcePath)
	router.PathPrefix(cfg.GetString(config.AdminPathPrefix)).Handler(
		http.StripPrefix(cfg.GetString(config.AdminPathPrefix),
			http.FileServer(http.Dir(adminPathDir))))
	//static path
	staticPathDir := deploymentDir + cfg.GetString(config.StaticResourcePath)
	router.PathPrefix(cfg.GetString(config.StaticPathPrefix)).Handler(
		http.StripPrefix(cfg.GetString(config.StaticPathPrefix),
			http.FileServer(http.Dir(staticPathDir))))
	//uploads path
	uploadPaths := cfg.GetString(config.FileUploadDirectory)
	if _, err := os.Stat(uploadPaths); os.IsNotExist(err) {
		_ = os.Mkdir(uploadPaths, os.ModeDir)
	}
	router.PathPrefix(cfg.GetString(config.UploadPathPrefix)).Handler(
		http.StripPrefix(cfg.GetString(config.UploadPathPrefix),
			http.FileServer(http.Dir(uploadPaths))))
	router.Schemes("http")
	//configure a logging middleware
	router.Use(middlewares.LoggingMiddleware, middlewares.PerformanceMonitor, middlewares.Secured)
	var webPort string
	if len(port) == 0 {
		webPort = cfg.GetString(config.ServerListen)
	}
	server := &http.Server{
		Handler: gh.CompressHandler(router),
		Addr:    webPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: cfg.GetDuration(config.ServerWriteTimeout) * time.Second,
		ReadTimeout:  cfg.GetDuration(config.ServerReadTimeout) * time.Second,
	}
	log.Infof("server listening on %v", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Infof("err = %v", err)
	}

	// accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	listenForShutDown(ch, server, time.Until(time.Now().Add(time.Minute*1)), func() {
		log.Info("closing connection to db")
	})
}

// listen for application shutdown and configure timeout to respond to request
// and then terminate
func listenForShutDown(ch chan os.Signal, server *http.Server, wait time.Duration, callback func()) {
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	callback()
	defer cancel()
	log.Infof("shutting down, received signal = %v", ctx)
	_ = server.Shutdown(ctx)
	os.Exit(0)
}
