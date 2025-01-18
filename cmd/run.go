package main

import (
	"github.com/Irooniam/sotailc/internal"
)

func main() {
	console := internal.NewConsole()
	console.SetLayout()
	console.App.Run()
}
