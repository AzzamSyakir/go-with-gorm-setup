package main

import (
	"fmt"
	"golang-api/api/routes"
)

func main() {
	fmt.Println("server start on port 9000")
	routes.StartServer()
}
