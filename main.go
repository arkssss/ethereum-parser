package main

import (
	"ethereum-parser/app/api"
	"log"
	"net/http"
)

func main() {
	// init
	api.Setup()
	http.HandleFunc("/CurrentBlock", api.GetCurrentBlock) //设置访问的路由
	http.HandleFunc("/Subscribe", api.Subscribe)          //设置访问的路由
	http.HandleFunc("/Transactions", api.GetTransactions) //设置访问的路由
	err := http.ListenAndServe(":9090", nil)              //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
