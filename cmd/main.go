package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crazyfrankie/kube-ctl/conf"
	"github.com/crazyfrankie/kube-ctl/ioc"
)

func main() {
	engine := ioc.InitServer()

	srv := &http.Server{
		Handler: engine,
		Addr:    conf.GetConf().Server.Addr,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server Listen: %s\n", err)
		}
	}()
	log.Println(fmt.Sprintf("Server is running at http://localhost%s", srv.Addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
