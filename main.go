package main // import "github.com/CenturyLinkLabs/imagelayers"

import (
	"flag"

	"github.com/CenturyLinkLabs/imagelayers/server"
)

func main() {
	port := flag.Int("p", 9000, "port on which the server will run")
	flag.Parse()

	s := server.NewServer()
	s.Start(*port)
}
