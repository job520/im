package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"im/http/config"
	"im/http/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal, 1) // 退出信号
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	r := router.NewRouter()
	conf := config.Config
	srv := &http.Server{
		Addr:    conf.Server.Address,
		Handler: r,
	}
	go shutdown(quit, srv)
	logrus.Infof("server running at:%s \n", conf.Server.Address)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatal(err)
	}
}

func shutdown(quit chan os.Signal, srv *http.Server) {
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal(err)
	}
	os.Exit(0)
}
