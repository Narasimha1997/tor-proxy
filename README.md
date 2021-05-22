## tor-proxy
An experimental standalone tor-proxy service built using Go, using go-proxy, go-libtor and bine. This is a simple replacement to Tor's original tor-proxy. (Since this is experimental, I would still recommend to use the original Tor-proxy)

### What this is for?
As you all know Tor is an anonymous network used to provide anonymous (tracking free) end-to-end network connectivity for users. In other words Tor network encrypts your IP with multiple layers such that the identity of the source cannot be determined anymore. Because of this capability of Tor, it is the de-facto network for hosting Dark-Web sites with almost 100% anonymity, clients use Tor browser to access these sites, usually those sites ending with `.onion`. Apart from this, Tor network can also provide anynomity while accessing normal HTTP/s sites. Using this prox, you can:

1. Provide a single point of contact for multiple devices in your to connect to External Tor Network (like an egress).
2. Access Dark-Web sites without having to install Tor browser on devices (use existing browsers, with proxy config)
3. Stay anonymous over the internet - use Tor's power across multiple devices.
4. Deploy this service as an ingress/egress on K8s to anynomize your microservices.
5. Access sites that are banned in your country (I don't recommend this and might now work always)

### Features:
1. Standalone binary - does not depend on `tor` application to be installed locally.
2. Dockerfiles provided to build and deploy the service in a container environment.
3. k8s deployment file provided to deploy tor-proxy on a cloud native environment. (For larger workloads).

### How to set-up and use:

#### 1. Install Locally (Not recommended)
Local installation requires a properly configured Go environment to be working and must be supporting Go modules. Just clone this repository:
```
```
Then run:
```
cd tor-proxy
export GOBIN=$PWD   # skip this if you don't want the binary to be installed in current location
go install
```

If everything goes well, you should see `tor-proxy` binary, you can run this simply by:
```
./tor-proxy
```

Changing port: When starting up, the code looks for `PROXY_PORT` env variable, by default it is set to `8000`, you can change the default port by running (for example):
```
export PROXY_PORT=8001
```

#### 2a. Running the binary in container environment:
First, follow step-1.
You can just place the pre-built go binary inside a container and execute it normally. `Dockerfile_binary` has the steps for this, just use this dockerfile to build the container image:

```
docker build . -f Dockerfile -t tor-proxy:latest
```

#### 2b. Build and run completely in a container environment:
This is the recommended approach, the default `Dockerfile` has steps for this process. It basically uses a pre-packaged Golang environment and builds the code and all it's dependencies inside the container. Just run:

```
docker build . -t tor-proxy:latest
```

#### 3. Run the container:
You can run this container just like any other normal container. (using docker for example)

```
docker run --rm -p 8000:8000 tor-proxy
```

If everything works as expected, you should be able to use `127.0.0.1:8000` as the tor proxy on local machine and `your.lan.ip.xx:8000` can be used across devices connected to a common network backbone. Or `some.public.ip.xx:8000` can be used worldwide, cosider deploying it on Heroku (example).

#### 4. Configure your browser to use the proxy:
This step is easy and varies across browsers. For firefox - [check-here](https://support.mozilla.org/en-US/kb/connection-settings-firefox)

### TODOs:
1. Support Socks555 connections.
2. Create a simple JavaScript sandbox that would use this proxy to avoid manual proxy configuration. (Should check how)

### Acknowledgements
1. [The Tor Project](https://www.torproject.org/)
2. [go-libtor](https://github.com/ipsn/go-libtor)
3. [bine](https://github.com/cretz/bine)
4. [goproxy](https://github.com/elazarl/goproxy)

### How to contribute?
There are no rules, just expecting contributions in the form of code, issues, PRs. (Not money, this is just a 2 hours weekend project and is completely focused on learning.)
