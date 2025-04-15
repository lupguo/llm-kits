package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/convert", conversionHandler)
	fmt.Println("MCP Service running on http://127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)
}
