package main

import (
	"go-api-opt/internal/api"
	"log"
	"net/http"
)

func main(){
	mux:=http.NewServeMux()
	api.NewHandler().Register(mux)

	addr:=":8080"
	log.Printf("api listening on %s",addr)

	if err:=http.ListenAndServe(addr, mux);err!=nil{
		log.Fatal(err)
	}
}