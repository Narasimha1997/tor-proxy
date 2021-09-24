package main

import (
	"os"

	"github.com/Narasimha1997/tor-proxy/core"
)

func getProxyProtocol() string {
	protocol, ok := os.LookupEnv("PROXY_PROTOCOL")
	if !ok {
		protocol = "socks5"
	}

	if protocol != "socks5" && protocol != "http" {
		protocol = "socks5"
	}

	return protocol
}

func main() {
	proxyProtocol := getProxyProtocol()
	if proxyProtocol == "socks5" {
		socks5Server := core.Socks5Serving{}
		socks5Server.ListenAndServe()
	} else {
		httpServer := core.HTTPServing{}
		httpServer.ListenAndServe()
	}
}
