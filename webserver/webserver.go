package webserver

import "net/http"

// Run Run webserver on specified port (passed as string the
// way regular http.ListenAndServe works)
func Run(addr string) {
	uh := getUrlsAndHandlers()
	for _, element := range uh {
		http.HandleFunc(element.url, element.handler)
	}
	http.ListenAndServe(addr, nil)
}
