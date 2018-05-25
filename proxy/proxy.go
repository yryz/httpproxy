package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/yryz/httpproxy/config"

	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	log "github.com/sirupsen/logrus"
)

type ProxyServer struct {
	ssServer string
	ssCipher *ss.Cipher
}

func NewProxyServer() *http.Server {
	p := &ProxyServer{
		ssServer: config.Conf.SsServer,
	}

	var err error
	p.ssCipher, err = ss.NewCipher(config.Conf.SsCipher, config.Conf.SsPassword)
	if err != nil {
		panic("init cipher error: " + err.Error())
	}

	return &http.Server{
		Addr:    config.Conf.Listen,
		Handler: p,
	}
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			log.Debugf("panic: %v\n", err)
		}
	}()

	if r.Method == "CONNECT" {
		p.HandleConnect(w, r)
	} else {
		p.HandleHttp(w, r)
	}
}

// 处理HTTPS、HTTP2代理请求
func (p *ProxyServer) HandleConnect(w http.ResponseWriter, r *http.Request) {
	log.Infof("%s %s", r.Method, r.Host)

	hj, _ := w.(http.Hijacker)
	conn, _, err := hj.Hijack()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ssConn, err := ss.Dial(r.URL.Host, p.ssServer, p.ssCipher.Copy())
	if err != nil {
		log.Error("ss dial: ", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	go ss.PipeThenClose(conn, ssConn, nil)
	ss.PipeThenClose(ssConn, conn, nil)
}

// 处理HTTP代理请求
func (p *ProxyServer) HandleHttp(w http.ResponseWriter, r *http.Request) {
	log.Infof("%s %s", r.Method, r.URL)

	// ss proxy
	tr := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			log.Infof("dial ss %v/%v", addr, network)
			return ss.Dial(addr, p.ssServer, p.ssCipher.Copy())
		},
	}

	// transport
	resp, err := tr.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("request error: ", err)
		return
	}
	defer resp.Body.Close()

	// copy headers
	for k, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// copy body
	n, err := io.Copy(w, resp.Body)
	if err != nil && err != io.EOF {
		log.Errorf("copy response body error: %v", err)
	}

	log.Infof("copied %v bytes from %v.", n, r.Host)
}
