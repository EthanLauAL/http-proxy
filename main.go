package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleProxy)
	http.ListenAndServe(":3128", nil)
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Host == "" { //Request to this server.
		http.NotFound(w,r)
		return
	}

	resp,err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write([]byte(err.Error()))
	} else {
		//Copy headers
		for key,values := range resp.Header {
			for _,value := range values {
				w.Header().Add(key,value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w,resp.Body)
	}
}
