package main

import "github.com/elazarl/goproxy"

type GssApiImplementation interface {
    GetTicket(ctx *goproxy.ProxyCtx, r *http.Response) []byte
}