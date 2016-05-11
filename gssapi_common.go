package main

import "github.com/elazarl/goproxy"

type GssApiImplementation interface {
    GetTicket(ctx *goproxy.ProxyCtx, host string) []byte
}