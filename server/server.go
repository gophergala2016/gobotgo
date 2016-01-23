package server

import "net/http"

func main() {
	http.ListenAndServe("", nil)
}
