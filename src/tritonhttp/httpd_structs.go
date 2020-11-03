package tritonhttp

type HttpServer struct {
	ServerPort string
	DocRoot    string
	MIMEPath   string
	MIMEMap    map[string]string
}

type HttpResponseHeader struct {
	// Add any fields required for the response here
	server        string
	lastModified  string
	contentType   string
	contentLength int64
	connection    string
}

type HttpRequestHeader struct {
	// Add any fields required for the request here
	host       string
	connection string
	others     bool
}
