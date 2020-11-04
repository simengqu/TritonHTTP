package tritonhttp

import (
	"bufio"
	"net"
)

// import (
// 	"net"
// 	// "strings"
// 	// "fmt"
// 	"bufio"
// 	"log"
// 	"os"
// )

func (hs *HttpServer) handleBadRequest(conn net.Conn) {
	// panic("todo - handleBadRequest")
	w := bufio.NewWriter(conn)
	s := "HTTP/1.1 400 Bad Request\r\n"
	s += "Server: Go-Triton-Server-1.0\r\n"
	w.WriteString(s)
	w.Flush()
	conn.Close()

}

func (hs *HttpServer) handleFileNotFoundRequest(conn net.Conn) {
	// panic("todo - handleFileNotFoundRequest")
	w := bufio.NewWriter(conn)
	s := "HTTP/1.1 404 Not Found\r\n"
	s += "Server: Go-Triton-Server-1.0\r\n\r\n"
	w.WriteString(s)
	w.Flush()
	// conn.Close()
}

// func (hs *HttpServer) handleResponse(requestHeader *HttpRequestHeader, conn net.Conn) (result string) {
// 	// panic("todo - handleResponse")
// 	// server := "GoTriton-Server-1.0\r\n"
// 	// lastModified :=
// 	// contentType :=
// 	return server
// }

// func (hs *HttpServer) sendResponse(responseHeader HttpResponseHeader, conn net.Conn) {
// 	panic("todo - sendResponse")

// 	// Send headers
// 	headers := responseHeader.server
// 	// Send file if required

// 	// Hint - Use the bufio package to write response
// 	w := bufio.NewWriter(os.Stdout)

// 	log.Println(w, headers)
// 	w.Flush()
// }
