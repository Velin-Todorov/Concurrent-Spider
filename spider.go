package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func saveHtmlFile(body []byte, file_path string) (interface{}, error) {
	file, err := os.Create(file_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to create file.")
	}

	f, err := file.Write(body)

	if err != nil{
		file.Close()
		return nil, fmt.Errorf("Failed to write to file")
	}

	err = file.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close file")
	}
	return f, nil
}


func spider(url string, file_path string, wg *sync.WaitGroup) (interface{}, error) {
	// this will be the concurrent fn
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error occurred when making the request")
		return nil, fmt.Errorf("Failed to build the request for %s", err)
	}

	defer resp.Body.Close() //this is used when we are done reading from the response body and we want to release the associated resources to prevent resource leaks

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request failed %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)

	f, err := saveHtmlFile(body, file_path)

	wg.Done()
	return f, nil

}


func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	url1 := os.Args[1]
	url2 := os.Args[2]

	go spider(url1,"page1.html", &wg)
	go spider(url2,"page2.html", &wg)

	wg.Wait()
}