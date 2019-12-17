package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var db *mgo.Database
var router = mux.NewRouter()

func setResponseHeader(response http.ResponseWriter) {
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
}
