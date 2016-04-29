package main

import (
    "log"
    "github.com/jackwakefield/gopac"
)

func autodetectProxy(host string) (string, error) {
    parser := new(gopac.Parser)

    // use parser.Parse(path) to parse a local file
    // or parser.ParseUrl(url) to parse a remote file
    // https://technet.microsoft.com/fr-fr/library/cc995261.aspx
    // https://en.wikipedia.org/wiki/Web_Proxy_Autodiscovery_Protocol
    if err := parser.ParseUrl("http://wpad/wpad.dat"); err != nil {
        log.Fatalf("Failed to parse PAC (%s)", err)
        return "", err
    }

    // find the proxy entry for host
    entry, err := parser.FindProxy("", host)

    if err != nil {
        log.Fatalf("Failed to find a proxy entry (%s)", err)
        return "", err
    }

    log.Println(entry)
    return entry, err
}