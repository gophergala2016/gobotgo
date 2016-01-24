package main

import "net/http"

func main() {
	http.ListenAndServe(":8101", http.FileServer(http.Dir("")))
}
