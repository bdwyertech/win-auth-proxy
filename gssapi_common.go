package main

type GssApiImplementation interface {
    GetTicket(host string) []byte
}