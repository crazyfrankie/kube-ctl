package main

import (
	"github.com/crazyfrankie/kube-ctl/conf"
	"github.com/crazyfrankie/kube-ctl/ioc"
)

func main() {
	srv := ioc.InitServer()

	srv.Run(conf.GetConf().Server.Addr)
}
