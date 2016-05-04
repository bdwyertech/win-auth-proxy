package main

import "flag"

type Arguments struct {
    proxy string
    autodetectProxy bool
}

func Parse(commandLine []string) Arguments {
    var args = Arguments{}
    flag.StringVar(&args.proxy, "proxy", "", "corporate proxy address")
    flag.StringVar(&args.proxy, "x", "", "corporate proxy address")
    flag.BoolVar(&args.autodetectProxy, "autodetect", false, "try to autodetect the corporate proxy")
    flag.BoolVar(&args.autodetectProxy, "a", false, "try to autodetect the corporate proxy")
    flag.Parse()
    return args
}