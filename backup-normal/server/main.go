package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	addr := ":5000"

	http.HandleFunc("/upload", handler)
	http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 1000)
	bytesRead := 0
	numPart := r.URL.Query()["part_number"][0]
	for {
		n, err := r.Body.Read(buf)
		bytesRead += n

		fmt.Println("Numer part: ", numPart)
		fmt.Println("Number bytes: ", len(buf[:n]))
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading HTTP response: ", err.Error())
		}
	}
	fmt.Printf("read %d bytes\n", bytesRead)
}
