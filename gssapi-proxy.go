package main

import (
"log"
"github.com/elazarl/goproxy"
"encoding/base64"
"net/http"
"os"
"net/url"
)

func HasNegotiateChallenge() goproxy.RespConditionFunc {
    return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
        return ((resp.StatusCode == 401) && (len(resp.Header["Www-Authenticate"])>0))
    }
}

func main() {
    var args = parse()
    proxy := goproxy.NewProxyHttpServer()
    proxy.Verbose = true

    if args.autodetectProxy {
        args.proxy, _ = autodetectProxy("google.com")
    }
    
    // transfer all the requests via the proxy specified on the command line as first positional argument
    proxy.Tr = &http.Transport { Proxy: func (req *http.Request) (*url.URL, error) { return url.Parse(os.Args[1]) } }

    proxy.OnResponse(HasNegotiateChallenge()).DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx)*http.Response {
        ctx.Logf("Received 401 and Www-Authenticate from server, proceeding to reply")

        impl := CurrentOsGssImplementation{}
        
        ticket:= impl.GetTicket(ctx)
                
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
    })
    log.Fatal(http.ListenAndServe(":8080", proxy))
}
