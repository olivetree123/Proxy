package proxy

import (
	"fmt"
	"time"

	// "crypto/tls"
	"io"
	"log"
	"net"
)

// Proxy - Manages a Proxy connection, piping data between local and remote.
type Proxy struct {
	sentBytes     uint64
	receivedBytes uint64
	laddr, raddr  *net.TCPAddr
	//lconn, rconn  io.ReadWriteCloser
	lconn, rconn *net.TCPConn
	erred        bool
	errsig       chan bool
	tlsAddress   string
}

// New - Create a new Proxy instance. Takes over local connection passed in,
// and closes it when finished.
func New(lconn *net.TCPConn, laddr, raddr *net.TCPAddr) *Proxy {
	return &Proxy{
		lconn:  lconn,
		laddr:  laddr,
		raddr:  raddr,
		erred:  false,
		errsig: make(chan bool),
	}
}

// Start - open connection to remote and start proxying data.
func (p *Proxy) Start() {
	defer p.lconn.Close()

	var err error
	//connect to remote
	p.rconn, err = net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		log.Fatalln("Remote connection failed: ", err)
		return
	}
	defer p.rconn.Close()
	//display both ends
	log.Println("Opened ", p.laddr.String(), " >>> ", p.raddr.String())

	//bidirectional copy
	go func() {
		//n, err := io.Copy(p.lconn, p.rconn)
		p.pipe(p.lconn, p.rconn)
	}()
	go func() {
		//n, err := io.Copy(p.rconn, p.lconn)
		p.pipe(p.rconn, p.lconn)
	}()

	//wait for close...
	<-p.errsig
	//<-p.errsig
	log.Println(fmt.Sprintf("Closed (%d bytes sent, %d bytes recieved)", p.sentBytes, p.receivedBytes))
}

func (p *Proxy) err(s string, err error) {
	if p.erred {
		return
	}
	if err != io.EOF {
		log.Println(s, err)
	}
	p.erred = true
	p.errsig <- true
}

func (p *Proxy) pipe(src, dst *net.TCPConn) {
	isLocal := src == p.lconn

	t, _ := time.ParseDuration("10s")
	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)
	for {
		_ = src.SetReadDeadline(time.Now().Add(t))
		n, err := src.Read(buff)
		if err != nil {
			p.err("Read failed\n", err)
			return
		}
		b := buff[:n]

		//write out result
		_ = dst.SetWriteDeadline(time.Now().Add(t))
		n, err = dst.Write(b)
		if err != nil {
			p.err("Write failed\n", err)
			return
		}
		if isLocal {
			p.sentBytes += uint64(n)
		} else {
			p.receivedBytes += uint64(n)
		}
	}
}
