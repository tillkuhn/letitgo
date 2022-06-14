package httpserver

import "net/http"

var uiRoot = "https://my.ui"

// AddCORSHeader set CORS http headers, so we can call methods from external domains (e.g. to pull reports)
func AddCORSHeader(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", uiRoot)
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
