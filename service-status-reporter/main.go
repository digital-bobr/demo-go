package main

import (
	"service-status-reporter/pkg/router"
)

func main() {
	router.GetInstance().Run(":8080")
}
