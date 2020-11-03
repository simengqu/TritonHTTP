package tritonhttp

import (
	"net"
	// "fmt"
	"strings"
	"log"
	// "time"
	"bufio"
)

/* 
For a connection, keep handling requests until 
	1. a timeout occurs or
	2. client closes connection or
	3. client sends a bad request
*/
func (hs *HttpServer) handleConnection(conn net.Conn) {

	// panic("todo - handleConnection")

	// Start a loop for reading requests continuously
	// Set a timeout for read operation
	// Read from the connection socket into a buffer
	// Validate the request lines that were read
	// Handle any complete requests
	// Update any ongoing requests
	// If reusing read buffer, truncate it before next read

	delimiter := "\n"
	remaining := ""
	fl := true
	// timeoutDuration := 5 * time.Second
	bufReader := bufio.NewReader(conn)
	for {
		w := bufio.NewWriter(conn)
		
		// conn.SetReadDeadline(time.Now().Add(timeoutDuration))
		buf := make([]byte, 32)
		size, err := bufReader.Read(buf)
		if err != nil {
			log.Println("err")
			break
		}
		data := buf[:size]	
		remaining = remaining + string(data)

		for strings.Contains(remaining, delimiter) {
			idx := strings.Index(remaining, delimiter)
			msg := remaining[:idx]
			// initial := remaining[idx+1:]
			initialLine := strings.Fields(msg)
			remaining = remaining[idx+1:]
			log.Println(msg)
			// log.Println(initialLine[0])
			// w.WriteString(initialLine[0])
			// w.Flush()
			if strings.HasPrefix(msg, "GET"){
				if len(initialLine) != 3 || initialLine[len(initialLine)-1] != "HTTP/1.1" || initialLine[0] != "GET" && fl{
					// req.validInitial = false
					// w.WriteString("400 Bad Request1")
					// w.Flush()
					log.Println("error1")
					log.Println(initialLine)
					hs.handleBadRequest(conn)
					break
				}else if !strings.HasPrefix(initialLine[1], "/") {
					// w.WriteString("400 Bad Request2")
					// w.Flush()
					log.Println("error2")
					hs.handleBadRequest(conn)
					break
				}else {
					log.Println("HTTP/1.1 200 OK\r\n")
					w.WriteString("HTTP/1.1 200 OK\r\n")
					w.Flush()
				}
			} else if strings.HasPrefix(msg, "Host:") {
				fl = false
				idxH := strings.Index(msg, ":")
				msgH := msg[idxH+1:]
				log.Println(msgH)
				w.WriteString("HOST\r\n")
				w.Flush()
				// req.host = strings.TrimSpace(msgH)
				// req.validInitial = true

			} else if strings.HasPrefix(msg, "Connection:") {
				fl = false
				idxH := strings.Index(msg, ":")
				msgH := msg[idxH+1:]
				connection := strings.TrimSpace(msgH)
				if connection == "Close" {
					conn.Close()
					log.Println("Connection closed by request.")
					break
				}
				w.WriteString("CONNECTION\r\n")
				w.Flush()
				// req.validInitial = true
			}
		}
		// w.WriteString("reading\r\n")
		// w.Flush()
	}
}
