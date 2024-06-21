package main

import "github.com/obimadu/ipc3-stage-2/internals/config"

const webPort string = ":80"

func main() {
	// Init
	config.Config()

	// Run http server
	router().Run(webPort)
}
