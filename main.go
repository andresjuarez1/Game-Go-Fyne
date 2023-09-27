package main

import (
    "fyne.io/fyne/v2/app"
    "juego2/juego"
)

func main() {
    myApp := app.New()
    w := myApp.NewWindow("Game")
	juego.GameMain(w)
    w.CenterOnScreen()
    w.ShowAndRun()
}
