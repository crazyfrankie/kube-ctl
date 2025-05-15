package main

import (
	"context"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/crazyfrankie/kube-ctl/conf"
	"github.com/crazyfrankie/kube-ctl/ioc"
)

func main() {
	app := ioc.InitApp()

	prometheus.MustRegister(app.Metrics)

	g := &run.Group{}

	g.Add(func() error {
		http.Handle("/metrics", promhttp.Handler())
		return http.ListenAndServe("0.0.0.0:8082", nil)
	}, func(err error) {
		//
	})

	srv := &http.Server{
		Handler: app.Engine,
		Addr:    conf.GetConf().Server.Addr,
	}

	g.Add(func() error {
		log.Println("Server is running at http://localhost:8083")
		return srv.ListenAndServe()
	}, func(err error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown main server: %v", err)
		}
	})

	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	if err := g.Run(); err != nil {
		log.Printf("program interrupted, err:%s", err)
	}
}
