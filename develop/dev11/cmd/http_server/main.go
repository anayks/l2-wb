package main

import (
	"fmt"

	httpserv "httpserv/internal/api"
)

func main() {
	err := httpserv.StartHTTP()
	if err != nil {
		fmt.Printf("Error on starting http: %v", err)
	}
	fmt.Printf("HTTP server ended work.")
}
