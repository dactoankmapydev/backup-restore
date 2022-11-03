package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const maxS = 15 * 1024 * 1024

func uploadFile(filePath string) error {
	client := &http.Client{}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file!!!")
	}
	stats, _ := file.Stat()
	fileSize := stats.Size()
	defer file.Close()

	nBytes, nChunks := int64(0), int64(0)
	r := bufio.NewReader(file)
	buf := make([]byte, maxS)

	remaining := int(fileSize)
	var partNum = 0
	var currentSize int

	for {
		if remaining < maxS {
			currentSize = remaining
		} else {
			currentSize = maxS
		}

		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)

		fileWriter, err := bodyWriter.CreateFormFile("file", file.Name())
		if err != nil {
			log.Println("error writing to buffer", err)
		}

		// read content to buffer
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		nChunks++
		nBytes += int64(len(buf))

		remaining -= currentSize
		partNum++

		log.Printf("Part %v complete, read %d bytes\n", partNum, len(buf))
		//log.Println(string(buf))

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		// io copy
		_, err = io.Copy(fileWriter, bytes.NewReader(buf))
		if err != nil {
			log.Println(err)
		}
		contentType := bodyWriter.FormDataContentType()
		_ = bodyWriter.Close()

		uri := fmt.Sprintf("http://localhost:5000/upload?part_number=%d", partNum)
		req, err := http.NewRequest("POST", uri, bytes.NewReader(buf))
		req.Header.Set("Content-Type", contentType)
		if err != nil {
			log.Println(err)
		}
		// log.Println("url", req.URL)

		_, err = client.Do(req)

		if err != nil {
			log.Fatal("Error doing request:", err)
		}
	}
	log.Println("Total bytes: ", nBytes, "Total chunks: ", nChunks)
	return nil
}

func main() {
	filename := "/home/toannd2/cms.zip"
	_ = uploadFile(filename)
}
