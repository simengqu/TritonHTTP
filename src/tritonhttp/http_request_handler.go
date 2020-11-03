package tritonhttp

import (
	"net"
	"os"
	"strconv"

	// "fmt"
	"log"
	"strings"

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

	delimiter := "\r\n"
	remaining := ""
	var res HttpResponseHeader
	// timeoutDuration := 5 * time.Second
	bufReader := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	for {

		// conn.SetReadDeadline(time.Now().Add(timeoutDuration))
		buf := make([]byte, 1024)
		size, err := bufReader.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}
		// defer conn.Close()
		data := buf[:size]
		remaining = remaining + string(data)
		// bufReader.Reset(bufReader)

		for strings.Contains(remaining, delimiter) {
			idx := strings.Index(remaining, delimiter)
			msg := remaining[:idx]
			// initial := remaining[idx+1:]
			initialLine := strings.Fields(msg)
			remaining = remaining[idx+1:]
			// log.Println(msg)
			// log.Println(initialLine[0])
			// w.WriteString(initialLine[0])
			// w.Flush()
			if strings.HasPrefix(msg, "GET") {
				log.Println(initialLine)
				checkHTTP := strings.TrimSpace(initialLine[len(initialLine)-1]) == "HTTP/1.1"
				checkLength := len(initialLine) == 3
				checkGet := initialLine[0] == "GET"
				if !checkLength {
					// req.validInitial = false
					// w.WriteString("400 Bad Request1")
					// w.Flush()
					log.Println("error1")
					log.Println(len(initialLine))
					hs.handleBadRequest(conn)
					break
				} else if !checkHTTP {
					log.Println("error2")
					log.Println(!checkHTTP)
					log.Println(initialLine[len(initialLine)-1])
					hs.handleBadRequest(conn)
					break
				} else if !checkGet {
					log.Println("error3")
					log.Println(initialLine[0])
					hs.handleBadRequest(conn)
					break
				} else if !strings.HasPrefix(initialLine[1], "/") {
					// w.WriteString("400 Bad Request2")
					// w.Flush()
					log.Println("error4")
					hs.handleBadRequest(conn)
					break
				} else {
					w.WriteString("HTTP/1.1 200 OK\r\n")
					w.WriteString("Server: Go-Triton-Server-1.0\r\n")
					w.Flush()

					url := hs.DocRoot + initialLine[1]
					if initialLine[1] == "/" {
						log.Println("url is 11/")
						url = hs.DocRoot + initialLine[1] + "index.html"
					} else {
						log.Println("url is " + url)
						url = hs.DocRoot + initialLine[1]
					}
					log.Println(url)

					lastIdx := strings.LastIndex(url, ".")
					extension := url[lastIdx:]

					if ext, ok := hs.MIMEMap[extension]; ok {
						res.contentType = ext
						log.Println("extension found: " + extension + " | " + ext)
					} else {
						res.contentType = "application/octet-stream"
						log.Println("extension not found: " + extension)
					}

					// lastIdx := strings.LastIndex(url, "/")
					// res.contentType = initialLine[1][lastIdx:]
					fi, err := os.Stat(url)
					if err != nil {
						log.Fatal(err)
					}
					// get the size
					size := fi.Size()
					log.Println(res.contentType)
					log.Println(size)

					modifiedtime := fi.ModTime()
					w.WriteString("Last-Modified: " + modifiedtime.String() + "\r\n")
					w.Flush()
					w.WriteString("Content-Length: " + strconv.FormatInt(size, 10) + "\r\n")
					w.Flush()
					w.WriteString("Contene-Type: " + res.contentType + "\r\n")
					w.Flush()

				}
			}
			if strings.HasPrefix(msg, "Host:") {
				idxH := strings.Index(msg, ":")
				msgH := msg[idxH+1:]
				log.Println(msgH + "end of msgH")
			}
			if strings.HasPrefix(msg, "Connection:") {
				idxH := strings.Index(msg, ":")
				msgH := msg[idxH+1:]
				connection := strings.TrimSpace(msgH)
				if connection == "close" {
					w.WriteString("Connection: closed\r\n")
					w.Flush()
					conn.Close()
					log.Println("Connection closed by request.")
					break
				}
			}
		}
	}
}
