package tritonhttp

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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
	timeoutDuration := time.Now().Add(5 * time.Second)
	bufReader := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	response := ""
	code := 0
	fmt.Println("\n\n======================================\nIn Go routine\n======================================")
	for time.Now().Before(timeoutDuration) {
		buf := make([]byte, 1024)
		defer conn.Close()
		conn.SetReadDeadline(timeoutDuration)

		size, err := bufReader.Read(buf)
		// size, err := conn.Read(buf)
		if err != nil {
			// fmt.Println(err)
		}

		data := buf[:size]
		remaining = remaining + string(data)
		url := ""
		if size != 0 {

			for strings.Contains(remaining, "\r\n\r\n") {
				fmt.Println("IN CONTAINS:::::::::::::::::::")
				fmt.Println("original: " + remaining)
				idx := strings.Index(remaining, "\r\n\r\n")
				msg := remaining[:idx] // whole requests
				remaining = remaining[idx+4:]
				if strings.HasPrefix(remaining, "GET") {
					fmt.Println("Has GET-----")
				}
				fmt.Println("remaining: " + remaining)
				reqSlice := strings.Split(msg, delimiter) // request
				fmt.Println("processing: ")
				fmt.Println(reqSlice)
				if len(reqSlice) < 2 {
					hs.handleBadRequest(conn)
					break
				}
				fmt.Println("message:" + msg)
				fmt.Println("left:" + remaining)
				initialLine := reqSlice[0] // get

				// initialLine = GET /index.html HTTP/1.1
				if strings.HasPrefix(initialLine, "GET") {
					firstR := strings.Split(initialLine, " ")
					fmt.Println(firstR)
					if len(firstR) != 3 {
						fmt.Println("error1")
						fmt.Println(firstR)
						fmt.Println(len(firstR))
						code = 400
						hs.handleBadRequest(conn)
						break
					} else if firstR[2] != "HTTP/1.1" {
						fmt.Println("error2")
						fmt.Println(firstR[2])
						code = 400
						hs.handleBadRequest(conn)
						break
					} else if firstR[0] != "GET" {
						fmt.Println("error3")
						fmt.Println(initialLine[0])
						code = 400
						hs.handleBadRequest(conn)
						break
					} else if !strings.HasPrefix(firstR[1], "/") {
						fmt.Println("error4")
						code = 400
						hs.handleBadRequest(conn)
						break
					} else {
						// check if file is valid
						if !strings.HasPrefix(firstR[1], hs.DocRoot) {
							url = hs.DocRoot + firstR[1]
						}

						idxFirstR := strings.LastIndex(firstR[1], "/")
						if firstR[1] == "/" {
							fmt.Println("url is 11/:")
							url = hs.DocRoot + firstR[1] + "index.html"
						} else if idxFirstR == len(firstR[1])-1 {
							fmt.Println("url is 22/:")
							url = hs.DocRoot + firstR[1] + "index.html"
						} else {
							url = hs.DocRoot + firstR[1]
						}
						fi, err := os.Open(url)
						defer fi.Close()
						if err != nil {
							code = 404
							hs.handleFileNotFoundRequest(conn)
							fmt.Println(err)
							break
						}

						lastIdx := strings.LastIndex(url, ".")
						extension := url[lastIdx:]

						if ext, ok := hs.MIMEMap[extension]; ok {
							res.contentType = ext
							fmt.Println("extension found: " + extension + " | " + ext)
						} else {
							res.contentType = "application/octet-stream"
							fmt.Println("extension not found: " + extension)
						}
					}
				}

				// check if headers valid
				for i := 2; i < len(reqSlice); i++ {
					kv := strings.Split(reqSlice[i], ":")
					if len(kv) != 2 {
						code = 400
						hs.handleBadRequest(conn)
						break
					}
				}
				// check if host valid
				secondLine := reqSlice[1]
				fmt.Println("second l: " + secondLine)
				if strings.HasPrefix(secondLine, "Host:") {
					fmt.Println("Has Host")
					code = 200

					fi, err := os.Open(url)
					defer fi.Close()
					if err != nil {
						code = 404
						hs.handleFileNotFoundRequest(conn)
						fmt.Println(err)
						break
					}
					fiStat, _ := fi.Stat()
					res.contentLength = fiStat.Size()
					response += "HTTP/1.1 200 OK\r\n"
					response += "Server: Go-Triton-Server-1.0\r\n"
					res.lastModified = fiStat.ModTime().String()
					response += "Last-Modified: " + res.lastModified + "\r\n"
					response += "Content-Length: " + strconv.FormatInt(res.contentLength, 10) + "\r\n"
					response += "Content-Type: " + res.contentType + "\r\n\r\n"
					fmt.Println("\n-=-=-=-=-=-==-=-=-=-=-=-=-=-=-=-=-=--=-\nWriting 200 OK response\n0-0-==0-=0-=0-=0-=0=0-=0-=0=")
					w.WriteString(response)
					w.Flush()
					// io.Copy(w, fi)
					// w.Flush()
				} else {
					fmt.Println("error6")
					hs.handleBadRequest(conn)
					code = 400
					break
				}

				// check if close connection
				for i := 2; i < len(reqSlice); i++ {
					kv := strings.Split(reqSlice[i], ":")
					if len(kv) != 2 {
						code = 400
						hs.handleBadRequest(conn)
						break
					}
					if strings.HasPrefix(reqSlice[i], "Connection:") {
						fmt.Println("handel connection2")
						idxH := strings.Index(reqSlice[i], ":")
						msgH := reqSlice[i][idxH+1:]
						connection := strings.TrimSpace(msgH)
						fmt.Println("conn msg:" + connection)
						if connection == "close" {
							res.connection = "close"
							conn.Close()
							fmt.Println("Connection closed by request.")
							return
						} else {
							fmt.Println("not close")
							res.connection = "no"
						}
					}
				}
				fmt.Println("-=-=-=-=-=-=-=WRITING RESPONSE -=-=-=-=-=-=-=")
				fmt.Println(code)
			}

		}
	}
	// has data left in buffer, partial request
	if remaining != "" {
		hs.handleBadRequest(conn)
		fmt.Println("Partial request: " + remaining)
	}

}
