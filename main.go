package main

import (
	"github.com/Narasimha1997/tor-proxy/core"
)

func main() {
	httpServing := core.HTTPServing{}
	httpServing.ListenAndServe()
}
