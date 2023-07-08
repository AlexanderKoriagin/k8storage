package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akrillis/k8storage/internal/entities"
)

//
// It's just a test module with a lot of printlns to demonstrate how the storage works.
//

const (
	textStart = `
We'll do two tests:

FIRST
1. Read the test01.json file
2. Put data into test storage
3. Get data from test storage
4. Print the data received from the test storage

SECOND
1. Read the storage.png file
2. Put data into test storage
3. Get data from test storage
4. Save the data received from the test storage as storage_from_storage.png`
)

func main() {
	fmt.Println(textStart)
	fmt.Println()
	fmt.Println()

	fmt.Println("*** FIRST TEST ***")
	firstTest()
	fmt.Println()
	fmt.Println()

	fmt.Println("*** SECOND TEST ***")
	secondTest()
}

func firstTest() {
	fmt.Println("1. Read the test01.json file")
	data, err := os.ReadFile("test01.json")
	if err != nil {
		log.Fatalf("couldn't read the file test01.json: %v", err)
	}
	fmt.Println(string(data))
	fmt.Println()

	fmt.Println("2. Put data into test storage")
	pfr := entities.PutFileRequest{
		ClientID: "1",
		Name:     "test01.json",
		Data:     data,
	}
	body, err := json.Marshal(pfr)
	if err != nil {
		log.Fatalf("couldn't marshal the request body: %v", err)
	}

	resp, err := http.Post(
		"http://127.0.0.1:58080/files",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Fatalf("couldn't put data into test storage: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("wrong status code: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	fmt.Println("Data has been put into test storage")
	fmt.Println()
	time.Sleep(200 * time.Millisecond)

	fmt.Println("3. Get data from test storage")
	resp, err = http.Get("http://127.0.0.1:58080/files?client_id=1&name=test01.json")
	if err != nil {
		log.Fatalf("couldn't put data into test storage: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("wrong status code: %d", resp.StatusCode)
	}

	var res entities.GetFileResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Fatalf("couldn't decode response body: %v", err)
	}
	_ = resp.Body.Close()

	fmt.Println("Data has been got from test storage")
	fmt.Println()

	fmt.Println("4. Print the data received from the test storage")
	fmt.Println(string(res.Data))
}

func secondTest() {
	fmt.Println("1. Read the storage.png file")
	data, err := os.ReadFile("storage.png")
	if err != nil {
		log.Fatalf("couldn't read the file storage.png: %v", err)
	}
	fmt.Println("Done")
	fmt.Println()

	fmt.Println("2. Put data into test storage")
	pfr := entities.PutFileRequest{
		ClientID: "1",
		Name:     "storage.png",
		Data:     data,
	}
	body, err := json.Marshal(pfr)
	if err != nil {
		log.Fatalf("couldn't marshal the request body: %v", err)
	}

	resp, err := http.Post(
		"http://127.0.0.1:58080/files",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Fatalf("couldn't put data into test storage: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("wrong status code: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	fmt.Println("Data has been put into test storage")
	fmt.Println()
	time.Sleep(200 * time.Millisecond)

	fmt.Println("3. Get data from test storage")
	resp, err = http.Get("http://127.0.0.1:58080/files?client_id=1&name=storage.png")
	if err != nil {
		log.Fatalf("couldn't put data into test storage: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("wrong status code: %d", resp.StatusCode)
	}

	var res entities.GetFileResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Fatalf("couldn't decode response body: %v", err)
	}
	_ = resp.Body.Close()

	fmt.Println("Data has been got from test storage")
	fmt.Println()

	fmt.Println("4. Save the data received from the test storage as storage_from_storage.png")
	err = os.WriteFile("storage_from_storage.png", res.Data, 0644)
	if err != nil {
		log.Fatalf("couldn't write the file floppy_from_storage.svg: %v", err)
	}
	fmt.Println("Data has been saved as storage_from_storage.png")
	fmt.Println()
}
