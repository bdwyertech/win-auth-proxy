package main

import (
    "fmt"
    "log"
    "github.com/elazarl/goproxy"
    "encoding/base64"
    "net/http"
    "os"
    "net/url"
    "strings"
)

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func hasNegotiateChallenge() goproxy.RespConditionFunc {
    return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
        return (resp.StatusCode == 401 && (contains(resp.Header["WWW-Authenticate"], "Negotiate") || contains(resp.Header["Www-Authenticate"], "Negotiate")))
    }
}

func hasProxyNegotiateChallenge() goproxy.RespConditionFunc {
    return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
        return (resp.StatusCode == 407 && contains(resp.Header["Proxy-Authenticate"], "Negotiate"))
    }
}

func authenticate(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
    ctx.Logf("Received 401 and Www-Authenticate from server, proceeding to reply")

    impl := CurrentOsGssImplementation{}

    spn := "http/" + strings.ToUpper(strings.Split(r.Request.Host,":")[0])

    ticket:= impl.GetTicket(ctx, spn)

    // Generate the Authorization header
    headerstr := "Negotiate " + base64.StdEncoding.EncodeToString(ticket)
    ctx.Logf("Generated header %s", headerstr)

    // Modify the original request, and rerun the request
    ctx.Req.Header["Authorization"] = []string{headerstr}
    client := http.Client{}
    newr, err := client.Do(ctx.Req)
    if err != nil {
        ctx.Warnf("New request failed: %v", err)
    }
    ctx.Logf("Got response, forwarding it back to client")

    // Return the new response in place of the original
    return newr
}

// when the client tries to establish an HTTPS connection
// INFO: Running 0 CONNECT handlers
// => memory leak
// the client gets a 502 (Bad Gateway)
// elazarl/goproxy/https.go: line 96: proxy.connectDial fails: the corporate proxy returns 407 with a Proxy-Authenticate header, but the OnResponse handler is not triggered!
// how to alter the connectDial request in order to inject the required Proxy-Authorization header?
func authenticateProxy(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
    ctx.Logf("Received 407 and Proxy-Authenticate from server, proceeding to reply")

    impl := CurrentOsGssImplementation{}

    // Forefront TMG sets the 'Via' header. This header contains the name of the host after a version number.
    // otherwise, the r.Header["Negotiate"] could contain a Basic realm="fqdn"  
    spn := strings.Split(r.Header["Via"][0], " ")[1]
    ticket:= impl.GetTicket(ctx, spn)

    // Generate the Authorization header
    headerstr := "Negotiate " + base64.StdEncoding.EncodeToString(ticket)
    ctx.Logf("Generated header %s", headerstr)

    // Modify the original request, and rerun the request
    ctx.Req.Header["Proxy-Authorization"] = []string{headerstr}
    client := getHTTPClient()
    newr, err := client.Do(ctx.Req)
    if err != nil {
        ctx.Warnf("New request failed: %v", err)
    }
    ctx.Logf("Got response, forwarding it back to client")

    // Return the new response in place of the original
    return newr
}

func getHTTPClient() http.Client {
    var client http.Client
    if (opts.proxy != "") {
        proxyURL, _ := url.Parse(opts.proxy)
        client = http.Client{ Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
    } else {
        client = http.Client{}
    }
    return client
}

var opts Arguments

func main() {
    opts = Parse(os.Args)
    if opts.autodetectProxy {   
        v, err := autodetectProxy("google.com")
        if err != nil {
            fmt.Printf("error raised during proxy autodetection: %v", err)
        } else {
            fmt.Printf("autodetected proxy: %v", v)     
        }
        os.Exit(0)
    }

    
    proxy := goproxy.NewProxyHttpServer()
    
    // behave as a MITM
    // the client fails with a TrustFailure (Le certificat distant n'est pas valide selon la procédure de validation.)
    // probably because the client expects an SSL Hello, but it receives a plain HTTP response
    /*proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(func (req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if req.URL.Scheme == "https" {
			req.URL.Scheme = "http"
		}
		return req, nil
	})*/
    
    proxy.Verbose = true

    if opts.proxy != "" {
        // transfer all the requests via the proxy specified on the command line as first positional argument
        proxy.Tr.Proxy = func (req *http.Request) (*url.URL, error) { return url.Parse(opts.proxy) }
    }

    proxy.OnResponse(hasNegotiateChallenge()).DoFunc(authenticate)
    proxy.OnResponse(hasProxyNegotiateChallenge()).DoFunc(authenticateProxy)
    
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", opts.listeningPort), proxy))
}
