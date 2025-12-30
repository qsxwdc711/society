package main

import "society/ioc"

func main() {
	ioc.InitViper()

	app := InitWebServer()
	ioc.InitApiColl(app.server, app.DB)
	app.server.Run(":9999")
}
