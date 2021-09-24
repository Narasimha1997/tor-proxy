package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cretz/bine/tor"
	"github.com/ipsn/go-libtor"
	"gopkg.in/elazarl/goproxy.v1"
)

// DefaultProxyPort default port to listen on
const DefaultProxyPort string = "8000"

// TorHandle Contains the context related information of Tor Connection
type TorHandle struct {
	torCtx       *tor.Tor
	torDialer    *tor.Dialer
	torTransport *http.Transport
}

// initTorHandle Create and return a pointer to new TorHandle
func initTorHandle() *TorHandle {

	ctx := context.Background()
	torCtx, err := tor.Start(
		ctx,
		&tor.StartConf{ProcessCreator: libtor.Creator, DebugWriter: os.Stderr},
	)
	if err != nil {
		log.Fatalf("Failed to create Tor Context = %v\n", err)
	}

	torHandle := TorHandle{}
	torHandle.torCtx = torCtx

	// create a new HTTP client based on Tor Context:
	dialer, err := torHandle.torCtx.Dialer(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create dialer for Tor Context - %v\n", err)
	}

	torHandle.torDialer = dialer

	torHandle.torTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyFromEnvironment,
		DialContext:     torHandle.torDialer.DialContext,
	}

	log.Println("Prepared Tor Handle")

	return &torHandle
}

// clearHandle clears up tor handle, just calls torCtx.Close()
func clearHandle(torHandle *TorHandle) {
	if torHandle != nil {
		torHandle.torCtx.Close()
		log.Println("Cleared Tor Context")
	}
}

type HTTPServing struct{}

func (httpProxy *HTTPServing) getProxy() (*goproxy.ProxyHttpServer, *TorHandle) {
	proxy := goproxy.NewProxyHttpServer()
	tcx := initTorHandle()

	proxy.Tr = tcx.torTransport
	proxy.ConnectDial = tcx.torTransport.Dial
	proxy.Verbose = true

	return proxy, tcx
}

// ListenAndServe start http-proxy server on a given port
func (httpProxy *HTTPServing) ListenAndServe() {
	port, exists := os.LookupEnv("PROXY_PORT")
	if !exists {
		port = DefaultProxyPort
	}

	proxy, tcx := httpProxy.getProxy()
	defer clearHandle(tcx)

	listenString := fmt.Sprintf(":%s", port)
	err := http.ListenAndServe(listenString, proxy)

	if err != nil {
		log.Fatalf("Failed to start proxy server on %s, reason - %v\n", port, err)
	}
}
