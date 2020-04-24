package main

import "os"

const containerDir = "/neuroide"

func main() {
	_, err := os.Create(containerDir + "/test")
	if err != nil {
		panic(err)
	}
}