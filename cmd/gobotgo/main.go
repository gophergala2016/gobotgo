// gobotgo is an API based gameroom for the playing of Go.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gophergala2016/gobotgo/server"
)

var port = flag.String("port", ":8100", "port to run service on")

func init() {
	flag.Parse()
}

func main() {
	log.Fatal(http.ListenAndServe(*port, server.MuxerAPIv1()))
}
