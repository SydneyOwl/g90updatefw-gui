package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.NewWithID("sydneyowl")
	gui := newFwgui()
	gui.loadUI(myApp)
	myApp.Run()
}
