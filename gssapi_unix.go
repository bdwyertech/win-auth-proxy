// +build !windows

package main

import "github.com/elazarl/goproxy"

// the idea here is to use something like https://github.com/apcera/gssapi/blob/master/test/client_access_test.go
// to get a kerberos ticket (given that the current session already has a TGT)
type CurrentOsGssImplementation struct {
}

func (t CurrentOsGssImplementation) GetTicket(ctx *goproxy.ProxyCtx) []byte {
    return []byte{}
}