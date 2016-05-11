package main

import "flag"

type Arguments struct {
    proxy string
    autodetectProxy bool
    listeningPort uint
}

func Parse(commandLine []string) Arguments {
    var args = Arguments{}
    flag.StringVar(&args.proxy, "proxy", "", "corporate proxy address")
    flag.StringVar(&args.proxy, "x", "", "corporate proxy address")
    flag.BoolVar(&args.autodetectProxy, "autodetect", false, "try to autodetect the corporate proxy")
    flag.BoolVar(&args.autodetectProxy, "a", false, "try to autodetect the corporate proxy")
    flag.UintVar(&args.listeningPort, "port", 8080, "listening port (default: 8080)")
    flag.UintVar(&args.listeningPort, "p", 8080, "listening port (default: 8080)")
    flag.Parse()
    return args
}