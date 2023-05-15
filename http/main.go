package main

import (
	"context"
	"im/http/middleware"
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
	r.Use(middleware.Test())
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go shutdown(quit, srv)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func shutdown(quit chan os.Signal, srv *http.Server) {
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	os.Exit(0)
}
