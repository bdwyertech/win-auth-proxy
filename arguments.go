package main

import "flag"

type arguments struct {
    proxy string
    autodetectProxy bool
}

func parse() arguments {
    var args = arguments{}
    flag.StringVar(&args.proxy, "proxy", "", "corporate proxy address")
    flag.StringVar(&args.proxy, "x", "", "corporate proxy address")
    flag.BoolVar(&args.autodetectProxy, "autodetect", false, "try to autodetect the corporate proxy")
    flag.BoolVar(&args.autodetectProxy, "a", false, "try to autodetect the corporate proxy")
    flag.Parse()
    return args
}