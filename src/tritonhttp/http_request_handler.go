package tritonhttp

import (
	"bufio"
	"fmt"
	"io"
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
					// w.WriteString("400 Bad Request")
					// w.Flush()
					// conn.Close()
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
						response += "HTTP/1.1 200 OK\r\n"
						response += "Server: Go-Triton-Server-1.0\r\n"

						if !strings.HasPrefix(firstR[1], hs.DocRoot) {
							url = hs.DocRoot + firstR[1]
						}

						idxFirstR := strings.LastIndex(firstR[1], "/")
						if firstR[1] == "/" {
							fmt.Println("url is 11/:")
							url = hs.DocRoot + firstR[1] + "index.html"
						} else if idxFirstR == len(firstR[1])-1 {
							// if strings.HasPrefix(firstR[1], "/index")
							fmt.Println("url is 22/:")
							url = hs.DocRoot + firstR[1] + "index.html"
						} else {
							// fmt.Println("url is " + url)
							url = hs.DocRoot + firstR[1]
						}
						// fmt.Println(url)

						lastIdx := strings.LastIndex(url, ".")
						extension := url[lastIdx:]

						if ext, ok := hs.MIMEMap[extension]; ok {
							res.contentType = ext
							fmt.Println("extension found: " + extension + " | " + ext)
						} else {
							res.contentType = "application/octet-stream"
							fmt.Println("extension not found: " + extension)
						}

						// lastIdx := strings.LastIndex(url, "/")
						// res.contentType = initialLine[1][lastIdx:]
						fi, err := os.Open(url)
						defer fi.Close()
						if err != nil {
							code = 404
							hs.handleFileNotFoundRequest(conn)
							fmt.Println(err)
							break
						}
						// get the size
						fiStat, _ := fi.Stat()
						res.contentLength = fiStat.Size()
						// io.Copy(w, fi)
						// w.Flush()
						// fmt.Println(res.contentType)
						// fmt.Println(size)

						res.lastModified = fiStat.ModTime().String()
						response += "Last-Modified: " + res.lastModified + "\r\n"
						response += "Content-Length: " + strconv.FormatInt(res.contentLength, 10) + "\r\n"
						response += "Content-Type: " + res.contentType + "\r\n\r\n"

						// check host
						secondLine := reqSlice[1]
						fmt.Println("second l: " + secondLine)
						if strings.HasPrefix(secondLine, "Host") {
							// idxH := strings.Index(msg, ":")
							// msgH := msg[idxH+1:]

							fmt.Println("Has Host")
							// w.WriteString(response)
							// w.Flush()
							// hs.sendResponse()
							code = 200
						} else {
							fmt.Println("error6")
							// hs.handleBadRequest(conn)
							code = 400
							break
						}
					}
				} else {
					code = 400
					fmt.Println("error5")
					fmt.Println(initialLine)
					hs.handleBadRequest(conn)
					break
				}

				// if len(reqSlice) > 2 {
				// 	fmt.Println("handel connection")
				// 	if strings.HasPrefix(reqSlice[2], "Connection:") {
				// 		fmt.Println("handel connection2")
				// 		idxH := strings.Index(reqSlice[2], ":")
				// 		msgH := reqSlice[2][idxH+1:]
				// 		connection := strings.TrimSpace(msgH)
				// 		fmt.Println("conn msg:" + connection)
				// 		if connection == "close" {
				// 			res.connection = "close"
				// 			// w.WriteString("Connection: closed\r\n")
				// 			// w.Flush()
				// 			conn.Close()
				// 			fmt.Println("Connection closed by request.")
				// 			return
				// 		} else {
				// 			fmt.Println("not close")
				// 			res.connection = "no"
				// 		}
				// 	}
				// }

				for i := 2; i < len(reqSlice); i++ {
					kv := strings.Split(reqSlice[i], ":")
					if len(kv) != 2 {
						code = 400
						hs.handleBadRequest(conn)
					} else {
						if strings.HasPrefix(reqSlice[i], "Connection:") {
							fmt.Println("handel connection2")
							idxH := strings.Index(reqSlice[i], ":")
							msgH := reqSlice[i][idxH+1:]
							connection := strings.TrimSpace(msgH)
							fmt.Println("conn msg:" + connection)
							if connection == "close" {
								res.connection = "close"
								// w.WriteString("Connection: closed\r\n")
								// w.Flush()
								conn.Close()
								fmt.Println("Connection closed by request.")
								return
							} else {
								fmt.Println("not close")
								res.connection = "no"
							}
						}
					}
				}

				// goodFormat := true
				// for i := 2; i < len(reqSlice); i++ {
				// 	idxKv := strings.Index(reqSlice[i], ":")
				// 	if idxKv == -1 {
				// 		fmt.Println("not valid header")
				// 		hs.handleBadRequest(conn)
				// 		goodFormat = false
				// 		break
				// 	}
				// 	// kv := strings.Split(reqSlice[i], ":")
				// 	fmt.Println("processing header")
				// 	if strings.HasPrefix(reqSlice[i], "Connection:") {
				// 		fmt.Println("handel connection2")
				// 		idxH := strings.Index(reqSlice[2], ":")
				// 		msgH := reqSlice[2][idxH+1:]
				// 		connection := strings.TrimSpace(msgH)
				// 		fmt.Println("conn msg:" + connection)
				// 		if connection == "close" {
				// 			res.connection = "close"
				// 			// w.WriteString("Connection: closed\r\n")
				// 			// w.Flush()
				// 			conn.Close()
				// 			fmt.Println("Connection closed by request.")
				// 			return
				// 		} else {
				// 			fmt.Println("not close")
				// 			res.connection = "no"
				// 		}
				// 	}
				// }
				// if goodFormat {
				// 	w.WriteString(response)
				// 	w.Flush()
				// }
				fmt.Println("-=-=-=-=-=-=-=WRITING RESPONSE -=-=-=-=-=-=-=")
				fmt.Println(code)
				if code == 200 {
					w.WriteString(response)
					w.Flush()
					fi, err := os.Open(url)
					defer fi.Close()
					if err != nil {
						code = 404
						// hs.handleFileNotFoundRequest(conn)
						fmt.Println(err)
						break
					}
					io.Copy(w, fi)
					w.Flush()
					// fmt.Println(res.contentType)
					// fmt.Println(size)
				} else if code == 400 {
					hs.handleBadRequest(conn)
				} else if code == 404 {
					hs.handleFileNotFoundRequest(conn)
				} else {
					fmt.Println("-=-=-=-=--=-=-==-=-=-=-==-error when handling requests-=-=-=-=--=-=-==-=-=-=-==-")
				}
			}

		}
	}
	// hs.handleBadRequest(conn)
}
