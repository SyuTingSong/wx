package main

import (
	"context"
	"github.com/apsdehal/go-logger"
	"github.com/gin-gonic/gin"
	"github.com/syutingsong/wx/config"
	"github.com/syutingsong/wx/web"
	"github.com/syutingsong/wx/web/login"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var _ = config.ParseConfig()

func main() {
	var log, _ = logger.New(
		config.Global.LogLevel,
		config.Global.LogColor,
	)

	log.NoticeF("Server start with %v", config.Global)
	r := gin.Default()
	r.GET("/", login.L2)
	rLogin := r.Group("/login")
	{
		rLogin.GET("/l2", login.L2)
	}

	r.GET("/access_token", web.AccessToken)

	server := &http.Server{
		Addr:    config.Global.Addr.String(),
		Handler: r,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Notice("Shutdown server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		log.FatalF("Server shutdown error %s", err.Error())
	}
	cancel()
	log.Notice("Exit")
}
