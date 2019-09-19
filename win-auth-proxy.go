package main

import (
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func main() {

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	var auth_data string
	var AlwaysMitmAuth goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		// Initialise authentication data
		auth_data = getAuthorizationHeader(os.Args[1])

		// Transfer all the requests via the proxy specified on the command line as first positional argument
		proxy.Tr = &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) { return url.Parse(os.Args[1]) },
			ProxyConnectHeader: http.Header{
				"Proxy-Authorization": {auth_data},
			},
		}
		return goproxy.MitmConnect, host
	}

	// Handle HTTP authenticate responses
	proxy.OnResponse(HasNegotiateChallenge()).DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		ctx.Logf("Received 407 and Proxy-Authenticate from server, proceeding to reply")

		headerstr := getAuthorizationHeader(os.Args[1])

		// Modify the original request, and rerun the request
		ctx.Req.Header["Proxy-Authorization"] = []string{headerstr}
		client := http.Client{
			Transport: proxy.Tr,
		}

		newr, err := client.Do(ctx.Req)

		if err != nil {
			ctx.Warnf("New request failed: %v", err)
		}

		ctx.Logf("Got response, forwarding it back to client")

		// Return the new response in place of the original
		return newr
	})

	// Handle HTTPS Connect Requests
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*"))).HandleConnect(AlwaysMitmAuth)

	// Handle HTTP Connect Requests
	proxy.OnRequest(goproxy.Not(goproxy.ReqHostMatches(regexp.MustCompile("^.*:443$")))).
		DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

			// Initialise authentication data
			auth_data = getAuthorizationHeader(os.Args[1])

			// Transfer all the requests via the proxy specified on the command line as first positional argument
			proxy.Tr = &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) { return url.Parse(os.Args[1]) },
				ProxyConnectHeader: http.Header{
					"Proxy-Authorization": {auth_data},
				},
			}

			return req, nil
		})

	log.Fatal(http.ListenAndServe(":53128", proxy))

}
