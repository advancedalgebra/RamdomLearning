package main

import (
	"RamdomLearning/conf"
	"fmt"
	"net/http"
)

func main() {
	routers := InitRouter()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: routers,
	}
	fmt.Println(conf.Port)
	server.ListenAndServe()
}