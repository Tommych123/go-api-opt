package main

import (
	"go-api-opt/internal/api"
	"log"
	_ "net/http/pprof"
	"net/http"
)

func main(){
	// API сервер на отдельном mux
	apiMux:=http.NewServeMux()
	api.NewHandler().Register(apiMux)

	//pprof регистрируется на DefaultServeMux из-за импорта net/http/pprof
	// поэтому admin сервер поднимаем с http.DefaultServeMux
	adminMux:=http.DefaultServeMux

	go func() {
		adminAddr:=":6060"
		log.Printf("pprof listening on %s", adminAddr)
		if err:=http.ListenAndServe(adminAddr, adminMux);err!=nil{
			log.Fatal(err)
		}
	}()
	apiAddr:=":8080"
	log.Printf("api listening on %s",apiAddr)

	if err:=http.ListenAndServe(apiAddr, apiMux);err!=nil{
		log.Fatal(err)
	}
}