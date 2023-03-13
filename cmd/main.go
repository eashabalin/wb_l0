package main

import "wb_l0/pkg/app"

func main() {
	server := app.NewApp()
	server.Run("7000")
}
