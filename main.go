package main

import (
	"fmt"
	"lazytask/ui"
)

func main() {
	app := ui.InitializeApp()
	if err := app.Run(); err != nil {
		fmt.Println("Unable to run the app: ", err)
	}
}
