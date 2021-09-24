package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/armon/go-socks5"
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
func initTorHandle(initHTTPTransport bool) *TorHandle {

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

	if initHTTPTransport {
		torHandle.torTransport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
			DialContext:     torHandle.torDialer.DialContext,
		}
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

// HTTPServing Type definition for HTTPServing
type HTTPServing struct{}

func (httpProxy *HTTPServing) getProxy() (*goproxy.ProxyHttpServer, *TorHandle) {
	proxy := goproxy.NewProxyHttpServer()
	tcx := initTorHandle(true)

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

// Socks5Serving Type definition for Socks5 server
type Socks5Serving struct{}
type ProxyAuthParams struct {
	username string
	password string
}

func (sserver *Socks5Serving) getAuthParams() *ProxyAuthParams {

	usernameEnv, ok := os.LookupEnv("PROXY_USERNAME")
	if !ok {
		usernameEnv = ""
	}

	passwordEnv, ok := os.LookupEnv("PROXY_PASSWORD")
	if !ok {
		passwordEnv = ""
		usernameEnv = ""
	}

	return &ProxyAuthParams{
		username: usernameEnv,
		password: passwordEnv,
	}
}

func (sserver *Socks5Serving) setTorDialer(config *socks5.Config) *TorHandle {
	tcx := initTorHandle(false)
	config.Dial = tcx.torDialer.DialContext

	return tcx
}

func (sserver *Socks5Serving) ListenAndServe() {
	socks5Config := socks5.Config{
		Logger: log.New(os.Stdin, "", log.LstdFlags),
	}

	// get auth parameters from env:
	authParams := sserver.getAuthParams()
	if authParams.username != "" {
		socks5Config.AuthMethods = []socks5.Authenticator{
			socks5.UserPassAuthenticator{Credentials: socks5.StaticCredentials{
				authParams.username: authParams.password,
			}},
		}
	}

	// set TorHandle
	torHandle := sserver.setTorDialer(&socks5Config)
	defer clearHandle(torHandle)

	// create the server object
	server, err := socks5.New(&socks5Config)

	if err != nil {
		log.Fatalln(err)
	}

	// start the server on default Port or custom provided port
	port, ok := os.LookupEnv("PROXY_PORT")
	if !ok {
		port = DefaultProxyPort
	}

	err = server.ListenAndServe("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln(err)
	}
}
