package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/k0kubun/pp"
)

/*
func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>\n")
}
*/

func handler(w http.ResponseWriter, r *http.Request) {
	pp.Printf("URL: %s\n", r.URL.String())
	pp.Printf("Version: %v\n", r.Proto)
	pp.Printf("Method: %s\n", r.Method)
	pp.Printf("Header: %s\n", r.Header)
	pp.Printf("Forms: %v\n", r.Form)
	q := r.URL.Query()
	pp.Printf("Query: %v\n", q)
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("--------------------body-------------------\n%s\n-------------------------------------------\n", string(body))
	fmt.Fprintf(w, "<html><body>hello world</body></html>")
}

func cookieHandler(w http.ResponseWriter, r *http.Request) {
	// h := http.Header{}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	// header := w.Header()
	w.Header().Add("User-Agent", "spectre")
	if _, ok := r.Header["Cookie"]; ok {
		//クッキーがあれば一度は訪問済み
		fmt.Fprintf(w, "<html><body>2回目以降</body></html>\n")
	} else {
		fmt.Fprintf(w, "<html><body>初訪問</body></html>\n")
	}
}

func handlerDigest(w http.ResponseWriter, r *http.Request) {
	pp.Printf("URL: %s\n", r.URL.String())
	pp.Printf("Query: %v\n", r.Proto)
	pp.Printf("Method: %s\n", r.Method)
	pp.Printf("Header: %v\n", r.Header)
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("--body--\n%s\n", string(body))
	if _, ok := r.Header["Authorization"]; !ok {
		w.Header().Add("WWW-Authenticate", `Digest realm="Secret Zone", nonce="TgLc25U2BQA=f510a2780473e18e6587be702c2e67fe2b04afd", algorithm=MD5, qop="auth"`)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		fmt.Fprintf(w, "<html><body>secret page</body></html>\n")
	}
}

func thubmnailHndler(w http.ResponseWriter, r *http.Request) {
	buffer, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(buffer[:10]))
	file, err := os.OpenFile("img/thubmnail.jpg", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	file.Write(buffer)
	defer file.Close()
	pp.Printf("URL: %s\n", r.URL.String())
	pp.Printf("Query: %v\n", r.Proto)
	pp.Printf("Method: %s\n", r.Method)
	pp.Printf("Header: %v\n", r.Header)
	fmt.Printf("body\n\n%v\n", r.Body)
	defer r.Body.Close()
}

//echo server
func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)
	http.HandleFunc("/cookie", cookieHandler)
	http.HandleFunc("/digest", handlerDigest)
	http.HandleFunc("/img", thubmnailHndler)
	log.Printf("start https listening:18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}
