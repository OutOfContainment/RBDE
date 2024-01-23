package main

import (
	"image/color"
	"log"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func clockUpdate_go(currTime *widget.Label) {
	for range time.Tick(time.Second) {
		time := time.Now().Format(time.TimeOnly)
		currTime.SetText(time)
	}

}

func skeleton(RBDE fyne.App, win fyne.Window) *fyne.Container {
	// add image
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
	go clockUpdate_go(currTime)
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

	screen := container.NewPadded(
		canvas.NewRectangle(color.NRGBA{R: 127, G: 20, B: 60, A: 155}),
		container.NewBorder(bar, nil, nil, nil, list),
	)

	buttons := container.NewCenter(
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
	skeleton := container.NewGridWithRows(
		2,
		screen,
		container.NewVBox(
			title,
			buttons,
		),
	)
	return skeleton
}
