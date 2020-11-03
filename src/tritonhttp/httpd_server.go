package tritonhttp

import (
	// "fmt"
	"log"
	"net"
)
/** 
	Initialize the tritonhttp server by populating HttpServer structure
**/
func NewHttpdServer(port, docRoot, mimePath string) (*HttpServer, error) {
	// panic("todo - NewHttpdServer")

	// Initialize mimeMap for server to refer
	var hs HttpServer
	var err error
	hs.ServerPort = port
	hs.MIMEPath = mimePath
	hs.DocRoot = docRoot
	hs.MIMEMap, err = ParseMIME(hs.MIMEPath)

	// Return pointer to HttpServer
	return &hs, err
}

/** 
	Start the tritonhttp server
**/
func (hs *HttpServer) Start() (err error) {
	// panic("todo - StartServer")

	// Start listening to the server port

	// Accept connection from client

	// Spawn a go routine to handle request
	
	log.Println("Handling new connection...")

	ln, err := net.Listen("tcp", hs.ServerPort)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		ln.Close()
		log.Println("Listener closed")
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		go hs.handleConnection(conn)
	}
	return err
}

