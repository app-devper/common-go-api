package main

import (
	"mgo-gin/app"
)

func main() {
	var server app.Routes
	server.StartGin()
}
