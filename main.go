package main

import (
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ClockUpdate(currTime *widget.Label) {
	go func() {
		for range time.Tick(time.Second) {
			time := time.Now().Format(time.TimeOnly)
			currTime.SetText(time)
		}
	}()

}

func main() {
	log.Println("Start")

	// Initialise window
	RBDE := app.New()
	win := RBDE.NewWindow("DiEmu")
	win.Resize(fyne.NewSize(320, 640))

	title := container.New(layout.NewCenterLayout(), canvas.NewText("RBDE", color.White))

	// Clock timer at the top
	currTime := widget.NewLabel("")
	ClockUpdate(currTime)
	clock := container.New(layout.NewCenterLayout(), currTime)

	buttons := container.New(
		layout.NewCenterLayout(),
		container.NewVBox(
			container.NewHBox(
				widget.NewButton("<<", func() {
					log.Println("<<")
				}),
				widget.NewButton("Record", func() {
					log.Println("Play")
				}),
				widget.NewButton("Play/Stop", func() {
					log.Println("Stop")
				}),
				widget.NewButton(">>", func() {
					log.Println(">>")
				}),
			),
			container.NewHBox(
				widget.NewButton("Menu", func() {
					log.Println("Menu")
				}),
				layout.NewSpacer(),
				widget.NewButton("Delete", func() {
					log.Println("Delete")
				}),
			),
		),
	)
	// GUI layout
	GUI := container.NewVBox(
		title,
		clock,
		buttons,
	)
	win.SetContent(GUI)

	// Open window
	log.Println("Show window!")

	win.ShowAndRun()
	log.Println("Goodbye ... :(")
}
