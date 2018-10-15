package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func checkStatus() (response string) {
	url := "http://localhost:9200/"

	resp, err := http.Get(url)
	if err != nil {
		Fatal.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	response = string(body)
	return
}

func deleteIndex(name string) (response string) {
	url := "http://localhost:9200/" + name

	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("X-Custom-Header", "Go-elkClient")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Fatal.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	response = string(body)
	return
}

func createIndex(name string, jsonData string) (response string) {
	url := "http://localhost:9200/" + name

	var jsonStr = []byte(jsonData)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "Go-elkClient")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Fatal.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	response = string(body)
	return
}

func postDoc(name string, jsonData string) (response string) {
	url := "http://localhost:9200/" + name + "/doc/"

	var jsonStr = []byte(jsonData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "Go-elkClient")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Fatal.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	response = string(body)
	return
}
