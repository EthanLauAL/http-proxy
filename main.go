package main

import (
	"encoding/base64"
	"fmt"
	"flag"
	"io"
	"log"
	"net/http"
)

var (
	argListen string
	argUserPass string
	basicAuthUserPass string
)

func main() {
	flag.StringVar(&argListen, "listen", "", "Proxy listen, required.")
	flag.StringVar(&argUserPass, "userpass", "", "user:pass")
	flag.Parse()

	if argListen == "" {
		log.Fatalln("-listen argument required.")
	}

	if argUserPass != "" {
		basicAuthUserPass = "Basic " + base64.StdEncoding.EncodeToString(
			[]byte(argUserPass))
	}

	http.HandleFunc("/", handleProxy)
	log.Fatalln(http.ListenAndServe(argListen, nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Host == "" { //Request to this server.
		http.NotFound(w,r)
		return
	}

	if r.Method == "CONNECT" {
		r.URL.Scheme = "http"
		r.URL.Opaque = "https"
	}

	if basicAuthUserPass != "" {
		auth := r.Header.Get("Proxy-Authorization")
		r.Header.Del("Proxy-Authorization")
		if auth != basicAuthUserPass {
			w.WriteHeader(http.StatusProxyAuthRequired)
			fmt.Fprintln(w, "ProxyAuthRequired")
			return
		}
		
	}

	resp,err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Println(r.Method, r.URL, err)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Println(r.Method, r.URL, resp.Status)
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
