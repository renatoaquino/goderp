package main

import (
	"github.com/renatoaquino/goderp"
)

func main() {
	gd := goderp.New()
	gd.Define("PORT", 8888, "Service Port", "Daemon")
	gd.Define("LOG_LEVEL", "info", "Log Level", "Logging")
	gd.Dump()
}
