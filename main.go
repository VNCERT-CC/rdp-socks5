package main

import (
	"flag"
	"log"
	"net"

	"h12.io/socks"
)

var lAddr = flag.String(`l`, `127.0.0.1:3388`, `Listen address`)
var socksURI = flag.String(`x`, `socks5://127.0.0.1:1080?timeout=5m`, `Socks URI`)
var proxyAddr = flag.String(`r`, `10.10.10.10:3389`, `Remote address`)
var dialFunc func(string, string) (net.Conn, error)

func copy2(dst net.Conn, src net.Conn) {
	defer dst.Close()
	defer src.Close()
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if err != nil {
			// log.Println(`Read err:`, err)
			return
		}
		// log.Println(src.RemoteAddr().String(), `=>`, dst.RemoteAddr().String(), `:`, len(buf[:n]))
		_, err = dst.Write(buf[:n])
		if err != nil {
			// log.Println(`Write errL`, err)
			return
		}
	}
}

func handleConn(lconn net.Conn) {
	defer lconn.Close()
	rconn, err := dialFunc(`tcp`, *proxyAddr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(`Connected to`, *proxyAddr)
	defer rconn.Close()
	go copy2(lconn, rconn)
	copy2(rconn, lconn)
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
