// Package proxy - a Simple and Powerful ReverseProxy in Go
// https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/
package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	targetUrl, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	proxy.ModifyResponse = modifyResponse()
	proxy.ErrorHandler = errorHandler()
	return proxy, nil
}

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		// return errors.New("response body is invalid")
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err.Error())
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		for k, v := range resp.Header {
			log.Printf("  RespHeader field %q, Value %q", k, v)
		}
		log.Printf("Response body: %v", len(buf))
		reader := io.NopCloser(bytes.NewBuffer(buf))
		resp.Body = reader

		return nil
	}
}

// RequestHandler handles the http request using proxy
func RequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving %s %s %s", r.Method, r.URL, r.Proto)
		// Iterate over all header fields
		for k, v := range r.Header {
			log.Printf("  ReqHeader field %q, Value %q", k, v)
		}
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Request body: %v", string(buf))
		reader := io.NopCloser(bytes.NewBuffer(buf))
		r.Body = reader
		proxy.ServeHTTP(w, r)
	}
}

func NewReverseProxy() {
	// initialize a reverse proxy and pass the actual backend server url here
	proxyAddr := "https://trivy.tsc-np.signintra.com"
	localPort := 8081
	proxy, err := NewProxy(proxyAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("Starting reverse proxy for %s on localPort %d", proxyAddr, localPort)
	// handle all requests to your server using the proxy
	http.HandleFunc("/", RequestHandler(proxy))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", localPort), nil))
}
