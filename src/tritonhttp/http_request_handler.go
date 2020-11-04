package tritonhttp

import (
	"net"
	"os"
	"strconv"
	"time"

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
	timeoutDuration := 5 * time.Second
	// bufReader := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	response := ""
	buf := make([]byte, 10)
	for {

		conn.SetReadDeadline(time.Now().Add(timeoutDuration))
		// defer conn.Close()

		// size, err := bufReader.Read(buf)
		size, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			// break
		}
		data := buf[:size]
		remaining = remaining + string(data)
		// bufReader.Reset(bufReader)
		// log.Println("original msg: " + remaining)
		for strings.Contains(remaining, "\r\n\r\n") {
			// conn.SetReadDeadline(time.Now().Add(timeoutDuration))
			// log.Println("delimiter: " + delimiter)
			// log.Println("original msg: " + remaining)
			log.Println("original:" + remaining)
			idx := strings.Index(remaining, "\r\n\r\n")
			msg := remaining[:idx] // whole requests
			remaining = remaining[idx+1:]
			reqSlice := strings.Split(msg, delimiter) // request
			log.Println(len(reqSlice))
			if len(reqSlice) < 2 {
				w.WriteString("400 Bad Request2")
				w.Flush()
				break
			}
			log.Println("message:" + msg)
			log.Println("left:" + remaining)
			initialLine := reqSlice[0] // get
			// log.Println("request: " + msg)
			// initial := remaining[idx+1:]

			// log.Println("left msg: " + remaining)
			// log.Println(initialLine[0])
			// w.WriteString(initialLine[0])
			// w.Flush()
			// log.Println("prefix host: ")
			// log.Println(strings.HasPrefix(msg, "Host"))
			// log.Println("prefix conn: ")
			// log.Println(strings.HasPrefix(msg, "Connection:"))
			// initialLine = GET /index.html HTTP/1.1
			if strings.HasPrefix(initialLine, "GET") {
				firstR := strings.Split(initialLine, " ")
				log.Println(firstR)
				// checkHTTP := firstR[0] == "HTTP/1.1"
				// checkLength := len(initialLine) == 3
				// checkGet := firstR[0] == "GET"
				if len(firstR) != 3 {
					// req.validInitial = false
					// w.WriteString("400 Bad Request1")
					// w.Flush()
					log.Println("error1")
					log.Println(firstR)
					log.Println(len(firstR))
					w.WriteString("400 Bad Request2")
					w.Flush()
					// hs.handleBadRequest(conn)
					break
				} else if firstR[2] != "HTTP/1.1" {
					log.Println("error2")
					log.Println(firstR[2])
					w.WriteString("400 Bad Request2")
					w.Flush()
					// hs.handleBadRequest(conn)
					break
				} else if firstR[0] != "GET" {
					log.Println("error3")
					log.Println(initialLine[0])
					w.WriteString("400 Bad Request2")
					w.Flush()
					// hs.handleBadRequest(conn)
					break
				} else if !strings.HasPrefix(firstR[1], "/") {
					w.WriteString("400 Bad Request2")
					w.Flush()
					log.Println("error4")
					// hs.handleBadRequest(conn)
					break
				} else {
					response += "HTTP/1.1 200 OK\r\n"
					response += "Server: Go-Triton-Server-1.0\r\n"
					// w.WriteString("HTTP/1.1 200 OK\r\n")
					// w.WriteString("Server: Go-Triton-Server-1.0\r\n")
					// w.Flush()

					url := hs.DocRoot + firstR[1]
					if firstR[1] == "/" {
						log.Println("url is 11/")
						url = hs.DocRoot + firstR[1] + "index.html"
					} else {
						log.Println("url is " + url)
						url = hs.DocRoot + firstR[1]
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
						w.WriteString("HTTP/1.1 404 Not Found\r\n")
						w.Flush()
						// hs.handleFileNotFoundRequest(conn)
						log.Fatal(err)
					}
					// get the size
					res.contentLength = fi.Size()
					// log.Println(res.contentType)
					// log.Println(size)

					res.lastModified = fi.ModTime().String()
					response += "Last-Modified: " + res.lastModified + "\r\n"
					response += "Content-Length: " + strconv.FormatInt(res.contentLength, 10) + "\r\n"
					response += "Contene-Type: " + res.contentType + "\r\n\r\n"
					// w.WriteString(response)
					// w.Flush()
					// w.WriteString("Last-Modified: " + rs.lastModified.String() + "\r\n")
					// w.Flush()
					// w.WriteString("Content-Length: " + strconv.FormatInt(size, 10) + "\r\n")
					// w.Flush()
					// w.WriteString("Contene-Type: " + res.contentType + "\r\n")
					// w.Flush()

				}
			} else {
				w.WriteString("400 Bad Request2")
				w.Flush()
				log.Println("error5")
				break
			}
			secondLine := reqSlice[1]
			log.Println("second l: " + secondLine)
			if strings.HasPrefix(secondLine, "Host") {
				// idxH := strings.Index(msg, ":")
				// msgH := msg[idxH+1:]

				// log.Println(msgH + " end of msgH")
				w.WriteString(response)
				w.Flush()
				// hs.sendResponse()
			} else {
				w.WriteString("400 Bad Request2")
				w.Flush()
				log.Println("error6")
				break
			}

			if len(reqSlice) > 2 {
				log.Println("handel connection")
				if strings.HasPrefix(reqSlice[2], "Connection:") {
					log.Println("handel connection2")
					idxH := strings.Index(reqSlice[2], ":")
					msgH := reqSlice[2][idxH+1:]
					connection := strings.TrimSpace(msgH)
					log.Println("conn msg:" + connection)
					if connection == "close" {
						res.connection = "close"
						w.WriteString("Connection: closed\r\n")
						w.Flush()
						conn.Close()
						log.Println("Connection closed by request.")
						return
					} else {
						log.Println("not close")
						res.connection = "no"
					}
				}
			}

		}
		log.Println("finish requests")
	}
}
