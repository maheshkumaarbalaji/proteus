package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/mkbworks/proteus/lib/http"
)

func main() {
	server := http.NewServer()
	CurrentDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Error occurred while getting current working directory: " + err.Error())
		os.Exit(1)
	}
	
	TargetDirectory := filepath.Join(CurrentDirectory, "Files")
	server.Static("/files", TargetDirectory)
	server.Get("/user/:name", func(req *http.HttpRequest, res *http.HttpResponse) error {
		names, _ := req.Segments.Get("name")
		fmt.Printf("The name value in the path is %s\n", strings.Join(names, ","))
		return nil
	})
	
	server.Listen(8080, "localhost")
}