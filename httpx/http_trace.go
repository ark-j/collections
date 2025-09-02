package httpx

import (
	"crypto/tls"
	"log"
	"net/http/httptrace"
)

func getTracer() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(addr string) {
			log.Printf("connection to host and port = %s\n", addr)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			log.Printf("connection acquired: %+v\n", info)
		},
		PutIdleConn: func(err error) {
			if err != nil {
				log.Printf("put idle conn: %+v\n", err)
			}
		},
		GotFirstResponseByte: func() {
			log.Println("got first response byte")
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			log.Printf("dns started for host: %s\n", info.Host)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			log.Printf("dns resvolver done: %+v\n", info)
		},
		ConnectStart: func(network, addr string) {
			log.Printf("connection started at network: %s and addr: %s\n", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			log.Printf("connection done at network: %s and addr: %s\n", network, addr)
		},
		TLSHandshakeStart: func() {
			log.Println("tls handshake started")
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			if err != nil {
				log.Printf("tls handshake done: %+v\n", state)
			} else {
				log.Println("tls handshake done with err:", err)
			}
		},
		WroteHeaderField: func(key string, value []string) {
			log.Printf("header field written key: %s, value: %v\n", key, value)
		},
		WroteHeaders: func() {
			log.Println("writing of headers completed")
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			log.Printf("writing of request completed: %+v\n", info)
		},
	}
}
