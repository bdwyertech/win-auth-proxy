package main

import (
    "fmt"
    "log"
    "github.com/nilleb/goproxy"
    "encoding/base64"
    "net/http"
    "os"
    "net/url"
    "strings"
    "bufio"
    "net"
    "errors"
    "io/ioutil"
    "io"
    "crypto/tls"
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
        return (resp != nil && resp.StatusCode == 401 && (contains(resp.Header["WWW-Authenticate"], "Negotiate") || contains(resp.Header["Www-Authenticate"], "Negotiate")))
    }
}

func hasProxyNegotiateChallenge() goproxy.RespConditionFunc {
    return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
        return (resp != nil && resp.StatusCode == 407 && contains(resp.Header["Proxy-Authenticate"], "Negotiate"))
    }
}

func authenticate(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
    ctx.Logf("Received 401 and Www-Authenticate from server, proceeding to reply")

    impl := CurrentOsGssImplementation{}

    spn := "http/" + strings.ToUpper(strings.Split(r.Request.Host,":")[0])

    ticket:= impl.GetTicket(spn)

    // Generate the Authorization header
    headerstr := "Negotiate " + base64.StdEncoding.EncodeToString(ticket)

    // Modify the original request, and rerun the request
    ctx.Req.Header["Authorization"] = []string{headerstr}
    client := http.Client{}
    newr, err := client.Do(ctx.Req)
    if err != nil {
        ctx.Warnf("New request failed: %v", err)
        return r
    }
    ctx.Logf("Got response, forwarding it back to client")

    // Return the new response in place of the original
    return newr
}

// when the client tries to establish an HTTPS connection
// INFO: Running 0 CONNECT handlers
// => memory leak
// the client gets a 502 (Bad Gateway) returned by goproxy httpError
// elazarl/goproxy/https.go: line 96: proxy.connectDial fails: the corporate proxy returns 407 with a Proxy-Authenticate header, but the OnResponse handler is not triggered!
// how to alter the connectDial request in order to inject the required Proxy-Authorization header?
func authenticateProxy(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
    ctx.Logf("Received 407 and Proxy-Authenticate from server, proceeding to reply")

    impl := CurrentOsGssImplementation{}

    // Forefront TMG sets the 'Via' header. This header contains the name of the host after a version number.
    // otherwise, the r.Header["Negotiate"] could contain a Basic realm="fqdn"  
    spn := strings.Split(r.Header["Via"][0], " ")[1]
    ticket:= impl.GetTicket(spn)

    // Generate the Authorization header
    headerstr := "Negotiate " + base64.StdEncoding.EncodeToString(ticket)

    // Modify the original request, and rerun the request
    ctx.Req.Header["Proxy-Authorization"] = []string{headerstr}
    client := getHTTPClient()
    newr, err := client.Do(ctx.Req)
    if err != nil {
        ctx.Warnf("New request failed: %v", err)
        return r
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


func traceRequests(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
    ctx.Logf("requesting %s, %s", r.Host, r.URL)
    return r, nil
}

/*
func authenticatedConnectRequest(addr string, resp *http.Response) *http.Request {
    impl := CurrentOsGssImplementation{}

    spn := strings.Split(resp.Header["Via"][0], " ")[1]
    ticket:= impl.GetTicket(spn)

    // Generate the Authorization header
    headerstr := "Negotiate " + base64.StdEncoding.EncodeToString(ticket)

    connectReq := &http.Request{
        Method: "CONNECT",
        URL:    &url.URL{Opaque: addr},
        Host:   addr,
        Header: make(http.Header),
    }
    // Modify the original request, and rerun the request
    connectReq.Header["Proxy-Authorization"] = []string{headerstr}
    return connectReq
}

func dial(proxy *goproxy.ProxyHttpServer, network, addr string) (c net.Conn, err error) {
	if proxy.Tr.Dial != nil {
		return proxy.Tr.Dial(network, addr)
	}
	return net.Dial(network, addr)
}

func myNewConnectDialToProxy(proxy *goproxy.ProxyHttpServer, https_proxy string) func(network, addr string) (net.Conn, error) {
	u, err := url.Parse(https_proxy)
	if err != nil {
		return nil
	}
	if u.Scheme == "" || u.Scheme == "http" {
		if strings.IndexRune(u.Host, ':') == -1 {
			u.Host += ":80"
		}
		return func(network, addr string) (net.Conn, error) {
			connectReq := &http.Request{
				Method: "CONNECT",
				URL:    &url.URL{Opaque: addr},
				Host:   addr,
				Header: make(http.Header),
			}
			c, err := dial(proxy, network, u.Host)
			if err != nil {
				return nil, err
			}
			connectReq.Write(c)
			// Read response.
			// Okay to use and discard buffered reader here, because
			// TLS server will not speak until spoken to.
			br := bufio.NewReader(c)
			resp, err := http.ReadResponse(br, connectReq)
			if err != nil {
				c.Close()
				return nil, err
			}
            if resp.StatusCode == 407 {
                c.Close()
                c, err := dial(proxy, network, u.Host)
                if err != nil {
                    return nil, err
                }
                connectReq = authenticatedConnectRequest(addr, resp)
                connectReq.Write(c)
                // Read response.
                // Okay to use and discard buffered reader here, because
                // TLS server will not speak until spoken to.
                br := bufio.NewReader(c)
                resp, err = http.ReadResponse(br, connectReq)
                if err != nil {
                    c.Close()
                    return nil, err
                }
            }
			if resp.StatusCode != 200 {
				resp, _ := ioutil.ReadAll(resp.Body)
				c.Close()
				return nil, errors.New("proxy refused connection " + string(resp))
			}
            // problem: the proxy doesn't return the remote answer. Just a 200, connection established.
			return c, nil
		}
	}
	if u.Scheme == "https" {
		if strings.IndexRune(u.Host, ':') == -1 {
			u.Host += ":443"
		}
		return func(network, addr string) (net.Conn, error) {
			c, err := dial(proxy, network, u.Host)
			if err != nil {
				return nil, err
			}
			c = tls.Client(c, proxy.Tr.TLSClientConfig)
			connectReq := &http.Request{
				Method: "CONNECT",
				URL:    &url.URL{Opaque: addr},
				Host:   addr,
				Header: make(http.Header),
			}
			connectReq.Write(c)
			// Read response.
			// Okay to use and discard buffered reader here, because
			// TLS server will not speak until spoken to.
			br := bufio.NewReader(c)
			resp, err := http.ReadResponse(br, connectReq)
			if err != nil {
				c.Close()
				return nil, err
			}
			if resp.StatusCode != 200 {
				body, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 500))
				resp.Body.Close()
				c.Close()
				return nil, errors.New("proxy refused connection" + string(body))
			}
			return c, nil
		}
	}
	return nil
}
*/

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
    
    proxy.OnRequest().DoFunc(traceRequests)
    
    proxy.Verbose = true

    if opts.proxy != "" {
        // transfer all the requests via the proxy specified on the command line as first positional argument
        proxy.Tr.Proxy = func (req *http.Request) (*url.URL, error) { return url.Parse(opts.proxy) }
        //proxy.ConnectDial = myNewConnectDialToProxy(proxy, opts.proxy)
    }

    proxy.OnResponse(hasNegotiateChallenge()).DoFunc(authenticate)
    proxy.OnResponse(hasProxyNegotiateChallenge()).DoFunc(authenticateProxy)
    
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", opts.listeningPort), proxy))
}
