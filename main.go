package main

import (
	"image/color"
	"log"
	"path/filepath"
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
	win.Resize(fyne.NewSize(250, 500))

	apspath, err := filepath.Abs("images/Logo.png")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(apspath)

	image := canvas.NewImageFromFile(apspath)
	image.FillMode = canvas.ImageFillOriginal
	title := container.New(layout.NewCenterLayout(), image)

	// Clock timer at the top
	currTime := widget.NewLabel("")
	ClockUpdate(currTime)
	bar := container.New(
		layout.NewCenterLayout(),
		canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 55}),
		currTime,
	)

	var testData = [10]string{"1.wav", "2.wav", "3.wav"}
	list := widget.NewList(
		func() int {
			return 10
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(testData[i])
		})

	rect := canvas.NewRectangle(color.NRGBA{R: 127, G: 20, B: 60, A: 155})
	screen := container.NewPadded(
		rect,
		container.NewBorder(bar, nil, nil, nil, list),
	)

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
	GUI := container.NewGridWithRows(
		2,
		screen,
		container.NewVBox(
			title,
			buttons,
		),
	)
	win.SetContent(GUI)

	// Open window
	log.Println("Show window!")

	win.ShowAndRun()
	log.Println("Goodbye ... :(")
}
