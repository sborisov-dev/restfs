package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"rest-fs/src/fs"
	"rest-fs/src/middlewares"
	"rest-fs/src/routes"
)

var workDir = flag.String("workdir", "/tmp", "working directory for rest fs")

func main() {
	flag.Parse()
	dir := *workDir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("working directiry %s is not exists", dir)
	}
	manager := fs.NewLocalFileManager(dir)
	h := middlewares.Log(routes.NewRouter(manager))
	log.Println("RestFS starts listening :8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
