package main

import (
	"os"
	"syscall"

	"github.com/yryz/httpproxy/config"
	"github.com/yryz/httpproxy/proxy"

	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "set" {
			proxyUrl := "http://" + config.Conf.Listen
			os.Setenv("http_proxy", proxyUrl)
			os.Setenv("https_proxy", proxyUrl)
			log.Info("set http_proxy & https_proxy success!")
			syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, syscall.Environ())
			return
		}
	}

	proxyServer := proxy.NewProxyServer()
	log.Infof("listen on %s", config.Conf.Listen)
	proxyServer.ListenAndServe()
}
