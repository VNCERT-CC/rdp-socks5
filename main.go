package main

import (
	"flag"
	"log"
	"net"
	"io"
	"time"

	"h12.io/socks"
)

var lAddr = flag.String(`l`, `127.0.0.1:3388`, `Listen address`)
var socksURI = flag.String(`x`, `socks5://127.0.0.1:1080?timeout=15m`, `Socks URI`)
var proxyAddr = flag.String(`r`, `10.10.10.10:3389`, `Remote address`)
var dialFunc func(string, string) (net.Conn, error)

func handleConn(lconn net.Conn) {
	rconn, err := dialFunc(`tcp`, *proxyAddr)
	if err != nil {
		lconn.Close()
		log.Println(err)
		return
	}
	log.Println(`Connected to`, *proxyAddr)
	defer func () {
		time.Sleep(time.Second)
		rconn.Close()
		lconn.Close()
	}()
	go io.Copy(lconn, rconn)
	io.Copy(rconn, lconn)
	log.Println(`Closed:`, lconn.RemoteAddr().String())
}

func main() {
	flag.Parse()
	dialFunc = socks.Dial(*socksURI)
	srv, err := net.Listen(`tcp`, *lAddr)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(`Proxing`, *lAddr, `to`, *proxyAddr, `via`, *socksURI)
	for {
		conn, err := srv.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(`Got connection from`, conn.RemoteAddr().String())
		go handleConn(conn)
	}
}
