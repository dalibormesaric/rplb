package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var delay = flag.Int("delay", 0, "Delay in milliseconds.")

func main() {
	flag.Parse()

	log.Println("Starting delay server on 8888")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if *delay > 0 {
			d := time.Duration(*delay * int(time.Millisecond))
			time.Sleep(d)
			fmt.Fprintf(w, "Response after %s milliseconds.\n", d)
		} else {
			fmt.Fprintf(w, "Response.\n")
		}
	})
	http.ListenAndServe(":8888", nil)
}
