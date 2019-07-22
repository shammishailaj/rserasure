package main

import "net/http"

func main() {
	http.HandleFunc("/", authHandler)
	http.HandleFunc("/restricted", restrictedHandler)

	http.ListenAndServe(":3000", nil)
}
