package main

import (
	"github.com/Irooniam/sotailc/internal"
)

func main() {
	console := internal.NewConsole()
	console.Run()
	/*
		console.SetLayout()
		if err := console.App.SetRoot(console.Form, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
			panic(err)
		}
	*/

}
