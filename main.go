package main

import (
	"lenscape-chat/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(":9091")
}
